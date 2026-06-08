package handler

import (
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
)

func SupabaseAuthMiddleware(authUsecase user.AuthUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			token := bearerToken(c.Request().Header.Get("Authorization"))
			if token == "" {
				return c.JSON(http.StatusUnauthorized, UnauthorizedError{
					Message: "missing token",
				})
			}

			userID, err := authUsecase.VerifyToken(ctx, token)
			if err != nil || userID == "" {
				return c.JSON(http.StatusUnauthorized, UnauthorizedError{
					Message: "invalid token",
				})
			}

			c.Set("user_id", userID)
			return next(c)
		}
	}
}
