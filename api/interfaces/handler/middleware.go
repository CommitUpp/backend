package handler

import (
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
)

func SupabaseAuthMiddleware(authUsecase user.AuthUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if authUsecase == nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Server Error"})
			}

			ctx := c.Request().Context()
			rawHeader := c.Request().Header.Get("Authorization")
			token := bearerToken(rawHeader)

			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "missing token"})
			}

			userID, err := authUsecase.VerifyToken(ctx, token)
			if err != nil || userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
			}

			c.Set("user_id", userID)
			c.Set("access_token", token)
			return next(c)
		}
	}
}
