package message

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nikoksr/notify"
	"github.com/panjf2000/ants/v2"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/message/channel"
	"github.com/usual2970/retell/pkg/config"
	"github.com/usual2970/retell/pkg/jwt"
	"github.com/usual2970/retell/pkg/logger"
	"github.com/usual2970/retell/pkg/sse"
	"github.com/usual2970/retell/pkg/utils"

	"golang.org/x/sync/errgroup"
)

const sseSystemPushTopic = "new_system_message:%s"

type noti struct {
	n *notify.Notify

	title, content, link string
}

type MessageRepository interface {
	GetTaskByUri(ctx context.Context, uri string) (*domain.MessageTask, error)
	BatchInsertHistory(ctx context.Context, tasks []*domain.MessageHistory) error
	Total(ctx context.Context, filters map[string]any) (int64, error)
	List(ctx context.Context, where map[string]any, pagination *constant.Pagination, order string) ([]domain.MessageHistory, error)
	Updates(ctx context.Context, where map[string]any, updates map[string]any) error
}

type UserAccountRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*domain.UserAccount, error)
}

type SettingRepository interface {
	GetSetting(ctx context.Context, key string, userId ...int64) (*domain.Setting, error)
	SetSetting(ctx context.Context, setting *domain.Setting) error
}

type Service struct {
	msgRepo         MessageRepository
	historyChan     chan *domain.MessageHistory
	processCancel   context.CancelFunc
	userAccountRepo UserAccountRepository
	settingRepo     SettingRepository
}

var service *Service
var onceService sync.Once

func NewService(msgRepo MessageRepository, userAccountRepo UserAccountRepository, settingRepo SettingRepository) *Service {
	onceService.Do(func() {
		service = &Service{
			msgRepo:         msgRepo,
			userAccountRepo: userAccountRepo,
			settingRepo:     settingRepo,
			historyChan:     make(chan *domain.MessageHistory, 1),
		}

		ctx, cancel := context.WithCancel(context.Background())
		service.processCancel = cancel
		service.processHistory(ctx)
	})

	return service
}

func (s *Service) GetUnreadNum(ctx context.Context) (int64, error) {
	userId, err := jwt.GetUserID(ctx)
	if err != nil {
		return 0, err
	}

	where := map[string]any{
		"user_id":   userId,
		"readed_at": utils.IsNull,
		"type":      domain.MessageTplTypeSystem,
	}

	total, err := s.msgRepo.Total(ctx, where)
	if err != nil {
		return 0, constant.NewXError(100, "failed to get total unread messages: %v", err)
	}

	return total, nil
}

func (s *Service) List(ctx context.Context, req *domain.MessageListReq) (*constant.PaginatedResponse[domain.MessageInfoResp], error) {
	userId, err := jwt.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	retentionDays, err := s.getRetentionDays(ctx)
	if err != nil {
		return nil, constant.NewXError(100, "failed to get retention days: %v", err)
	}

	retentionBegin := time.Now().AddDate(0, 0, -retentionDays).UTC()

	where := map[string]any{
		"user_id":      userId,
		"type":         domain.MessageTplTypeSystem,
		"created_at >": retentionBegin,
	}

	if req.CreatedAtStart != 0 {
		where["created_at >"] = time.Unix(req.CreatedAtStart, 0).UTC()
	}
	if req.CreatedAtEnd != 0 {
		where["created_at <"] = time.Unix(req.CreatedAtEnd, 0).UTC()
	}
	if req.KWD != "" {
		where["content like"] = "%" + req.KWD + "%"
	}
	if req.Level != 0 {
		where["level"] = req.Level
	}
	if req.Readed != nil {
		if *req.Readed {
			where["readed_at"] = utils.IsNotNull
		} else {
			where["readed_at"] = utils.IsNull
		}
	}

	msgs, err := s.msgRepo.List(ctx, where, &req.Pagination, "created_at desc")
	if err != nil {
		return nil, constant.NewXError(100, "failed to list messages: %v", err)
	}

	total, err := s.msgRepo.Total(ctx, where)
	if err != nil {
		return nil, constant.NewXError(100, "failed to get total messages: %v", err)
	}

	items := make([]domain.MessageInfoResp, 0, len(msgs))
	for _, msg := range msgs {
		msg := &msg
		items = append(items, *msg.Trans2InfoResp())
	}

	return constant.NewPaginatedResponse(&req.Pagination, items, total), nil

}
func (s *Service) BatchSetReaded(ctx context.Context) error {
	userId, err := jwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	where := map[string]any{
		"user_id":   userId,
		"readed_at": utils.IsNull,
		"type":      domain.MessageTplTypeSystem,
	}

	updates := map[string]any{
		"readed_at": time.Now().UTC(),
	}

	if err := s.msgRepo.Updates(ctx, where, updates); err != nil {
		return constant.NewXError(100, "failed to batch set readed: %v", err)
	}

	return nil
}

