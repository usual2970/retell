package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (c *CustomValidator) Validate(i any) error {
	c.lazyInit()
	return c.validate.Struct(i)
}

func (c *CustomValidator) lazyInit() {
	c.once.Do(func() {
		c.validate = validator.New()
	})
}

func Register(e *echo.Echo) {
	e.Validator = &CustomValidator{}
}
