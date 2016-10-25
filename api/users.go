package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth-server/domain"
	"github.com/mainflux/mainflux-auth-server/services"
)

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &userReq{}
	if err := json.Unmarshal(body, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := services.RegisterUser(data.Username, data.Password)
	if err != nil {
		serviceErr := err.(*domain.ServiceError)
		w.WriteHeader(serviceErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func addUserKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	key, err := services.AddUserKey(userId, apiKey, data)
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
