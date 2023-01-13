package repository

import (
	"context"
	"database/sql"
	"file/internal/domain"
	"time"
)

// UserRepository represent the user's repository contract
type FileRepository interface {
	Insert(ctx context.Context, file *domain.File) (*domain.File, error)
	GetOneByID(ctx context.Context, id int) (*domain.File, error)
	GetOneByName(ctx context.Context, name string) (*domain.File, error)
}

type FileRepositoryImpl struct {
	DB *sql.DB
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewFileRepository(db *sql.DB) FileRepository {
	return &FileRepositoryImpl{
		DB: db,
	}
}

func (fr *FileRepositoryImpl) Insert(ctx context.Context, file *domain.File) (fileDomain *domain.File, err error) {
	stmt := `insert into files (uuid, name, url, user_id, is_temp, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	fileOld, _ := fr.GetOneByName(ctx, file.Name)

	if fileOld != nil {
		return fileOld, nil
	}

	var newID int

	err = fr.DB.QueryRowContext(ctx, stmt,
		file.Uuid,
		file.Name,
		file.Url,
		file.UserID,
		file.IsTemp,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return nil, err
	}

	file, err = fr.GetOneByID(ctx, newID)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fr *FileRepositoryImpl) GetOneByID(context context.Context, id int) (res *domain.File, err error) {
	stmt, err := fr.DB.PrepareContext(context, "SELECT id, uuid, name, url, user_id, is_temp FROM files WHERE id=$1")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(context, id)
	var file domain.File

	err = row.Scan(
		&file.ID,
		&file.Uuid,
		&file.Name,
		&file.Url,
		&file.UserID,
		&file.IsTemp,
	)

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (fr *FileRepositoryImpl) GetOneByName(context context.Context, name string) (res *domain.File, err error) {
	stmt, err := fr.DB.PrepareContext(context, "SELECT id, uuid, name, url, user_id, is_temp FROM files WHERE name=$1")

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(context, name)

	var file domain.File

	err = row.Scan(
		&file.ID,
		&file.Uuid,
		&file.Name,
		&file.Url,
		&file.UserID,
		&file.IsTemp,
	)

	if err != nil {
		return nil, err
	}

	return &file, nil
}
