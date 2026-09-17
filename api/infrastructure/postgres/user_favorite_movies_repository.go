package postgres

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/CommitUpp/backend/api/lib/tmdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type favoriteMovieRepository struct {
	db *pgxpool.Pool
}

func NewFavoriteMovieRepository(db *pgxpool.Pool) repository.FavoriteMovieRepository {
	return &favoriteMovieRepository{
		db: db,
	}
}

func (r *favoriteMovieRepository) GetFavoriteMovies(ctx context.Context, userID string) ([]repository.Movie, error) {
	query := `
		SELECT m.id, m.tmdb_id, m.title, m.trailer_url
		FROM favorite_movies fm
		JOIN movies m ON fm.movie_id = m.id
		WHERE fm.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favoriteMovies []repository.Movie

	for rows.Next() {
		var movie repository.Movie
		var posterPath string

		if err := rows.Scan(&movie.MovieID, &movie.TMDBID, &movie.Title, &posterPath); err != nil {
			return nil, err
		}

		movie.TrailerURL = tmdb.BuildBackdropURL(posterPath)

		favoriteMovies = append(favoriteMovies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favoriteMovies, nil
}

func (r *favoriteMovieRepository) UpdateFavorites(
	ctx context.Context,
	userID string,
	addMovieIDs []string,
	deleteMovieIDs []string,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, userID); err != nil {
		return err
	}

	// 削除
	for _, movieID := range deleteMovieIDs {
		_, err := tx.Exec(
			ctx,
			`DELETE FROM favorite_movies
			 WHERE user_id = $1
			   AND movie_id = $2`,
			userID,
			movieID,
		)
		if err != nil {
			return err
		}
	}

	// 追加
	for _, movieID := range addMovieIDs {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO favorite_movies(user_id, movie_id)
			 VALUES($1, $2)`,
			userID,
			movieID,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *favoriteMovieRepository) AddFavoriteMovie(ctx context.Context, userID string, movieID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, userID); err != nil {
		return err
	}

	var movieExists bool
	if err := tx.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM movies WHERE id = $1)`,
		movieID,
	).Scan(&movieExists); err != nil {
		return err
	}
	if !movieExists {
		return repository.ErrMovieNotFound
	}

	var favoriteExists bool
	if err := tx.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1 FROM favorite_movies
			WHERE user_id = $1 AND movie_id = $2
		)`,
		userID,
		movieID,
	).Scan(&favoriteExists); err != nil {
		return err
	}
	if favoriteExists {
		return repository.ErrFavoriteMovieAlreadyExists
	}

	var favoriteCount int
	if err := tx.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM favorite_movies WHERE user_id = $1`,
		userID,
	).Scan(&favoriteCount); err != nil {
		return err
	}
	if favoriteCount >= 3 {
		return repository.ErrFavoriteMovieLimitExceeded
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO favorite_movies(user_id, movie_id) VALUES($1, $2)`,
		userID,
		movieID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return repository.ErrFavoriteMovieAlreadyExists
			case "23503":
				return repository.ErrMovieNotFound
			}
		}
		return err
	}

	return tx.Commit(ctx)
}

func (r *favoriteMovieRepository) DeleteFavoriteMovie(ctx context.Context, userID string, movieID string) error {
	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM favorite_movies
		 WHERE user_id = $1
		   AND movie_id = $2`,
		userID,
		movieID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrFavoriteMovieNotFound
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrFavoriteMovieNotFound
	}

	return nil
}
