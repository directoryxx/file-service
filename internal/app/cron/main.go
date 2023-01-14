package main

import (
	"context"
	"file/config"
	"file/infrastructure"
	"file/internal/repository"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	log.Println("[INFO] Starting FILE Cron")

	envSource := "SYSTEM"

	if os.Getenv("BYPASS_ENV_FILE") == "" {
		log.Println("[INFO] Load Config")
		config.LoadConfig()
		envSource = "FILE"
	}

	log.Println("[INFO] Loaded Config : " + envSource)

	log.Println("[INFO] Loading Database")
	dbSQL, err := infrastructure.Open()

	if err != nil {
		log.Fatalf("Could not initialize Database connection using sqlx %s", err)
	}

	defer dbSQL.Close()

	log.Println("[INFO] Loading Minio")
	minioClient := infrastructure.MinioConnection()

	log.Println("[INFO] Loading Repository")
	fileRepo := repository.NewFileRepository(dbSQL, minioClient)

	s := gocron.NewScheduler(time.Local)

	s.Every(1).Day().At("00:30").Do(func() {
		tempFile, _ := fileRepo.GetAllTemp()

		for _, s := range tempFile {
			fileRepo.DeleteByID(context.Background(), s.ID)
			fileRepo.Deletefile(context.Background(), s.Name)
		}
	})

	s.StartBlocking()
}
