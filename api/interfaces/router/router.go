package router

import (
	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/CommitUpp/backend/api/interfaces/handler"
	"github.com/labstack/echo/v4"
)

type RouterConfig struct {
	AuthHandler        *handler.AuthHandler
	MovieStatusHandler *handler.MovieStatusHandler
	AuthUsecase        user.AuthUsecase
}

func NewRouter(cfg RouterConfig) *echo.Echo {
	e := echo.New()

	e.Use(handler.CORSMiddleware())

	apiServer := handler.NewServer(cfg.AuthHandler, cfg.MovieStatusHandler)

	authMiddleware := handler.SupabaseAuthMiddleware(cfg.AuthUsecase)
	handler.RegisterHandlersWithOptions(e, apiServer, handler.RegisterHandlersOptions{
		BaseURL: "/api/v1",
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"loginCallback":     {authMiddleware},
			"logout":            {authMiddleware},
			"updateMovieStatus": {authMiddleware},
		},
	})

	return e
}
