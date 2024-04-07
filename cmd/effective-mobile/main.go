package main

import (
	"github.com/Verce11o/effective-mobile-test/internal/app"
	"github.com/Verce11o/effective-mobile-test/internal/config"
	"github.com/Verce11o/effective-mobile-test/internal/lib/logger"
)

// @title EffectiveMobile Test API
// @description CRUD Cars service
// @version 1.0
// @host localhost:3010
// @BasePath /api/v1
func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg)

	app.Run(log, cfg)
}
