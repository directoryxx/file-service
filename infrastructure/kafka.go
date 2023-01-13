package infrastructure

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConnectKafka() (kafkaConn *kafka.Producer, err error) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_URI")})
	if err != nil {
		panic(err)
	}

	return p, err
}

func ConnectKafkaConsumer() (kafkaConn *kafka.Consumer, err error) {
	hostname, _ := os.Hostname()

	p, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_URI"),
		"group.id":          hostname,
		"auto.offset.reset": "earliest",
	})

	return p, err
}
