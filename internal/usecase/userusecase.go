package usecase

import (
	"context"
	"file/internal/repository"
)

// UserUseCase represent the user's usecase contract
type UserUseCase interface {
	DeleteUUID(ctx context.Context, uuid string)
}

type UserUseCaseImpl struct {
	UserRepo repository.UserRepository
}

// // NewMysqlAuthorRepository will create an implementation of author.Repository
func NewUserUseCase(UserRepo repository.UserRepository) UserUseCase {
	return &UserUseCaseImpl{
		UserRepo: UserRepo,
	}
}

func (uc *UserUseCaseImpl) DeleteUUID(ctx context.Context, uuid string) {
	uc.UserRepo.DeleteUUID(ctx, uuid)
}
