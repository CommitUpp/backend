package user

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
)

type UserMovieStatusUsecase interface {
	Execute(ctx context.Context, movieID string, userID string, status string, accessToken string) error
	GetStatuses(ctx context.Context, userID string, status *string, accessToken string) ([]repository.UserMovieStatus, error)
}

type userMovieStatusUsecase struct {
	userMovieStatusRepo repository.UserMovieStatusRepository
}

func NewUserMovieStatusUsecase(repo repository.UserMovieStatusRepository) UserMovieStatusUsecase {
	return &userMovieStatusUsecase{
		userMovieStatusRepo: repo,
	}
}

// POST
func (u *userMovieStatusUsecase) Execute(ctx context.Context, movieID string, userID string, status string, accessToken string) error {
	if status != "watched" && status != "wanna_watch" {
		return errors.New("invalid status value")
	}

	err := u.userMovieStatusRepo.WatchStatus(ctx, movieID, userID, status, accessToken)
	if err != nil {
		return err
	}

	return nil
}

// GET
func (u *userMovieStatusUsecase) GetStatuses(ctx context.Context, userID string, status *string, accessToken string) ([]repository.UserMovieStatus, error) {
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

	statuses, err := u.userMovieStatusRepo.GetWatchStatuses(ctx, userID, status, accessToken)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
