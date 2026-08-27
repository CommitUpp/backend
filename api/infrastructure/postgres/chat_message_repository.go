package postgres

import (
	"context"
	"errors"

	domainrepo "github.com/CommitUpp/backend/api/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatMessageRepository struct {
	db *pgxpool.Pool
}

func NewChatMessageRepository(db *pgxpool.Pool) domainrepo.ChatMessageRepository {
	return &ChatMessageRepository{db: db}
}

func (r *ChatMessageRepository) CreateChatMessage(
	ctx context.Context,
	input domainrepo.CreateChatMessageInput,
) (domainrepo.ChatMessage, error) {
	const query = `
		INSERT INTO chat_messages (room_id, user_id, content)
		SELECT cr.id, $2, $3
		FROM chat_rooms cr
		WHERE cr.id = $1
			AND EXISTS (
				SELECT 1
				FROM group_members gm
				WHERE gm.group_id = cr.group_id
					AND gm.user_id = $2
					AND gm.is_active = true
			)
		RETURNING id::text, room_id::text, user_id::text, content, created_at
	`

	var message domainrepo.ChatMessage
	err := r.db.QueryRow(ctx, query, input.RoomID, input.UserID, input.Content).Scan(
		&message.ID,
		&message.RoomID,
		&message.UserID,
		&message.Content,
		&message.CreatedAt,
	)
	if err == nil {
		return message, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domainrepo.ChatMessage{}, err
	}

	var roomExists bool
	if err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM chat_rooms WHERE id = $1)`,
		input.RoomID,
	).Scan(&roomExists); err != nil {
		return domainrepo.ChatMessage{}, err
	}
	if !roomExists {
		return domainrepo.ChatMessage{}, domainrepo.ErrChatRoomNotFound
	}

	return domainrepo.ChatMessage{}, domainrepo.ErrNotGroupMember
}
