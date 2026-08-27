package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CommitUpp/backend/api/application/usecase/group"
	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type stubChatMessageUsecase struct {
	input   group.CreateChatMessageInput
	message repository.ChatMessage
	err     error
}

func (s *stubChatMessageUsecase) CreateChatMessage(
	_ context.Context,
	input group.CreateChatMessageInput,
) (repository.ChatMessage, error) {
	s.input = input
	return s.message, s.err
}

func TestChatMessageHandlerCreateChatMessage(t *testing.T) {
	t.Parallel()

	messageID := uuid.New()
	roomID := uuid.New()
	userID := uuid.New()
	createdAt := time.Date(2026, time.August, 16, 12, 0, 0, 0, time.UTC)
	usecase := &stubChatMessageUsecase{message: repository.ChatMessage{
		ID:        messageID.String(),
		RoomID:    roomID.String(),
		UserID:    userID.String(),
		Content:   "message",
		CreatedAt: createdAt,
	}}
	handler := NewChatMessageHandler(usecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat-rooms/"+roomID.String()+"/messages", bytes.NewBufferString(`{"content":"message"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID.String())

	if err := handler.CreateChatMessage(c, roomID); err != nil {
		t.Fatalf("CreateChatMessage() unexpected error: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if usecase.input.UserID != userID.String() || usecase.input.RoomID != roomID.String() || usecase.input.Content != "message" {
		t.Fatalf("usecase input = %#v", usecase.input)
	}
}

func TestChatMessageHandlerCreateChatMessageRequiresAuthentication(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"content":"message"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewChatMessageHandler(&stubChatMessageUsecase{})
	if err := handler.CreateChatMessage(c, uuid.New()); err != nil {
		t.Fatalf("CreateChatMessage() unexpected error: %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
