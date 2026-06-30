package handler

import (
	"errors"
	"net/http"

	groupusecase "github.com/CommitUpp/backend/api/application/usecase/group"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// グループ関連APIのHTTPリクエストを受け取るハンドラー
type GroupHandler struct {
	groupUsecase groupusecase.GroupUsecase
}

// グループ関連エンドポイントのハンドラーを生成
func NewGroupHandler(gu groupusecase.GroupUsecase) *GroupHandler {
	return &GroupHandler{
		groupUsecase: gu,
	}
}

// グループ作成リクエストを受け取り、呼び出し元を検証
func (h *GroupHandler) CreateGroup(c echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "invalid user",
		})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
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
		UserID:      userID,
		Name:        req.Name,
		AccessToken: accessToken,
	})
	if err != nil {
		if errors.Is(err, groupusecase.ErrInvalidUserID) || errors.Is(err, groupusecase.ErrInvalidAccessToken) {
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

func (h *GroupHandler) JoinGroup(c echo.Context, groupId openapi_types.UUID) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "invalid user",
		})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "invalid token",
		})
	}

	output, err := h.groupUsecase.JoinGroup(ctx, groupusecase.JoinGroupInput{
		GroupID:     groupId.String(),
		UserID:      userID,
		AccessToken: accessToken,
	})
	if err != nil {
		if errors.Is(err, groupusecase.ErrInvalidUserID) || errors.Is(err, groupusecase.ErrInvalidAccessToken) {
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{
				Message: "invalid user",
			})
		}

		if errors.Is(err, groupusecase.ErrInvalidGroupID) {
			return c.JSON(http.StatusBadRequest, BadRequestError{
				Message: "group_id is required",
			})
		}

		if errors.Is(err, groupusecase.ErrGroupNotFound) {
			return c.JSON(http.StatusNotFound, BadRequestError{
				Message: "group not found",
			})
		}

		if errors.Is(err, groupusecase.ErrGroupMemberAlreadyExists) {
			return c.JSON(http.StatusConflict, ConflictError{
				Message: "already joined group",
			})
		}

		if errors.Is(err, groupusecase.ErrGroupRepositoryNotConfigured) {
			return c.JSON(http.StatusInternalServerError, InternalServerError{
				Message: "group repository is not configured",
			})
		}

		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: err.Error(),
		})
	}

	res := JoinGroupResponse{
		Group: Group{
			Id:          output.ID,
			Name:        output.Name,
			MonthlyGoal: int32(output.MonthlyGoal),
			CreatedAt:   output.CreatedAt,
		},
	}

	return c.JSON(http.StatusOK, res)
}
