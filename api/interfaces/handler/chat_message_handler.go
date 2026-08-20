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

type ChatMessageHandler struct {
	chatMessageUsecase group.ChatMessageUsecase
}

func NewChatMessageHandler(u group.ChatMessageUsecase) *ChatMessageHandler {
	return &ChatMessageHandler{chatMessageUsecase: u}
}

func (h *ChatMessageHandler) CreateChatMessage(c echo.Context, roomId openapi_types.UUID) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	var req CreateChatMessageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError{Message: "リクエストが不正です"})
	}

	message, err := h.chatMessageUsecase.CreateChatMessage(ctx, group.CreateChatMessageInput{
		RoomID:  roomId.String(),
		UserID:  userID,
		Content: req.Content,
	})
	if err != nil {
		switch {
		case errors.Is(err, group.ErrUserIDRequired):
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
		case errors.Is(err, group.ErrChatMessageRoomIDRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "チャットルームIDが指定されていません"})
		case errors.Is(err, group.ErrChatMessageContentRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "メッセージを入力してください"})
		case errors.Is(err, group.ErrChatMessageContentTooLong):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "メッセージは1000文字以内で入力してください"})
		case errors.Is(err, group.ErrNotGroupMember):
			return c.JSON(http.StatusForbidden, ForbiddenError{Message: "このチャットルームへの投稿権限がありません"})
		case errors.Is(err, group.ErrChatRoomNotFound):
			return c.JSON(http.StatusNotFound, NotFoundError{Message: "チャットルームが見つかりません"})
		default:
			log.Printf("failed to create chat message: user_id=%s room_id=%s err=%v", userID, roomId.String(), err)
			return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "メッセージの送信に失敗しました"})
		}
	}

	res, err := chatMessageToResponse(message)
	if err != nil {
		log.Printf("failed to build chat message response: message_id=%s err=%v", message.ID, err)
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "メッセージの送信に失敗しました"})
	}

	return c.JSON(http.StatusCreated, CreateChatMessageResponse{Message: res})
}

func chatMessageToResponse(message repository.ChatMessage) (ChatMessage, error) {
	id, err := uuid.Parse(message.ID)
	if err != nil {
		return ChatMessage{}, err
	}
	roomID, err := uuid.Parse(message.RoomID)
	if err != nil {
		return ChatMessage{}, err
	}
	userID, err := uuid.Parse(message.UserID)
	if err != nil {
		return ChatMessage{}, err
	}

	return ChatMessage{
		Id:        id,
		RoomId:    roomID,
		UserId:    userID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}, nil
}
