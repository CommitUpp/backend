package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/CommitUpp/backend/api/interfaces/grpc/pb"
)

// SupabaseAuthMiddleware はEchoの共通ミドルウェア
func SupabaseAuthMiddleware(authClient pb.AuthServiceClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. フロント(axios等)から送られてきたヘッダーからトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" || authHeader == token {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "missing token"})
			}

			// apiコンテナ自身のgRPCクライアントを使い、ネットワーク越しにauthコンテナのgRPCサーバーへトークン検証を依頼
			res, err := authClient.VerifyToken(c.Request().Context(), &pb.VerifyTokenRequest{
				AccessToken: token,
			})
			if err != nil || res.UserId == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
			}

			// 3. 認証成功。返ってきた user_id をコンテキストにセットして、後ろの処理（chatやgroup等）へ流す
			c.Set("user_id", res.UserId)
			return next(c)
		}
	}
}
