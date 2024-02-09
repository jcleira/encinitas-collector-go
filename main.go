package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/jcleira/encinitas-collector-go/config"
	agentServices "github.com/jcleira/encinitas-collector-go/internal/app/agent/services"
	agentHandlers "github.com/jcleira/encinitas-collector-go/internal/infra/http/agent/handlers"
	agentRepositoriesRedis "github.com/jcleira/encinitas-collector-go/internal/infra/repositories/agent/redis"
)

var errSignalQuit = errors.New("signal quit")

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("can't initialize zap logger: ", err)
	}

	var config config.Config
	err = envconfig.Process("", &config)
	if err != nil {
		logger.Fatal("can't process envconfig: ", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Pass,
		DB:       config.Redis.DB,
	})

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		eventCollector := agentServices.NewEventCollector(
			agentRepositoriesRedis.New(redisClient),
		)

		logger.Info("starting event collector")
		eventCollector.Collect(ctx)
		logger.Info("event collector stopped")

		return nil
	})

	g.Go(func() error {
		router := gin.Default()
		router.Use(cors.Default())

		router.POST("/agent/events",
			agentHandlers.NewEventsCreatorHandler(
				agentServices.NewEventPublisher(
					agentRepositoriesRedis.New(redisClient),
				),
			).Handle,
		)
		return router.Run(":3001")
	})

	g.Go(func() error {
		defer signal.Reset()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

		select {
		case sig := <-quit:
			return fmt.Errorf("signal received: %s: %w", sig.String(), errSignalQuit)
		case <-ctx.Done():
			return ctx.Err() //nolint:wrapcheck
		}
	})

	err = g.Wait()
	if !errors.Is(err, errSignalQuit) {
		logger.Fatal("error while waiting for errgroup: ", zap.Error(err))
	}
}
