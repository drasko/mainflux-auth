package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

func Server() http.Handler {
	mux := httprouter.New()
	mux.GET("/status", HealthCheck)
	mux.POST("/users", RegisterUser)
	mux.POST("/users/:user_id/api-keys", AddUserKey)
	mux.POST("/users/:user_id/devices/:device_id/api-keys", AddDeviceKey)

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
