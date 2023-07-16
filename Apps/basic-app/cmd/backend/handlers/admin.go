package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hairizuanbinnoorazman/basic-app/logger"
)

type Status struct {
	Logger logger.Logger
}

func (h Status) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Start Status Handler")
	defer h.Logger.Info("End Status Handler")

	type resp struct {
		Status string `json:"status"`
	}
	raw, _ := json.Marshal(resp{Status: "ok"})
	w.WriteHeader(http.StatusOK)
	w.Write(raw)
}
