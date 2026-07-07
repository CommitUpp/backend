package postgres

import (
	"context"
	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/CommitUpp/backend/api/lib/tmdb"
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
		SELECT m.id, m.tmdb_id, m.title, m.poster_url
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

		movie.PosterURL = tmdb.BuildPosterURL(posterPath)

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
