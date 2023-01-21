package worker

import (
	"context"
	"encoding/json"
	"file/infrastructure"
	"file/internal/controller"
	"file/internal/domain"
	"file/internal/repository"
	"file/internal/usecase"
	"fmt"
	"log"
	"sync"
	"time"
)

func RunWorker() {
	log.Println("[INFO] Starting File Service on port")

	log.Println("[INFO] Loading Kafka Consumer")
	kafkaConn, _ := infrastructure.ConnectKafkaConsumer()

	log.Println("[INFO] Loading Redis")
	redisConnect := infrastructure.OpenRedis()

	defer redisConnect.Close()

	log.Println("[INFO] Loading Kafka Producer")
	kafkaProducer, err := infrastructure.ConnectKafka()

	if err != nil {
		log.Fatalf("Could not initialize connection to kafka producer %s", err)
	}

	defer kafkaProducer.Close()

	log.Println("[INFO] Loading Repository")
	userRepo := repository.NewUserRepository(redisConnect, kafkaProducer)

	log.Println("[INFO] Loading Usecase")
	userUsecase := usecase.NewUserUseCase(userRepo)

	log.Println("[INFO] Loading Controller")
	userController := controller.NewUserController(userUsecase)

	var wg sync.WaitGroup

	wg.Add(2)

	// go func() {
	// 	consumerKafkaAuthLogout("auth-logout", kafkaConn, userController, &wg)
	// }()
	// go func() {
	// 	consumerKafkaAuthLogin("auth-login", kafkaConn, userController, &wg)
	// }()

	go consumerKafkaAuthLogin("auth-login", userController, &wg)
	go consumerKafkaAuthLogout("auth-logout", userController, &wg)

	// // go consumerKafkaAuth("auth", kafkaConn, userController, &wg)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	kafkaConn.Close()
	fmt.Println("Done!")
}

func consumerKafkaAuthLogout(topic string, userController controller.UserController, wg *sync.WaitGroup) {
	log.Println("[INFO] Loading Kafka Consumer Auth Logout")
	kafkaConsumer, _ := infrastructure.ConnectKafkaConsumer()
	kafkaConsumer.SubscribeTopics([]string{topic}, nil)

	defer wg.Done()

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := kafkaConsumer.ReadMessage(time.Second)
		if err == nil {
			action := domain.PublishAuthLogout{}
			json.Unmarshal(msg.Value, &action)
			if action.Action == "logout" {
				fmt.Println(action.Data.Uuid + " deleted")
				userController.DeleteUUID(context.Background(), action.Data.Uuid)
			}
		}
	}

	// kafkaConsumer.Close()
	defer kafkaConsumer.Close()
}

func consumerKafkaAuthLogin(topic string, userController controller.UserController, wg *sync.WaitGroup) {
	log.Println("[INFO] Loading Kafka Consumer Auth Login")
	kafkaConsumer, _ := infrastructure.ConnectKafkaConsumer()
	kafkaConsumer.SubscribeTopics([]string{topic}, nil)

	defer wg.Done()

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := kafkaConsumer.ReadMessage(time.Second)
		if err == nil {
			action := domain.PublishAuthLogin{}
			json.Unmarshal(msg.Value, &action)
			if action.Action == "login" {
				// fmt.Println(action)
				fmt.Println(action.Data.Uuid + " saved")
				userController.RememberUUID(context.Background(), &action)
			}
		}
	}

	defer kafkaConsumer.Close()
}

// func consumerKafkaFile(topic string, kafkaConsumer *kafka.Consumer, mailController controller.MailController, wg *sync.WaitGroup) {
// 	kafkaConsumer.SubscribeTopics([]string{topic}, nil)

// 	defer wg.Done()

// 	// A signal handler or similar could be used to set this to false to break the loop.
// 	run := true

// 	for run {
// 		msg, err := kafkaConsumer.ReadMessage(time.Second)
// 		if err == nil {
// 			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
// 			mailController.SendMail(context.TODO(), string(msg.Value))
// 		}
// 	}

// 	kafkaConsumer.Close()
// }
