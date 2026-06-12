package router

import (
	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/CommitUpp/backend/api/interfaces/handler"
	"github.com/labstack/echo/v4"
)

type RouterConfig struct {
	AuthHandler        *handler.AuthHandler
	MovieStatusHandler *handler.MovieStatusHandler
	AuthUsecase user.AuthUsecase
}

func NewRouter(cfg RouterConfig) *echo.Echo {
	e := echo.New()

	e.Use(handler.CORSMiddleware())

	e.Use(handler.SupabaseAuthMiddleware(cfg.AuthUsecase))

	apiServer := handler.NewServer(cfg.AuthHandler, cfg.MovieStatusHandler)

	handler.RegisterHandlersWithBaseURL(e, apiServer, "/api/v1")

	return e
}
