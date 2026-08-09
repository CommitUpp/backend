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

func (h *ChatRoomHandler) GetGroupChatRooms(c echo.Context, groupId openapi_types.UUID) error {
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

	groupID := groupId.String()

	chatRooms, err := h.chatRoomUsecase.GetChatRooms(ctx, userIDStr, groupID, accessToken)
	if err != nil {
		switch {
		case errors.Is(err, group.ErrUserIDRequired), errors.Is(err, group.ErrAccessTokenRequired):
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
		case errors.Is(err, group.ErrGroupIDRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "グループIDが指定されていません"})
		case errors.Is(err, group.ErrForbidden):
			return c.JSON(http.StatusForbidden, ForbiddenError{Message: "このグループの閲覧権限がありません"})
		default:
			log.Printf("failed to get group chat rooms: user_id=%s group_id=%s err=%v", userIDStr, groupID, err)
			return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "チャットルーム一覧の取得に失敗しました"})
		}
	}

	res, err := chatRoomResponse(chatRooms)
	if err != nil {
		log.Printf("failed to build group chat rooms response: user_id=%s group_id=%s err=%v", userIDStr, groupID, err)
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "チャットルーム一覧の取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, GetGroupChatRoomsResponse{ChatRooms: res})
}

func (h *ChatRoomHandler) CreateGroupChatRoom(c echo.Context, groupId openapi_types.UUID) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	var req CreateChatRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError{Message: "リクエストが不正です"})
	}
	if req.MovieId == uuid.Nil {
		return c.JSON(http.StatusBadRequest, BadRequestError{Message: "映画IDが指定されていません"})
	}

	chatRoom, err := h.chatRoomUsecase.CreateChatRoom(ctx, group.CreateChatRoomInput{
		UserID:      userID,
		GroupID:     groupId.String(),
		MovieID:     req.MovieId.String(),
		AccessToken: accessToken,
	})
	if err != nil {
		switch {
		case errors.Is(err, group.ErrUserIDRequired), errors.Is(err, group.ErrAccessTokenRequired):
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
		case errors.Is(err, group.ErrGroupIDRequired), errors.Is(err, group.ErrMovieIDRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "必須項目が指定されていません"})
		case errors.Is(err, group.ErrForbidden):
			return c.JSON(http.StatusForbidden, ForbiddenError{Message: "このグループの操作権限がありません"})
		case errors.Is(err, group.ErrMovieNotWatched):
			return c.JSON(http.StatusForbidden, ForbiddenError{Message: "視聴済みの映画のみチャットルームを作成できます"})
		case errors.Is(err, group.ErrMovieNotFound):
			return c.JSON(http.StatusNotFound, NotFoundError{Message: "映画が見つかりません"})
		case errors.Is(err, group.ErrChatRoomAlreadyExists):
			return c.JSON(http.StatusConflict, ConflictError{Message: "この映画のチャットルームは既に存在します"})
		default:
			log.Printf(
				"failed to create group chat room: user_id=%s group_id=%s movie_id=%s err=%v",
				userID,
				groupId.String(),
				req.MovieId.String(),
				err,
			)
			return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "チャットルームの作成に失敗しました"})
		}
	}

	res, err := chatRoomToResponse(chatRoom)
	if err != nil {
		log.Printf("failed to build chat room response: chat_room_id=%s err=%v", chatRoom.ID, err)
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "チャットルームの作成に失敗しました"})
	}

	return c.JSON(http.StatusCreated, CreateChatRoomResponse{ChatRoom: res})
}

func chatRoomResponse(chatRooms []repository.ChatRoom) ([]ChatRoom, error) {
	res := make([]ChatRoom, 0, len(chatRooms))

	for _, chatRoom := range chatRooms {
		item, err := chatRoomToResponse(chatRoom)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil
}

func chatRoomToResponse(chatRoom repository.ChatRoom) (ChatRoom, error) {
	id, err := uuid.Parse(chatRoom.ID)
	if err != nil {
		return ChatRoom{}, err
	}

	groupID, err := uuid.Parse(chatRoom.GroupID)
	if err != nil {
		return ChatRoom{}, err
	}

	createdBy, err := uuid.Parse(chatRoom.CreatedBy)
	if err != nil {
		return ChatRoom{}, err
	}

	var movieID *openapi_types.UUID
	if chatRoom.MovieID != nil {
		parsedMovieID, err := uuid.Parse(*chatRoom.MovieID)
		if err != nil {
			return ChatRoom{}, err
		}
		movieID = &parsedMovieID
	}

	return ChatRoom{
		Id:         id,
		GroupId:    groupID,
		MovieId:    movieID,
		MovieTitle: chatRoom.MovieTitle,
		CreatedBy:  createdBy,
		CreatedAt:  chatRoom.CreatedAt,
	}, nil
}
