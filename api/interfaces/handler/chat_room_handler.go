package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/group"
	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ChatRoomHandler struct {
	chatRoomUsecase group.ChatRoomUsecase
}

func NewChatRoomHandler(u group.ChatRoomUsecase) *ChatRoomHandler {
	return &ChatRoomHandler{
		chatRoomUsecase: u,
	}
}

func (h *ChatRoomHandler) GetGroupChatRooms(c echo.Context, groupId string) error {
	ctx := c.Request().Context()

	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "ユーザーIDの解析に失敗しました"})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(
			http.StatusUnauthorized,
			UnauthorizedError{Message: "認証情報が見つかりません"},
		)
	}

	if groupId == "" {
		return c.JSON(http.StatusBadRequest, BadRequestError{Message: "グループIDが指定されていません"})
	}

	chatRooms, err := h.chatRoomUsecase.GetChatRooms(ctx, userIDStr, groupId, accessToken)
	if err != nil {
		switch {
		case errors.Is(err, group.ErrUserIDRequired), errors.Is(err, group.ErrAccessTokenRequired):
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
		case errors.Is(err, group.ErrGroupIDRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "グループIDが指定されていません"})
		case errors.Is(err, group.ErrForbidden):
			return c.JSON(http.StatusForbidden, ForbiddenError{Message: "このグループの閲覧権限がありません"})
		default:
			log.Printf("failed to get group chat rooms: user_id=%s group_id=%s err=%v", userIDStr, groupId, err)
			return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "チャットルーム一覧の取得に失敗しました"})
		}
	}

	res, err := chatRoomResponse(chatRooms)
	if err != nil {
		log.Printf("failed to build group chat rooms response: user_id=%s group_id=%s err=%v", userIDStr, groupId, err)
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "チャットルーム一覧の取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, GetGroupChatRoomsResponse{ChatRooms: res})
}

func chatRoomResponse(chatRooms []repository.ChatRoom) ([]ChatRoom, error) {
	res := make([]ChatRoom, 0, len(chatRooms))

	for _, chatRoom := range chatRooms {
		id, err := uuid.Parse(chatRoom.ID)
		if err != nil {
			return nil, err
		}

		groupID, err := uuid.Parse(chatRoom.GroupID)
		if err != nil {
			return nil, err
		}

		createdBy, err := uuid.Parse(chatRoom.CreatedBy)
		if err != nil {
			return nil, err
		}

		var movieID *openapi_types.UUID
		if chatRoom.MovieID != nil {
			parsedMovieID, err := uuid.Parse(*chatRoom.MovieID)
			if err != nil {
				return nil, err
			}
			movieID = &parsedMovieID
		}

		res = append(res, ChatRoom{
			Id:         id,
			GroupId:    groupID,
			MovieId:    movieID,
			MovieTitle: chatRoom.MovieTitle,
			CreatedBy:  createdBy,
			CreatedAt:  chatRoom.CreatedAt,
		})
	}

	return res, nil
}
