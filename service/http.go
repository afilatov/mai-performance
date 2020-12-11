package service

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RequestHandler interface {
	RegisterHandlers(*http.ServeMux) error
}

const metricsEndpoint = "/metrics"

type HTTPServer struct {
	port int

	mux *http.ServeMux
	metrics MetricsCollectioner
}

func NewHTTPServer(port int, requestHandler RequestHandler, metrics MetricsCollectioner) (*HTTPServer, error) {
	serv := &HTTPServer{
		mux: http.NewServeMux(),
		port: port,
		metrics: metrics,
	}

	if err := requestHandler.RegisterHandlers(serv.mux); err != nil {
		return nil, err
	}

	serv.mux.Handle(metricsEndpoint, promhttp.Handler())

	return serv, nil
}

func (h *HTTPServer) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", h.port), h.mux)
}
