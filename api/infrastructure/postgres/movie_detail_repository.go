package postgres

import (
	"context"

	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type movieDetailRepository struct {
	db *pgxpool.Pool
}

func NewMovieDetailRepository(db *pgxpool.Pool) *movieDetailRepository {
	return &movieDetailRepository{db: db}
}

func (r *movieDetailRepository) GetMovieDetail(ctx context.Context, movieID string) (*repository.MovieDetail, error) {
	var movieDetail repository.MovieDetail
	var posterPath string

	// 映画本体の情報を取得
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			tmdb_id,
			title,
			poster_url,
			trailer_url,
			overview,
			release_date
		FROM movies
		WHERE id = $1
	`, movieID).Scan(
		&movieDetail.MovieID,
		&movieDetail.TMDBID,
		&movieDetail.Title,
		&posterPath,
		&movieDetail.TrailerURL,
		&movieDetail.Overview,
		&movieDetail.ReleaseDate,
	)
	if err != nil {
		return nil, err
	}

	// poster_url をフルパスに変換する場合はここで組み立てる
	movieDetail.PosterURL = posterPath

	// watched_user を取得
	watchedUsers, err := r.getWatchedUsers(ctx, movieID)
	if err != nil {
		return nil, err
	}
	movieDetail.WatchedUser = watchedUsers

	// streaming_services を取得
	streamingServices, err := r.getStreamingServices(ctx, movieID)
	if err != nil {
		return nil, err
	}
	movieDetail.StreamingServices = streamingServices

	return &movieDetail, nil
}

func (r *movieDetailRepository) getWatchedUsers(ctx context.Context, movieID string) ([]repository.WatchedUser, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			u.id,
			u.avatar_url
		FROM watch_statuses ws
		INNER JOIN users u
			ON u.id = ws.user_id
		WHERE ws.movie_id = $1
			AND ws.status = 'watched'
		ORDER BY u.id
	`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	watchedUsers := make([]repository.WatchedUser, 0)
	for rows.Next() {
		var user repository.WatchedUser
		if err := rows.Scan(
			&user.UserID,
			&user.AvatarURL,
		); err != nil {
			return nil, err
		}
		watchedUsers = append(watchedUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return watchedUsers, nil
}

func (r *movieDetailRepository) getStreamingServices(ctx context.Context, movieID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
				ss.name
		FROM movie_streamings ms
		INNER JOIN streaming_services ss
				ON ss.id = ms.service_id
		WHERE ms.movie_id = $1
		ORDER BY ss.name
	`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	streamingServices := make([]string, 0)
	for rows.Next() {
		var serviceName string

		if err := rows.Scan(&serviceName); err != nil {
			return nil, err
		}

		streamingServices = append(streamingServices, serviceName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return streamingServices, nil
}
