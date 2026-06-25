package postgres

import (
	"context"

	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type movieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *movieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) GetMovies(ctx context.Context, keyword string) ([]repository.Movie, error) {

	var (
		rows pgx.Rows
		err  error
	)

	if keyword == "" {
		rows, err = r.db.Query(ctx, `
			SELECT id, tmdb_id, title, poster_url
			FROM movies
		`)
	} else {
		rows, err = r.db.Query(ctx, `
			SELECT id, tmdb_id, title, poster_url
			FROM movies
			WHERE title ILIKE '%' || $1 || '%'
		`, keyword)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []repository.Movie

	for rows.Next() {
		var m repository.Movie
		if err := rows.Scan(&m.MovieID, &m.TMDBID, &m.Title, &m.PosterURL); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return movies, nil
}
