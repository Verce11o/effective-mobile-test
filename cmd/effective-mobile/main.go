package main

import (
	"github.com/Verce11o/effective-mobile-test/internal/app"
	"github.com/Verce11o/effective-mobile-test/internal/config"
	"github.com/Verce11o/effective-mobile-test/internal/lib/logger"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg)

	app.Run(log, cfg)
}
