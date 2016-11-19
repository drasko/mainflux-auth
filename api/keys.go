package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func addKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	spec, err := readPayload(r)
	if err != nil {
		sErr := err.(*domain.AuthError)
		w.WriteHeader(sErr.Code)
		return
	}

	key, err := services.AddKey(header[1], spec)
	if err != nil {
		sErr := err.(*domain.AuthError)
		w.WriteHeader(sErr.Code)
		return
	}

	rep := fmt.Sprintf(`{"key":"%s"}`, key)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rep))
}

func updateKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: implement the key update
}
