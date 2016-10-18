package core

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/mainflux/mainflux-auth-server/models"
)

const (
	key    string = "mainflux-api-key"
	issuer string = "mainflux"
)

type tokenClaims struct {
	models.Scopes
	jwt.StandardClaims
}

func Create(scopes models.Scopes) (string, error) {
	claims := tokenClaims{
		scopes,
		jwt.StandardClaims{Issuer: issuer},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return raw, nil
}

func ScopesOf(rawToken string) (*models.Scopes, error) {
	token, err := jwt.ParseWithClaims(
		rawToken,
		&tokenClaims{},
		func(val *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return &claims.Scopes, nil
	}

	return nil, errors.New("Failed to retrieve claims.")
}
