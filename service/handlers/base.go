package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/afilatov/mai-performance/service"
)

type handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type wrappedHandler struct {
	endpoint string
	origH handler
	metrics service.MetricsCollectioner
}

func newWrappedHandler(endpoint string, h handler, metrics service.MetricsCollectioner) (*wrappedHandler, error) {
	if err := metrics.RegisterEndpoint(endpoint); err != nil {
		return nil, err
	}

	return &wrappedHandler{
		endpoint: endpoint,
		origH: h,
		metrics: metrics,
	}, nil
}

func (wh *wrappedHandler) Handle(w http.ResponseWriter, r *http.Request) {
	wh.metrics.CollectRequest(wh.endpoint)

	t := time.Now()
	wh.origH.Handle(w,r)
	reqDur := time.Since(t)

	wh.metrics.CollectRequestTime(wh.endpoint, float64(reqDur.Milliseconds()))
}

type requestHandler struct {
	logger *log.Logger
	metrics service.MetricsCollectioner
}

func NewRequestHandler(logger *log.Logger, metrics service.MetricsCollectioner) (*requestHandler, error) {
	return &requestHandler{
		logger: logger,
		metrics: metrics,
	}, nil
}

func (rh *requestHandler) RegisterHandlers(mux *http.ServeMux) error {
	if err := rh.registerHandler(mux, "", newIndexHandler()); err != nil {
		return err
	}

	if err := rh.registerHandler(mux, "_info", newInfoHandler()); err != nil {
		return err
	}

	if err := rh.registerHandler(mux, "delay", newDelayHandler()); err != nil {
		return err
	}

	return nil
}

func (rh *requestHandler) registerHandler(mux *http.ServeMux, endpoint string, h handler) error {
	wh, err := newWrappedHandler(endpoint, h, rh.metrics)
	if err != nil {
		return err
	}

	mux.HandleFunc("/" + endpoint, wh.Handle)

	return nil
}

func sendResponse(w http.ResponseWriter, code int, data interface{}) error {
	respBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal response json: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := fmt.Fprint(w, string(respBody)); err != nil {
		return fmt.Errorf("unable to send response: %w", err)
	}

	return nil
}

func sendErrorResponse(w http.ResponseWriter, code int, error string) error {
	resp := struct{
		Code int `json:"code"`
		Error string `json:"error"`
	}{Code: code, Error: error}

	respBody, err := json.Marshal(&resp)
	if err != nil {
		return fmt.Errorf("unable to marshal response json: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := fmt.Fprint(w, string(respBody)); err != nil {
		return fmt.Errorf("unable to send response: %w", err)
	}

	return nil
}
