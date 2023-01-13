package controller

import (
	"context"
	"file/internal/domain"
	"file/internal/usecase"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

// interface
type FileController interface {
	Upload(ec echo.Context) error
	// Get(ec echo.Context) error
	// Delete(ec echo.Context) error
}

// implement interface
type FileControllerImpl struct {
	FileUseCase usecase.FileUseCase
	Minio       *minio.Client
}

func NewFileController(fileUseCase usecase.FileUseCase, minio *minio.Client) FileController {
	return &FileControllerImpl{
		FileUseCase: fileUseCase,
		Minio:       minio,
	}
}

func (fc *FileControllerImpl) Upload(c echo.Context) error {
	// Convert Echo Context
	con := c.Request().Context()
	ctx, cancel := context.WithTimeout(con, 10000*time.Second)
	defer cancel()

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	user := c.Get("user").(domain.User)

	objectName := file.Filename
	fileBuffer := src
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size
	bucketName := os.Getenv("MINIO_BUCKET")

	info, errMinio := fc.Minio.PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

	fileModel := &domain.File{
		Uuid:   (uuid.New()).String(),
		Name:   info.Key,
		Url:    "http://" + os.Getenv("MINIO_ENDPOINT") + "/" + info.Bucket + "/" + info.Key,
		UserID: user.ID,
		IsTemp: 1,
	}

	fileModel, errInsert := fc.FileUseCase.Insert(ctx, fileModel)

	if errMinio != nil {
		return c.JSON(http.StatusInternalServerError, errMinio)
	}

	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, errInsert)
	}

	return c.JSON(http.StatusOK, fileModel)
}
