package postgres

import (
	"context"
	"database/sql"

	domainrepo "github.com/CommitUpp/backend/api/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRoomRepository struct {
	db *pgxpool.Pool
}

func NewChatRoomRepository(db *pgxpool.Pool) domainrepo.ChatRoomRepository {
	return &ChatRoomRepository{
		db: db,
	}
}

func (r *ChatRoomRepository) GetChatRooms(
	ctx context.Context,
	groupID string,
) ([]domainrepo.ChatRoom, error) {
	const query = `
		SELECT
			id::text,
			group_id::text,
			movie_id::text,
			movie_title,
			created_by::text,
			created_at
		FROM chat_rooms
		WHERE group_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatRooms := make([]domainrepo.ChatRoom, 0)
	for rows.Next() {
		var (
			chatRoom domainrepo.ChatRoom
			movieID  sql.NullString
		)

		if err := rows.Scan(
			&chatRoom.ID,
			&chatRoom.GroupID,
			&movieID,
			&chatRoom.MovieTitle,
			&chatRoom.CreatedBy,
			&chatRoom.CreatedAt,
		); err != nil {
			return nil, err
		}

		if movieID.Valid {
			chatRoom.MovieID = &movieID.String
		}

		chatRooms = append(chatRooms, chatRoom)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatRooms, nil
}
