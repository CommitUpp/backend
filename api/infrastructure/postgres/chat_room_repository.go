package postgres

import (
	"context"
	"database/sql"
	"errors"

	domainrepo "github.com/CommitUpp/backend/api/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *ChatRoomRepository) CreateChatRoom(
	ctx context.Context,
	input domainrepo.CreateChatRoomInput,
) (domainrepo.ChatRoom, error) {
	const query = `
		INSERT INTO chat_rooms (group_id, movie_id, movie_title, created_by)
		SELECT $1, m.id, m.title, $3
		FROM movies m
		WHERE m.id = $2
			AND EXISTS (
				SELECT 1
				FROM watch_statuses ws
				WHERE ws.user_id = $3
					AND ws.movie_id = m.id
					AND ws.status = 'watched'
			)
		RETURNING
			id::text,
			group_id::text,
			movie_id::text,
			movie_title,
			created_by::text,
			created_at
	`

	var chatRoom domainrepo.ChatRoom
	var movieID string
	err := r.db.QueryRow(
		ctx,
		query,
		input.GroupID,
		input.MovieID,
		input.CreatedBy,
	).Scan(
		&chatRoom.ID,
		&chatRoom.GroupID,
		&movieID,
		&chatRoom.MovieTitle,
		&chatRoom.CreatedBy,
		&chatRoom.CreatedAt,
	)
	if err == nil {
		chatRoom.MovieID = &movieID
		return chatRoom, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domainrepo.ChatRoom{}, domainrepo.ErrChatRoomAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domainrepo.ChatRoom{}, err
	}

	var movieExists bool
	if err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM movies WHERE id = $1)`,
		input.MovieID,
	).Scan(&movieExists); err != nil {
		return domainrepo.ChatRoom{}, err
	}
	if !movieExists {
		return domainrepo.ChatRoom{}, domainrepo.ErrMovieNotFound
	}

	return domainrepo.ChatRoom{}, domainrepo.ErrMovieNotWatched
}
