package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth-server/domain"
	"github.com/mainflux/mainflux-auth-server/services"
)

func AddDeviceKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	apiKey := header[1]

	data, err := readPayload(r)
	if err != nil {
		sErr := err.(*domain.ServiceError)
		w.WriteHeader(sErr.Code)
		return
	}

	userId := ps.ByName("user_id")
	devId := ps.ByName("device_id")
	key, err := services.AddDeviceKey(userId, devId, apiKey, data)
	if err != nil {
		sErr := err.(*domain.ServiceError)
		w.WriteHeader(sErr.Code)
		return
	}

	rep := fmt.Sprintf(`{"key":"%s"}`, key)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rep))
}
