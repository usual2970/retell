package constant

import (
	"errors"
	"fmt"
)

var ErrRecordNotFound = errors.New("record not found")

var ErrParamWrongTel = errors.New("param wrong tel")

var ErrOPerateTooFast = errors.New("操作频繁")

var ErrChannelNotSupported = NewXError(3001, "消息渠道不支持")
var ErrCodeHasSend = NewXError(3002, "验证码已发送")

var ErrUnsupportedLoginType = NewXError(4000, "不支持的登录方式")
var ErrParamWrongCode = NewXError(4005, "验证码错误或已失效")
var ErrCodeExpired = NewXError(4006, "验证码错误或已失效")
var ErrAccountOrPasswordIncorrect = NewXError(4007, "账号或密码错误")
var ErrAccountNotExists = NewXError(4008, "账号不存在或已注销")
var ErrPasswordIncorrect = NewXError(4009, "密码错误")
var ErrDeveloperInfo = NewXError(4010, "请先完善开发者信息")
var ErrPurchaserInfo = NewXError(4011, "请先完善采购商信息")
var ErrProductNotFound = NewXError(4012, "产品不存在")
var ErrPurchaserBankAccountNotFound = NewXError(4013, "采购商银行账户不存在")
var ErrPasswordNotMatch = NewXError(4014, "两次输入的密码不一致")
var ErrAccountAlreadyExists = NewXError(4015, "账号已存在")
var ErrDeveloperChecked = NewXError(4016, "开发者信息已审核通过，无法修改")

var ErrNotLogin = NewXError(4999, "not login")

var ErrGameNotExists = NewXError(2000, "Game not exists")
var ErrGameVersionNotBuilded = NewXError(2001, "Game version not builded")
var ErrGameVersionNotExists = NewXError(2002, "Game version not exists")
var ErrGameVersionInitSqlURLNotExists = NewXError(2003, "请上传数据库初始化文件")

var ErrOrderNotExists = NewXError(5000, "实例不存在")
var ErrRuntimeSettingNotReady = NewXError(5001, "运行环境正在准备中，请 1 分钟后再试")
var ErrOrderNotStoped = NewXError(5002, "未停止的应用无法启动")
var ErrOrderNotDeployed = NewXError(5003, "未运行的应用无法停止")
var ErrOrderNotDeployedCannotRestart = NewXError(5004, "未运行的应用无法重启")
var ErrOrderDeploying = NewXError(5005, "应用正在部署中，请稍后再试")
var ErrOrderScheduled = NewXError(5006, "应用已经延时部署，请勿重复操作")
var ErrOrderDelayTime = NewXError(5007, "延时时间必须1分钟以上")
var ErrDeploymentNotScheduled = NewXError(5008, "非延时部署的应用")
var ErrOrderDeleteProdEnvironment = NewXError(5009, "生产环境的应用无法删除")
var ErrRuntimeConfig = NewXError(5010, "运行环境配置错误")
var ErrOrderTestEnvironmentLimit = NewXError(5011, "测试环境的应用数量已达上限")
var ErrOrderDeployed = NewXError(5012, "部署流程已完成，请勿重复操作")

var ErrAccessKeyNotExists = NewXError(6000, "AccessKey not exists")

var ErrResourcePositionNotExists = NewXError(7000, "ResourcePosition not exists")
var ErrResourceNotExists = NewXError(7001, "Resource not exists")

type XError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewXError(code int, msg string, params ...any) *XError {
	if len(params) > 0 {
		msg = fmt.Sprintf(msg, params...)
	}
	return &XError{code, msg}
}

func (e *XError) Error() string {
	return e.Msg
}

func (e *XError) GetCode() int {
	if e.Code == 0 {
		return 100
	}
	return e.Code
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}
