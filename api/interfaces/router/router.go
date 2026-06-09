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

	e.Use(handler.CORSMiddleware())
	handler.RegisterHandlersWithBaseURL(e, cfg.AuthHandler, "/api/v1")

	return e
}
