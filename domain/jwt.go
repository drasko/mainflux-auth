package domain

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

const issuer string = "mainflux"

var secretKey string = "mainflux-api-key"

type KeyData struct {
	AccessSpec
	jwt.StandardClaims
}

// SetSecretKey sets the secret key that will be used for decoding and encoding
// generated tokens. If not invoked, a default key will be used instead.
func SetSecretKey(key string) {
	secretKey = key
}

// CreateKey creates new JSON Web Token containing provided subject and
// access specification.
func CreateKey(subject string, access *AccessSpec) (string, error) {
	data := KeyData{
		*access,
		jwt.StandardClaims{
			Issuer:  issuer,
			Subject: subject,
		},
	}

	key := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	raw, err := key.SignedString([]byte(secretKey))
	if err != nil {
		return "", &ServiceError{http.StatusInternalServerError, err.Error()}
	}

	return raw, nil
}

// ContentOf extracts the key data given its string representation.
func ContentOf(key string) (KeyData, error) {
	data := KeyData{}

	token, err := jwt.ParseWithClaims(
		key,
		&data,
		func(val *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return data, &ServiceError{http.StatusBadRequest, err.Error()}
	}

	if token.Valid {
		return data, nil
	}

	return data, &ServiceError{http.StatusForbidden, err.Error()}
}
