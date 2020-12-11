package main

import (
	"log"
	"os"

	"github.com/afilatov/mai-performance/service"
	"github.com/afilatov/mai-performance/service/handlers"
)

const (
	servicePort = 8080
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	metricsCollector, err := service.NewMetricsCollector()
	if err != nil {
		logger.Fatalf("unable to init metrics collector:", err)
	}

	reqHandler, err := handlers.NewRequestHandler(logger, metricsCollector)
	if err != nil {
		logger.Fatalf("unable to init request handler:", err)
	}

	logger.Printf("Starting HTTP server on port %d", servicePort)

	httpServer, err := service.NewHTTPServer(servicePort, reqHandler, metricsCollector)
	if err != nil {
		logger.Fatalf("unable to init HTTP server:", err)
	}

	if err := httpServer.Start(); err != nil {
		logger.Fatalf("unable to start HTTP server:", err)
	}
}
