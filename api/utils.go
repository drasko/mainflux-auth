package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mainflux/mainflux-auth/domain"
)

func readPayload(r *http.Request) (domain.KeySpec, error) {
	data := domain.KeySpec{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return data, &domain.AuthError{Code: http.StatusBadRequest}
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return data, &domain.AuthError{Code: http.StatusBadRequest}
	}

	if valid := data.Validate(); !valid {
		return data, &domain.AuthError{Code: http.StatusBadRequest}
	}

	return data, nil
}
