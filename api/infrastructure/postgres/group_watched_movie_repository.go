package postgres

import (
	"context"

	domainrepo "github.com/CommitUpp/backend/api/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupWatchedMovieRepository struct {
	db *pgxpool.Pool
}

func NewGroupWatchedMovieRepository(db *pgxpool.Pool) domainrepo.GroupWatchedMovieRepository {
	return &GroupWatchedMovieRepository{
		db: db,
	}
}

func (r *GroupWatchedMovieRepository) GetWatchedMovies(
	ctx context.Context,
	groupID string,
	excludeUserID string,
) ([]domainrepo.GroupWatchedMovieRow, error) {
	const query = `
		SELECT
			gm.group_id,
			m.id AS movie_id,
			m.title,
			m.poster_url,
			u.id AS user_id,
			u.avatar_url
		FROM group_members gm
		INNER JOIN watch_statuses ms
			ON ms.user_id = gm.user_id
		INNER JOIN movies m
			ON m.id = ms.movie_id
		INNER JOIN users u
			ON u.id = gm.user_id
		WHERE gm.group_id = $1
			AND ms.status = 'watched'
			AND gm.user_id <> $2
			AND NOT EXISTS (
				SELECT 1
				FROM watch_statuses self_ws
				WHERE self_ws.user_id = $2
					AND self_ws.movie_id = ms.movie_id
					AND self_ws.status = 'watched'
			)
		ORDER BY m.title, u.id
	`

	rows, err := r.db.Query(ctx, query, groupID, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	watchedMovies := make([]domainrepo.GroupWatchedMovieRow, 0)
	for rows.Next() {
		var row domainrepo.GroupWatchedMovieRow
		if err := rows.Scan(
			&row.GroupID,
			&row.MovieID,
			&row.Title,
			&row.PosterURL,
			&row.UserID,
			&row.AvatarURL,
		); err != nil {
			return nil, err
		}
		watchedMovies = append(watchedMovies, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return watchedMovies, nil
}
