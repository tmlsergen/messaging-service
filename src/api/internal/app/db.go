package app

import (
	"os"

	"github.com/tmlsergen/messaging-service-api/pkg/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func InitRedis() *redis.RedisClient {
	dsn := os.Getenv("REDIS_DSN")

	return redis.NewRedisClient(dsn)
}
