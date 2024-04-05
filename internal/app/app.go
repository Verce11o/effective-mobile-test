package app

import (
	"context"
	"github.com/Verce11o/effective-mobile-test/internal/config"
	"github.com/Verce11o/effective-mobile-test/internal/lib/postgres"
	"github.com/Verce11o/effective-mobile-test/internal/lib/redis"
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

	srv := server.NewServer(log, db, cfg)

	go func() {
		if err := srv.Run(srv.InitRoutes()); err != nil {
			log.Fatalf("Error while start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}

}
