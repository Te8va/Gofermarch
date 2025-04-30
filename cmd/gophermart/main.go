package main

import (
	"github.com/Te8va/Gofermarch/internal/server"
	"github.com/Te8va/Gofermarch/pkg/logger"
)

func main() {
	log := logger.Logger()
	log.Infoln("Starting application...")

	if err := server.Serve(); err != nil {
		log.Error("Error running server", "error", err)
	}
}
