package routes

import (
	"github.com/usual2970/retell/auth"
	"github.com/usual2970/retell/code"
	"github.com/usual2970/retell/internal/repository"
	"github.com/usual2970/retell/internal/rest/common"
	"github.com/usual2970/retell/message"
	resourceposition "github.com/usual2970/retell/resource-position"
	"github.com/usual2970/retell/setting"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {

	// 公共路由
	RegisterCommon(e)

}
func RegisterCommon(e *echo.Echo) {
	group := e.Group("/common/v1")

	userAccontRepo := repository.NewUserAccountRepository()

	settingRepo := repository.NewSettingRepository()
	settingSvc := setting.NewService(settingRepo)

	messageRepo := repository.NewMessageRepository()
	messageSvc := message.NewService(messageRepo, userAccontRepo, settingRepo)

	codeRepo := repository.NewCodeRepository()
	codeSvc := code.NewService(codeRepo)

	accessTokenRepo := repository.NewAccessTokenRepository()
	accountRepo := repository.NewUserAccountRepository()

	authSvc := auth.NewService(codeSvc, messageSvc, accessTokenRepo, accountRepo)

	resourcePositionRepo := repository.NewResourcePositionRepository()
	resourcePositionSvc := resourceposition.NewService(resourcePositionRepo)

	common.NewAuthHandler(group, authSvc)

	common.NewCodeHandler(group, codeSvc)

	common.NewMessageHandler(group, messageSvc)

	common.NewResourceHandler(group, resourcePositionSvc)

	common.NewSettingHandler(group, settingSvc)

}
