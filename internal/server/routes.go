package server

import (
	"context"
	"encoding/gob"
	"fmt"
	_ "github.com/Verce11o/effective-mobile-test/docs"
	"github.com/Verce11o/effective-mobile-test/internal/cars/handler"
	"github.com/Verce11o/effective-mobile-test/internal/cars/repository"
	"github.com/Verce11o/effective-mobile-test/internal/cars/service"
	"github.com/Verce11o/effective-mobile-test/internal/config"
	"github.com/Verce11o/effective-mobile-test/internal/lib/communicator"
	"github.com/Verce11o/effective-mobile-test/internal/lib/tracer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	log        *zap.SugaredLogger
	db         *pgxpool.Pool
	redis      *redis.Client
	tracer     *tracer.JaegerTracing
	cfg        *config.Config
	httpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config, tracer *tracer.JaegerTracing) *Server {
	return &Server{log: log, db: db, redis: redis, cfg: cfg, tracer: tracer}
}

func (s *Server) Run(handler http.Handler) error {
	addr := fmt.Sprintf("%v:%v", s.cfg.Server.Host, s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr:         ":" + s.cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.log.Infof("Server running on: %v", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) InitRoutes() *gin.Engine {
	router := gin.New()

	gob.Register(map[string]interface{}{})

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	carRepo := repository.NewCarRepository(s.db, s.tracer.Tracer)
	carCache := repository.NewCarCacheRepository(s.redis, s.tracer.Tracer)

	carCommunicator := communicator.NewCommunicator(s.cfg.ExternalCarsApi.URL)
	carService := service.NewService(s.log, carRepo, carCache, s.tracer.Tracer, carCommunicator)
	carHandler := handler.NewHandler(s.log, carService, s.tracer.Tracer)

	api := router.Group("/api/v1")

	cars := api.Group("/cars")
	{
		cars.POST("", carHandler.CreateCar)
		cars.GET("", carHandler.GetCars)
		cars.PUT("", carHandler.UpdateCar)
		cars.DELETE("", carHandler.DeleteCar)
	}

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
