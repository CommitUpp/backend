package router

import (
	"github.com/CommitUpp/backend/api/interfaces/handler"
	"github.com/labstack/echo/v4"
)

type RouterConfig struct {
	AuthHandler *handler.AuthHandler
}

func NewRouter(cfg RouterConfig) *echo.Echo {
	e := echo.New()

	v1 := e.Group("/api/v1")

	user := v1.Group("/user")

	user.POST("/callback", cfg.AuthHandler.LoginCallback)

	return e
}
