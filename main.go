package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/jcleira/encinitas-collector-go/config"
	agentServices "github.com/jcleira/encinitas-collector-go/internal/app/agent/services"
	managerServices "github.com/jcleira/encinitas-collector-go/internal/app/manager/services"
	metricsServices "github.com/jcleira/encinitas-collector-go/internal/app/metrics/services"
	solanaServices "github.com/jcleira/encinitas-collector-go/internal/app/solana/services"
	agentHandlers "github.com/jcleira/encinitas-collector-go/internal/infra/http/agent/handlers"
	managerHandlers "github.com/jcleira/encinitas-collector-go/internal/infra/http/manager/handlers"
	metricsHandlers "github.com/jcleira/encinitas-collector-go/internal/infra/http/metrics/handlers"
	agentRepositoriesRedis "github.com/jcleira/encinitas-collector-go/internal/infra/repositories/agent/redis"
	managerRepositoriesSQL "github.com/jcleira/encinitas-collector-go/internal/infra/repositories/manager/sql"
	metricsRepositoriesInflux "github.com/jcleira/encinitas-collector-go/internal/infra/repositories/metrics/influx"
	solanaRepositoriesRedis "github.com/jcleira/encinitas-collector-go/internal/infra/repositories/solana/redis"
	solanaRepositoriesSQL "github.com/jcleira/encinitas-collector-go/internal/infra/repositories/solana/sql"
)

var errSignalQuit = errors.New("signal quit")

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var config config.Config
	if err := envconfig.Process("", &config); err != nil {
		slog.Error("can't process envconfig: ", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Pass,
		DB:       config.Redis.DB,
	})

	sqlx, err := sqlx.Connect("postgres", config.Postgres.URL())
	if err != nil {
		slog.Error("can't connect to postgres: ", err)
		os.Exit(1)
	}

	influx := influxdb.NewClient(config.InfluxDB.URL, config.InfluxDB.Token)

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
		transactionsCollector := solanaServices.NewTransactionsCollector(
			solanaRepositoriesSQL.New(sqlx),
			solanaRepositoriesRedis.New(redisClient),
		)

		logger.Info("starting transactions collector")
		transactionsCollector.Collect(ctx)
		logger.Info("transactions collector stopped")

		return nil
	})

	g.Go(func() error {
		ingester := metricsServices.NewIngester(
			solanaRepositoriesRedis.New(redisClient),
			agentRepositoriesRedis.New(redisClient),
			metricsRepositoriesInflux.New(
				influx,
				config.InfluxDB.TelegrafURL,
				metricsRepositoriesInflux.TransactionsBucket,
			),
			solanaRepositoriesSQL.New(sqlx),
		)

		logger.Info("starting ingester")
		ingester.Ingest(ctx)
		logger.Info("ingester stopped")

		return nil
	})

	g.Go(func() error {
		router := gin.Default()

		corsConfig := cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}

		router.Use(cors.New(corsConfig))

		router.POST("/agent/events",
			agentHandlers.NewEventsCreatorHandler(
				agentServices.NewEventPublisher(
					agentRepositoriesRedis.New(redisClient),
				),
			).Handle,
		)

		router.GET("/metrics/query",
			metricsHandlers.NewMetricsRetriever(
				metricsRepositoriesInflux.New(
					influx,
					config.InfluxDB.TelegrafURL,
					metricsRepositoriesInflux.TransactionsBucket,
				),
			).Handle,
		)

		router.GET("/metrics/programs/query",
			metricsHandlers.NewMetricsProgramRetrieverHandler(
				metricsRepositoriesInflux.New(
					influx,
					config.InfluxDB.TelegrafURL,
					metricsRepositoriesInflux.ProgramsBucket,
				),
			).Handle,
		)

		router.GET("/transactions/query",
			metricsHandlers.NewTransactionsRetriever(
				solanaRepositoriesSQL.New(sqlx),
			).Handle,
		)

		router.GET("/manager/programs",
			managerHandlers.NewProgramGetterHandler(
				managerServices.NewProgramGetter(
					managerRepositoriesSQL.New(sqlx),
				),
			).Handle,
		)

		router.POST("/manager/programs",
			managerHandlers.NewProgramsCreatorHandler(
				managerServices.NewProgramCreator(
					managerRepositoriesSQL.New(sqlx),
				),
			).Handle,
		)

		router.POST("/manager/emails",
			managerHandlers.NewEmailsCreatorHandler(
				managerServices.NewEmailCreator(
					managerRepositoriesSQL.New(sqlx),
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
		slog.Error("error while waiting for errgroup: ", err)
		os.Exit(1)
	}
}
