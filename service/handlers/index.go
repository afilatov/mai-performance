package handlers

import (
	"net/http"
)

type indexHandler struct {}

func newIndexHandler() *indexHandler {
	return &indexHandler{}
}

type indexResponse struct {
	Path string `json:"path"`
}

func (h *indexHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := sendResponse(w, http.StatusOK, indexResponse{Path: r.URL.Path}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

