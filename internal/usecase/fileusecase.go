package usecase

import (
	"context"
	"file/internal/domain"
	"file/internal/repository"
)

// UserUseCase represent the user's usecase contract
type FileUseCase interface {
	Insert(ctx context.Context, file *domain.File) (fileDomain *domain.File, err error)
}

type FileUseCaseImpl struct {
	FileRepo repository.FileRepository
}

// // NewMysqlAuthorRepository will create an implementation of author.Repository
func NewFileUseCase(FileRepo repository.FileRepository) FileUseCase {
	return &FileUseCaseImpl{
		FileRepo: FileRepo,
	}
}

func (fuc *FileUseCaseImpl) Insert(ctx context.Context, file *domain.File) (fileDomain *domain.File, err error) {
	file, err = fuc.FileRepo.Insert(ctx, file)

	return file, err
}
