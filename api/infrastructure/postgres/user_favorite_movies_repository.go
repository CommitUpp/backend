package postgres

import (
	"context"
	"github.com/CommitUpp/backend/api/domain/repository"
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
	return nil, nil
}

func (r *favoriteMovieRepository) Create(ctx context.Context, userID string, movieID string) error {
	return nil
}

func (r *favoriteMovieRepository) Delete(ctx context.Context, userID string, movieID string) error {
	return nil
}