func (s *Service) Read(ctx context.Context, req *domain.MessageReadReq) error {
	if req.ID <= 0 {
		return constant.NewXError(100, "invalid message ID")
	}
	userId, err := jwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	msgs, _ := s.msgRepo.List(ctx, map[string]any{
		"id":      req.ID,
		"user_id": userId,
	}, nil, "id desc")

	if len(msgs) == 0 {
		return constant.NewXError(100, "msg not exist")
	}

	if err := s.msgRepo.Updates(ctx, map[string]any{
		"id":      req.ID,
		"user_id": userId,
	}, map[string]any{
		"readed_at": time.Now().UTC(),
	}); err != nil {
		return constant.NewXError(100, "failed to set message readed: %v", err)
	}

	return nil

}

func (s *Service) SetRetentionDays(ctx context.Context, req *domain.MessageSetRetentionDaysReq) error {
	userId, err := jwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	setting, err := s.settingRepo.GetSetting(ctx, domain.SettingMessageRetentionDaysKey, userId)
	if err != nil {
		setting = domain.InitSetting(domain.SettingMessageRetentionDaysKey, userId)
	}

	setting.SetValue(&domain.SettingMessageRetentionDays{
		Days: req.Days,
	})

	if err := s.settingRepo.SetSetting(ctx, setting); err != nil {
		return constant.NewXError(100, "failed to set retention days: %v", err)
	}

	return nil
}

