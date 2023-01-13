package worker

import (
	"context"
	"file/infrastructure"
	"file/internal/controller"
	"file/internal/repository"
	"file/internal/usecase"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func RunWorker() {
	log.Println("[INFO] Starting File Service on port")

	log.Println("[INFO] Loading Kafka Consumer")
	kafkaConn, err := infrastructure.ConnectKafkaConsumer()

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

	wg.Add(1)

	go consumerKafkaAuth("auth", kafkaConn, userController, &wg)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")
}

func consumerKafkaAuth(topic string, kafkaConsumer *kafka.Consumer, userController controller.UserController, wg *sync.WaitGroup) {
	kafkaConsumer.SubscribeTopics([]string{topic}, nil)

	defer wg.Done()

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := kafkaConsumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			userController.DeleteUUID(context.TODO(), string(msg.Value))
		}
	}

	kafkaConsumer.Close()
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
