package handlers

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	maxDelayMs = 2000
	defaultCode = 200
)

type delayHandler struct {}

func newDelayHandler() *delayHandler {
	return &delayHandler{}
}

type delayRequest struct {
	DelayMs int64 `json:"delay_ms"`
	Code int `json:"code"`
}

type delayResponse struct {
	DelayMs int64 `json:"delay_ms"`
	Code int `json:"code"`
}

func (h *delayHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusBadRequest, "Bad request method")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Unable to read request body")
		return
	}

	var req delayRequest
	if err := json.Unmarshal(body, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Bad request body")
		return
	}

	if req.Code == 0 {
		req.Code = defaultCode
	}

	if req.DelayMs == -1 {
		req.DelayMs = rand.Int63n(maxDelayMs)
	}

	time.Sleep(time.Duration(req.DelayMs) * time.Millisecond)

	if err := sendResponse(w, req.Code, delayResponse{Code: req.Code, DelayMs: req.DelayMs}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

