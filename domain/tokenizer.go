package domain

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

const issuer string = "mainflux"

// default key value
var key string = "mainflux-api-key"

type tokenClaims struct {
	Payload
	jwt.StandardClaims
}

func SetKey(newKey string) {
	key = newKey
}

func CreateToken(subject string, payload *Payload) (string, error) {
	claims := tokenClaims{
		*payload,
		jwt.StandardClaims{
			Issuer:  issuer,
			Subject: subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return raw, nil
}

func ScopesOf(rawToken string) (Payload, error) {
	var payload Payload

	token, err := jwt.ParseWithClaims(
		rawToken,
		&tokenClaims{},
		func(val *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
	)

	if err != nil {
		return payload, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.Payload, nil
	}

	return payload, errors.New("failed to retrieve claims")
}
