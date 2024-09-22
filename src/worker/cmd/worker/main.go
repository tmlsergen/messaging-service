package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tmlsergen/messaging-service-worker/internal/app"
	"github.com/tmlsergen/messaging-service-worker/internal/message"
	msgSendler "github.com/tmlsergen/messaging-service-worker/pkg/msg-sendler"
)

const (
	cronConfigKey = "cron"
)

var mod int = 0

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	db := app.InitDB()
	redisClient := app.InitRedis()
	msgSendler := msgSendler.NewClient(os.Getenv("MSG_SENDLER_BASE_URL"), os.Getenv("MSG_SENDLER_AUTH_KEY"))

	messageRepository := message.NewRepository(db)
	messageService := message.NewService(messageRepository, redisClient, msgSendler, logger)

	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Error("failed to create scheduler", "error", err)
		panic(err)
	}

	ctx := context.Background()

	_, err = s.NewJob(
		gocron.DurationJob(2*time.Minute),
		gocron.NewTask(messageService.SendPendingMessages, ctx),
	)
	if err != nil {
		logger.Error("failed to create scheduler", "error", err)
		panic(err)
	}

	err = redisClient.Set(ctx, cronConfigKey, "start")
	if err != nil {
		logger.Error("failed to set redis key", "error", err)
		panic(err)
	}

	for {
		action, err := redisClient.Get(ctx, cronConfigKey)
		if err != nil {
			logger.Error("failed to get redis key", "error", err)
			panic(err)
		}

		if action == "start" && mod%2 == 0 {
			mod++
			s.Start()
			logger.Info("starting scheduler")
		}

		if action == "stop" && mod%2 == 1 {
			mod++
			s.StopJobs()
			logger.Info("stopping scheduler")
		}

		time.Sleep(2 * time.Second)
	}
}
