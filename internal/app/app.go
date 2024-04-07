package app

import (
	"context"
	"github.com/Verce11o/effective-mobile-test/cmd/carinfo"
	"github.com/Verce11o/effective-mobile-test/internal/config"
	"github.com/Verce11o/effective-mobile-test/internal/lib/postgres"
	"github.com/Verce11o/effective-mobile-test/internal/lib/redis"
	"github.com/Verce11o/effective-mobile-test/internal/lib/tracer"
	"github.com/Verce11o/effective-mobile-test/internal/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *zap.SugaredLogger, cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := postgres.New(ctx, cfg)
	defer db.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()

	trace := tracer.InitTracer(ctx, "cars service", cfg.Jaeger.Endpoint)

	srv := server.NewServer(log, db, redisClient, cfg, trace)

	go func() {
		if err := carinfo.Run(log, "localhost:3009"); err != nil {
			log.Fatalf("error while start carinfo server: %v", err)
		}
		log.Infof("External CarInfo api running on: %v", cfg.ExternalCarsApi.URL)
	}()

	go func() {
		if err := srv.Run(srv.InitRoutes()); err != nil {
			log.Fatalf("error while start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}

}
