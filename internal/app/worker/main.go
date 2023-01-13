package main

import (
	"file/config"
	"file/delivery/worker"
	"log"
	"os"
)

func main() {
	envSource := "SYSTEM"

	if os.Getenv("BYPASS_ENV_FILE") == "" {
		log.Println("[INFO] Load Config")
		config.LoadConfig()
		envSource = "FILE"
	}

	log.Println("[INFO] Loaded Config : " + envSource)

	worker.RunWorker()
}
