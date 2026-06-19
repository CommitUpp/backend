package handler

import (
	"errors"
	"net/http"

	groupusecase "github.com/CommitUpp/backend/api/application/usecase/group"
	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
)

// グループ関連APIのHTTPリクエストを受け取るハンドラー
type GroupHandler struct {
	authUsecase  user.AuthUsecase
	groupUsecase groupusecase.GroupUsecase
}

// グループ関連エンドポイントのハンドラーを生成
func NewGroupHandler(au user.AuthUsecase, gu groupusecase.GroupUsecase) *GroupHandler {
	return &GroupHandler{
		authUsecase:  au,
		groupUsecase: gu,
	}
}

// グループ作成リクエストを受け取り、呼び出し元を検証
func (h *GroupHandler) CreateGroup(c echo.Context) error {
	ctx := c.Request().Context()

	accessToken := bearerToken(c.Request().Header.Get("Authorization"))
	if accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "missing token",
		})
	}

	// Bearer tokenから認証済みユーザーIDを取得し、作成者としてusecaseへ渡します。
	userID, err := h.authUsecase.VerifyToken(ctx, accessToken)
	if err != nil || userID == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "invalid token",
		})
	}

	// OpenAPIから生成されたリクエスト型へJSON bodyをBindします。
	var req CreateGroupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError{
			Message: "invalid request body",
		})
	}

	// グループ名の正規化と月間目標の初期値設定はusecase層に任せます。
	output, err := h.groupUsecase.CreateGroup(ctx, groupusecase.CreateGroupInput{
		UserID: userID,
		Name:   req.Name,
	})
	if err != nil {
		if errors.Is(err, groupusecase.ErrInvalidUserID) {
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{
				Message: "invalid user",
			})
		}

		if errors.Is(err, groupusecase.ErrGroupRepositoryNotConfigured) {
			return c.JSON(http.StatusInternalServerError, InternalServerError{
				Message: "group repository is not configured",
			})
		}

		if errors.Is(err, groupusecase.ErrInvalidName) {
			return c.JSON(http.StatusBadRequest, BadRequestError{
				Message: "name is required",
			})
		}

		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: err.Error(),
		})
	}

	res := CreateGroupResponse{
		Group: Group{
			Id:          output.ID,
			Name:        output.Name,
			MonthlyGoal: int32(output.MonthlyGoal),
			CreatedAt:   output.CreatedAt,
		},
	}

	return c.JSON(http.StatusCreated, res)
}
