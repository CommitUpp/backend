package user

import "context"

type AuthUseCase interface {
	Auth(ctx context.Context)
}
