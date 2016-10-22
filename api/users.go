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

func RegisterUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func AddUserKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	apiKey := header[1]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := domain.Payload{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(data.Scopes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, s := range data.Scopes {
		if s.Actions == "" || s.Id == "" || s.Resource == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	uid := ps.ByName("user_id")
	key, err := services.AddUserKey(uid, apiKey, data)
	if err != nil {
		serviceErr := err.(*domain.ServiceError)
		w.WriteHeader(serviceErr.Code)
		return
	}

	rep := fmt.Sprintf(`{"key":"%s"}`, key)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rep))
}
