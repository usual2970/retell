package domain

const (
	AuthLoginTypeCode     = "code"
	AuthLoginTypePassword = "password"
)

type AuthRole int8

const (
	_ AuthRole = iota
	AuthRolePurchaser
	AuthRoleDeveloper
	AuthRoleBusinessDeveloper
)

type AuthLoginReq struct {
	Identifier string `json:"identifier" validate:"required,email"`
	Type       string `json:"type" validate:"required,oneof=code password"`
	Data       string `json:"data" validate:"required"`
}

type AuthLoginResp struct {
	AccessToken string            `json:"accessToken"`
	ExpiresAt   int64             `json:"expiresIn"`
	Extend      map[string]string `json:"extend"`
}

type AuthRegisterReq struct {
	Email           string   `json:"email" validate:"required,email"`
	Password        string   `json:"password" validate:"required"`
	ConfirmPassword string   `json:"confirmPassword" validate:"required,eqfield=Password"`
	Code            string   `json:"code" validate:"required"`
	Role            AuthRole `json:"role"`
}

type AuthRegisterResp struct{}

type AuthCodeReq struct {
	Receiver string `json:"receiver"`
}

type AuthForgetPasswordReq struct {
	Email           string `json:"email" validate:"required,email"`
	Code            string `json:"code" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type AuthResetPasswordReq struct {
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type AuthResetHeadImgReq struct {
	Headimgurl string `json:"headimgurl" validate:"required,url"`
}
