package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/panjf2000/ants/v2"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/config"
	"github.com/usual2970/retell/pkg/encode"
	"github.com/usual2970/retell/pkg/logger"
	"github.com/usual2970/retell/pkg/redlock"
	"github.com/usual2970/retell/pkg/utils"

	pkgJwt "github.com/usual2970/retell/pkg/jwt"
)

type CodeService interface {
	GenerateCode(ctx context.Context, param *domain.CodeGenerateReq) (*domain.Code, error)
	CheckCode(ctx context.Context, param *domain.CodeCheckReq) error
}

type MessageService interface {
	Send(ctx context.Context, req *domain.MessageSendReq) (*domain.MessageSendResp, error)
}

type AccessTokenRepository interface {
	SetAccessToken(ctx context.Context, accessToken string, userID int64) error
	DelAccessToken(ctx context.Context, token string) error
}

type Service struct {
	codeSvc     CodeService
	msgSvc      MessageService
	atRepo      AccessTokenRepository
	accountRepo AccountRepository
}

func NewService(
	codeSvc CodeService,
	msgSvc MessageService,
	atRepo AccessTokenRepository,
	accountRepo AccountRepository,

) *Service {
	return &Service{
		codeSvc:     codeSvc,
		msgSvc:      msgSvc,
		atRepo:      atRepo,
		accountRepo: accountRepo,
	}
}

func (s *Service) Login(ctx context.Context, param *domain.AuthLoginReq) (*domain.AuthLoginResp, error) {
	l := logger.WithField("module", "auth.Login").WithField("param", param)

	// 加锁
	locKey := fmt.Sprintf("lock_login_%s", param.Identifier)
	locVal, err := redlock.Lock(ctx, locKey, time.Second*2)
	if err != nil {
		return nil, constant.ErrOPerateTooFast
	}

	defer redlock.UnLock(ctx, locKey, locVal)

	loginer := NewLoginer(param.Type, s.accountRepo, s.codeSvc)
	if loginer == nil {
		l.Error("unsupported login type")
		return nil, constant.ErrUnsupportedLoginType
	}

	account, err := loginer.Login(ctx, param)
	if err != nil {
		l.WithField("err", err).Error("login failed")
		return nil, err
	}

	expiresAt := time.Now().Add(constant.AuthExpireDuration)
	claims := &jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%d", account.UserID),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	conf := config.GetConfig().Auth
	key := conf.JwtSecret
	if key == "" {
		l.Error("get auth key failed:")
		return nil, errors.New("get auth key failed")
	}
	accessToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}
	rs := &domain.AuthLoginResp{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt.Unix(),
		Extend: map[string]string{
			"openid": account.Openid,
			"role":   fmt.Sprint(account.UserProfile.Role),
		},
	}

	err = s.atRepo.SetAccessToken(ctx, accessToken, account.UserID)
	if err != nil {
		l.WithField("err", err).Error("save accesstoken failed")
		return nil, err
	}

	return rs, nil
}
func (s *Service) Register(ctx context.Context, param *domain.AuthRegisterReq) (*domain.AuthRegisterResp, error) {
	// 1. check code
	if err := s.codeSvc.CheckCode(ctx, &domain.CodeCheckReq{
		Receiver: param.Email,
		Code:     param.Code,
		Purpose:  domain.CodePurposeRegister,
	}); err != nil {
		return nil, err
	}

	// 2. check email exist
	openid := encode.Openid(domain.UserAccountPlatformEmail, param.Email)
	_, err := s.accountRepo.GetOneByOpenid(ctx, openid)
	if err != nil && !errors.Is(err, constant.ErrRecordNotFound) {
		return nil, err
	}

	// 3. if not exist, create account. else update account
	if err != nil {
		account := domain.NewAccount(openid, domain.UserAccountPlatformEmail)
		profile := domain.InitProfile(param.Email, param.Role)

		privateInfo := domain.NewPrivateInfoWithPassword(param.Email, param.Password)

		if err := s.accountRepo.Save(ctx, account, profile, privateInfo, func(account *domain.UserAccount) bool {
			if account.UserID == 0 {
				account.UserID = account.ID
			}
			if profile != nil {
				profile.UserID = account.ID
			}
			if privateInfo != nil {
				privateInfo.UserID = account.ID
			}
			return false
		}); err != nil {
			return nil, err
		}
	} else {
		return nil, constant.ErrAccountAlreadyExists
	}

	s.accountRepo.DeleteFromCache(ctx, openid)

	return nil, nil
}
func (s *Service) Logout(ctx context.Context) error {
	token, err := pkgJwt.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	return s.atRepo.DelAccessToken(ctx, token)
}
func (s *Service) LoginCode(ctx context.Context, param *domain.AuthCodeReq) error {
	return s.sendCode(ctx, param, domain.CodePurposeLogin)
}
func (s *Service) RegisterCode(ctx context.Context, param *domain.AuthCodeReq) error {
	// 1. check email exist
	openid := encode.Openid(domain.UserAccountPlatformEmail, param.Receiver)
	account, err := s.accountRepo.GetOneByOpenid(ctx, openid)
	if err != nil && !errors.Is(err, constant.ErrRecordNotFound) {
		return err
	}
	if err == nil && account != nil {
		return constant.ErrAccountAlreadyExists
	}
	return s.sendCode(ctx, param, domain.CodePurposeRegister)
}

