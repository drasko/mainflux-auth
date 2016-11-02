package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func checkCredentials(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &domain.AccessRequest{}
	if err = json.Unmarshal(body, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = services.CheckPermissions(data)
	if err != nil {
		sErr := err.(*domain.ServiceError)
		w.WriteHeader(sErr.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
}
