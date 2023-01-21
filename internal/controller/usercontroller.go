package controller

import (
	"context"
	"file/internal/domain"
	"file/internal/usecase"
)

// interface
type UserController interface {
	DeleteUUID(ctx context.Context, uuid string)
	RememberUUID(ctx context.Context, user *domain.PublishAuthLogin)
}

// implement interface
type UserControllerImpl struct {
	UserUsecase usecase.UserUseCase
	// Minio       *minio.Client
}

func NewUserController(userUsecase usecase.UserUseCase) UserController {
	return &UserControllerImpl{
		UserUsecase: userUsecase,
		// Minio:       minio,
	}
}

func (uc *UserControllerImpl) DeleteUUID(ctx context.Context, uuid string) {
	uc.UserUsecase.DeleteUUID(ctx, uuid)
}

func (uc *UserControllerImpl) RememberUUID(ctx context.Context, user *domain.PublishAuthLogin) {
	uc.UserUsecase.RememberUUID(ctx, user)
}
