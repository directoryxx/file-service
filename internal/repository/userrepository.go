package repository

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
)

// UserRepository represent the user's repository contract
type UserRepository interface {
	DeleteUUID(ctx context.Context, uuid string)
}

type UserRepositoryImpl struct {
	Redis *redis.Client
	Kafka *kafka.Producer
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewUserRepository(Redis *redis.Client, kafkaProducer *kafka.Producer) UserRepository {
	return &UserRepositoryImpl{
		// DB:    db,
		Redis: Redis,
		Kafka: kafkaProducer,
	}
}

func (m *UserRepositoryImpl) DeleteUUID(ctx context.Context, uuid string) {
	m.Redis.Del(ctx, uuid)
}
