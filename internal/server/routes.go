package server

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/cars/handler"
	"github.com/Verce11o/effective-mobile-test/internal/cars/repository"
	"github.com/Verce11o/effective-mobile-test/internal/cars/service"
	"github.com/Verce11o/effective-mobile-test/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	log        *zap.SugaredLogger
	db         *pgxpool.Pool
	cfg        *config.Config
	httpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, db *pgxpool.Pool, cfg *config.Config) *Server {
	return &Server{log: log, db: db, cfg: cfg}
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
		AllowOrigins:     []string{"https://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	carRepo := repository.NewCarRepository(s.db)
	carService := service.NewService(s.log, carRepo, s.cfg.ExternalCarsApi.URL)
	carHandler := handler.NewHandler(s.log, carService)

	api := router.Group("/api/v1")

	cars := api.Group("/cars")
	{
		cars.POST("", carHandler.CreateCar)
		cars.GET("", carHandler.GetCars)
	}

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
