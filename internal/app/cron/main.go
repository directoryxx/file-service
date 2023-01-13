package main

import (
	"file/config"
	"file/infrastructure"
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

	timezone, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(timezone)

	s.Every(1).Day().At("01:30").Do(func() {

	})

	s.StartBlocking()
}
