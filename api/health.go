package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json := []byte(`{"status": "OK"}`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