func (s *Service) ForgetPasswordCode(ctx context.Context, param *domain.AuthCodeReq) error {
	return s.sendCode(ctx, param, domain.CodePurposeForget)
}
func (s *Service) ForgetPassword(ctx context.Context, param *domain.AuthForgetPasswordReq) error {
	// 1. check code
	if err := s.codeSvc.CheckCode(ctx, &domain.CodeCheckReq{
		Receiver: param.Email,
		Code:     param.Code,
		Purpose:  domain.CodePurposeForget,
	}); err != nil {
		return err
	}

	// 2. check email exist
	openid := encode.Openid(domain.UserAccountPlatformEmail, param.Email)
	account, err := s.accountRepo.GetOneByOpenid(ctx, openid)
	if err != nil && !errors.Is(err, constant.ErrRecordNotFound) {
		return err

	}
	if err != nil {
		return constant.ErrAccountNotExists
	}

	// 3. update password
	privateInfo := account.UserPrivateInfo
	privateInfo.SetPassword(param.Password)
	if err := s.accountRepo.Save(ctx, account, nil, privateInfo, func(account *domain.UserAccount) bool {
		if privateInfo != nil {
			privateInfo.UserID = account.ID
		}
		return true
	}); err != nil {
		return err
	}

	return nil
}
func (s *Service) ResetPassword(ctx context.Context, param *domain.AuthResetPasswordReq) error {
	userID, err := pkgJwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	account, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, constant.ErrRecordNotFound) {
		return err
	}

	if err != nil {
		return constant.ErrAccountNotExists
	}

	privateInfo := account.UserPrivateInfo

	privateInfo.SetPassword(param.Password)

	if err := s.accountRepo.Save(ctx, account, nil, privateInfo, nil); err != nil {
		return err
	}

	s.accountRepo.DeleteFromCache(ctx, account.Openid)
	s.accountRepo.DelByUserIdFromCache(ctx, account.UserID)

	return nil
}

func (s *Service) ResetHeadImg(ctx context.Context, param *domain.AuthResetHeadImgReq) error {
	userID, err := pkgJwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	account, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, constant.ErrRecordNotFound) {
		return err
	}

	if err != nil {
		return constant.ErrAccountNotExists
	}

	userProfile := account.UserProfile
	userProfile.Headimgurl = utils.ParseUrl(param.Headimgurl)

	if err := s.accountRepo.Save(ctx, account, userProfile, nil, nil); err != nil {
		return err
	}

	s.cacheClear(account)

	return nil
}

func (s *Service) cacheClear(account *domain.UserAccount) error {
	ants.Submit(func() {
		ctx := context.Background()

		s.accountRepo.DeleteFromCache(ctx, account.Openid)
		s.accountRepo.DelByUserIdFromCache(ctx, account.UserID)

	})

	return nil
}

func (s *Service) sendCode(ctx context.Context, param *domain.AuthCodeReq, purpose domain.CodePurpose) error {
	// 1. 生成验证码
	code, err := s.codeSvc.GenerateCode(ctx, &domain.CodeGenerateReq{
		Receiver: param.Receiver,
		Purpose:  purpose,
	})
	if err != nil {
		return err
	}
	// 2. 发送验证码
	_, err = s.msgSvc.Send(ctx, &domain.MessageSendReq{
		UserID:  "",
		TaskURI: "check_code_email",
		ToUserInfo: &domain.ToUserInfo{
			Email: param.Receiver,
		},
		Param: map[string]string{
			"code": code.Code,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
