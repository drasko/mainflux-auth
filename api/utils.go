package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mainflux/mainflux-auth-server/domain"
)

func readPayload(r *http.Request) (domain.Payload, error) {
	data := domain.Payload{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return data, &domain.ServiceError{Code: http.StatusBadRequest}
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return data, &domain.ServiceError{Code: http.StatusBadRequest}
	}

	if len(data.Scopes) == 0 {
		return data, &domain.ServiceError{Code: http.StatusBadRequest}
	}

	for _, s := range data.Scopes {
		if s.Actions == "" || s.Id == "" || s.Resource == "" {
			return data, &domain.ServiceError{Code: http.StatusBadRequest}
		}
	}

	return data, nil
}
