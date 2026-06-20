package movie

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
)

type MovieStatusUsecase interface {
	Execute(ctx context.Context, movieID string, userID string, status string, accessToken string) error
	GetStatuses(ctx context.Context, userID string, status *string, accessToken string) ([]repository.MovieStatus, error)
}

type movieStatusUsecase struct {
	movieStatusRepo repository.MovieStatusRepository
}

func NewMovieStatusUsecase(repo repository.MovieStatusRepository) MovieStatusUsecase {
	return &movieStatusUsecase{
		movieStatusRepo: repo,
	}
}

// POST
func (u *movieStatusUsecase) Execute(ctx context.Context, movieID string, userID string, status string, accessToken string) error {
	if status != "watched" && status != "wanna_watch" {
		return errors.New("invalid status value")
	}

	err := u.movieStatusRepo.WatchStatus(ctx, movieID, userID, status, accessToken)
	if err != nil {
		return err
	}

	return nil
}

// GET
func (u *movieStatusUsecase) GetStatuses(ctx context.Context, userID string, status *string, accessToken string) ([]repository.MovieStatus, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if accessToken == "" {
		return nil, errors.New("access token is required")
	}

	if status != nil {
		if *status != "watched" && *status != "wanna_watch" {
			return nil, errors.New("invalid status filter value")
		}
	}

	statuses, err := u.movieStatusRepo.GetWatchStatuses(ctx, userID, status, accessToken)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
