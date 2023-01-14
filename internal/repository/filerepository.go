package repository

import (
	"context"
	"database/sql"
	"file/internal/domain"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
)

// UserRepository represent the user's repository contract
type FileRepository interface {
	Insert(ctx context.Context, file *domain.File) (*domain.File, error)
	GetOneByID(ctx context.Context, id int) (*domain.File, error)
	GetOneByName(ctx context.Context, name string) (*domain.File, error)
	GetAllTemp() ([]domain.File, error)
	Deletefile(ctx context.Context, name string)
	DeleteByID(ctx context.Context, id int) error
}

type FileRepositoryImpl struct {
	DB    *sql.DB
	Minio *minio.Client
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewFileRepository(db *sql.DB, minio *minio.Client) FileRepository {
	return &FileRepositoryImpl{
		DB:    db,
		Minio: minio,
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

func (fr *FileRepositoryImpl) GetAllTemp() ([]domain.File, error) {
	rows, err := fr.DB.Query("SELECT id, uuid, name, url, user_id, is_temp FROM files WHERE is_temp=1")

	if err != nil {
		return nil, err
	}

	files := []domain.File{}

	for rows.Next() {
		file := domain.File{}
		// create an instance of `Bird` and write the result of the current row into it
		if err := rows.Scan(&file.ID, &file.Uuid, &file.Name, &file.Url, &file.UserID, &file.IsTemp); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		// append the current instance to the slice of birds
		files = append(files, file)
	}

	return files, nil
}

func (fr *FileRepositoryImpl) Deletefile(ctx context.Context, name string) {
	bucketName := os.Getenv("MINIO_BUCKET")
	err := fr.Minio.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions{})
	fmt.Println(err)
}

func (fr *FileRepositoryImpl) DeleteByID(ctx context.Context, id int) (err error) {
	stmt, err := fr.DB.PrepareContext(ctx, "DELETE FROM files WHERE id=$1")
	stmt.Exec(id)

	return err
}