func (s *Service) Send(ctx context.Context, req *domain.MessageSendReq) (*domain.MessageSendResp, error) {

	task, err := s.msgRepo.GetTaskByUri(ctx, req.TaskURI)
	if err != nil {
		return nil, err
	}

	toUserInfo, err := s.buildToUserInfo(req)
	if err != nil {
		return nil, constant.NewXError(100, "failed to build to user info: %v", err)
	}

	notifiers, err := s.buildNotifiers(ctx, req, task, toUserInfo)
	if err != nil {
		return nil, constant.NewXError(100, "failed to build notifiers: %v", err)
	}

	histories, err := s.buildHistories(ctx, req, task, toUserInfo)
	if err != nil {
		return nil, constant.NewXError(100, "failed to build histories: %v", err)
	}

	historyCtx := context.Background()
	s.pushMsg(historyCtx, histories...)

	var eg errgroup.Group
	for _, noti := range notifiers {
		n := noti
		eg.Go(func() error {
			return n.n.Send(ctx, n.title, n.content)
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// 如果是系统消息，需要额外推送消息更新通知
	s.sendSystemNotice(context.Background(), task, toUserInfo)

	return &domain.MessageSendResp{}, nil
}

type SystemNotice struct {
	TotalUnRead int64 `json:"totalUnread"`
}

func (s *Service) sendSystemNotice(ctx context.Context, task *domain.MessageTask, toUserInfo *domain.ToUserInfo) {
	ants.Submit(func() {
		topic := fmt.Sprintf(sseSystemPushTopic, toUserInfo.UserUri)
		lg := logger.WithField("topic", topic).WithField("module", "message.sendSystemNotice")
		systemIndex := utils.FindIndex(task.MessaageTpls, func(tpl domain.MessageTpl) bool {
			return tpl.Type == domain.MessageTplTypeSystem
		})

		if systemIndex < 0 {
			return
		}

		retentionDays, err := s.getRetentionDays(ctx, toUserInfo.UserId)
		if err != nil {
			lg.WithField("error", err).Error("failed to get retention days")
			return
		}

		retentionBegin := time.Now().AddDate(0, 0, -retentionDays).UTC()

		total, err := s.msgRepo.Total(ctx, map[string]any{
			"readed_at":    utils.IsNull,
			"user_id":      toUserInfo.UserId,
			"created_at >": retentionBegin,
		})
		if err != nil {
			lg.WithField("error", err).Error("failed to get total unread messages")
			return
		}

		msg := SystemNotice{
			TotalUnRead: total,
		}

		content, _ := json.Marshal(msg)

		if err := sse.Publish(topic, string(content)); err != nil {
			logger.WithField("err", err).WithField("content", string(content)).Error("failed to publish system notice")
		}

	})
}

func (s *Service) buildHistories(_ context.Context, param *domain.MessageSendReq, task *domain.MessageTask, toUserInfo *domain.ToUserInfo) ([]*domain.MessageHistory, error) {

	if err := task.CheckParam(param.Param); err != nil {
		return nil, err
	}

	tpls := task.MessaageTpls
	if len(tpls) == 0 {
		return nil, constant.ErrRecordNotFound
	}

	histories := make([]*domain.MessageHistory, 0, len(tpls))

	// 只有系统消息需要保存到数据库
	systemTpls := utils.Filter(tpls, func(tpl domain.MessageTpl, index int) bool {
		return tpl.Type == domain.MessageTplTypeSystem
	})

	for _, tpl := range systemTpls {

		title := tpl.GetTitle(param.Param)
		content, err := tpl.GetContent(param)
		if err != nil {
			return nil, err
		}

		link, err := tpl.GetLink(param)
		if err != nil {
			return nil, err
		}
		toUser, err := tpl.GetToUser(toUserInfo)
		if err != nil {
			return nil, err
		}
		history := domain.NewMessageHistory(&domain.InitMessageHistoryParams{
			Task:     task,
			Tpl:      &tpl,
			UserID:   toUserInfo.UserId,
			UserUri:  toUserInfo.UserUri,
			Receiver: toUser,
			Title:    title,
			Content:  content,
			Link:     link,
		})

		histories = append(histories, history)
	}

	return histories, nil
}

func (s *Service) buildNotifiers(_ context.Context, param *domain.MessageSendReq, task *domain.MessageTask, toUserInfo *domain.ToUserInfo) ([]noti, error) {

	if err := task.CheckParam(param.Param); err != nil {
		return nil, err
	}

	tpls := task.MessaageTpls
	if len(tpls) == 0 {
		return nil, constant.ErrRecordNotFound
	}

	rs := make([]noti, 0, len(tpls))
	for _, tpl := range tpls {
		notifier := notify.New()
		title := tpl.GetTitle(param.Param)
		content, err := tpl.GetContent(param)
		if err != nil {
			return nil, err
		}

		link, err := tpl.GetLink(param)
		if err != nil {
			return nil, err
		}
		toUser, err := tpl.GetToUser(toUserInfo)
		if err != nil {
			return nil, err
		}
		// 替换toUser
		title = strings.ReplaceAll(title, "{{toUser}}", toUser)
		content = strings.ReplaceAll(content, "{{toUser}}", toUser)

		// 如果是系统消息只需要保存到数据库就行了
		if tpl.Type == domain.MessageTplTypeSystem {
			continue
		}

		sendChannel, err := s.getChannel(tpl.Type, toUser, &channel.ChannelConfig{
			Task: task,
			Tpl:  &tpl,
		})
		if err != nil {
			return nil, err
		}

		notifier.UseServices(sendChannel)

		n := noti{
			n:       notifier,
			title:   title,
			content: content,
			link:    link,
		}

		rs = append(rs, n)
	}

	return rs, nil
}

func (s *Service) buildToUserInfo(param *domain.MessageSendReq) (*domain.ToUserInfo, error) {
	rs := param.ToUserInfo
	if rs == nil {
		rs = &domain.ToUserInfo{}
	}
	if param.UserID == "" {
		return rs, nil
	}

	userID, err := strconv.Atoi(param.UserID)
	if err != nil {
		return rs, nil
	}

	userAccount, err := s.userAccountRepo.GetByUserID(context.Background(), int64(userID))
	if err != nil {
		return rs, nil
	}

	if rs.Email == "" {
		rs.Email = userAccount.UserPrivateInfo.Email
	}

	rs.UserId = userAccount.UserID

	rs.UserUri = userAccount.UserProfile.URI

	return rs, nil
}

func (s *Service) pushMsg(_ context.Context, msgs ...*domain.MessageHistory) {
	ants.Submit(func() {
		for _, msg := range msgs {
			s.historyChan <- msg
		}
	})

}

func (s *Service) processHistory(ctx context.Context) {
	l := logger.WithField("module", "message.processHistory")
	for i := 0; i < 5; i++ {
		ants.Submit(func() {
			tick := time.NewTicker(200 * time.Millisecond)
			defer tick.Stop()
			jobs := make([]*domain.MessageHistory, 0)
			canDo := false
			for {
				select {
				case msgHistory := <-s.historyChan:
					jobs = append(jobs, msgHistory)
					canDo = len(jobs) >= 50
				case <-tick.C:
					canDo = len(jobs) > 0
				case <-ctx.Done():
					return
				}

				if !canDo {
					continue
				}
				if err := s.msgRepo.BatchInsertHistory(ctx, jobs); err != nil {
					l.WithField("jobs", jobs).WithField("err", err).Error("batch insert failed")
				}
				jobs = jobs[0:0]
			}
		})

	}
}

func (s *Service) getChannel(t domain.MessageTplType, receiver string, channelConfig *channel.ChannelConfig) (notify.Notifier, error) {
	var n notify.Notifier
	var err error
	switch t {
	case domain.MessageTplTypeEmail:
		conf := config.GetConfig().EmailAuth
		rs := channel.NewMail(conf.Username, conf.Host+":"+conf.Port)
		rs.AuthenticateSMTP(conf.Username, conf.Username, conf.Password, conf.Host)
		rs.AddReceivers(receiver)
		rs.BodyFormat(channel.HTML)
		rs.SetTLS(&tls.Config{
			ServerName: conf.Host,
		})
		n = rs
	case domain.MessageTplTypeFeishu:
		n = channel.NewFeishu(receiver)
	case domain.MessageTplTypeSystemPush:
		n = channel.NewSystemPush(receiver, channelConfig)
	default:
		err = constant.ErrChannelNotSupported
	}

	return n, err
}

func (s *Service) getRetentionDays(ctx context.Context, userIds ...int64) (int, error) {
	var userId int64
	if len(userIds) == 0 {
		var err error
		userId, err = jwt.GetUserID(ctx)
		if err != nil {
			return 0, err
		}
	} else {
		userId = userIds[0]
	}

	setting, err := s.settingRepo.GetSetting(ctx, domain.SettingMessageRetentionDaysKey, userId)
	if err != nil {
		return 90, nil
	}

	return setting.GetRetentionDays(), nil
}
