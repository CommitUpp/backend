package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/CommitUpp/backend/api/interfaces/grpc/pb"
)

func SupabaseAuthMiddleware(authClient pb.AuthServiceClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			authHeader := c.Request().Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" || authHeader == token {
				return c.JSON(http.StatusUnauthorized, UnauthorizedError{
					Message: "missing token",
				})
			}

			res, err := authClient.VerifyToken(ctx, &pb.VerifyTokenRequest{
				AccessToken: token,
			})
			if err != nil || res.UserId == "" {
				return c.JSON(http.StatusUnauthorized, UnauthorizedError{
					Message: "invalid token",
				})
			}

			c.Set("user_id", res.UserId)
			return next(c)
		}
	}
}
