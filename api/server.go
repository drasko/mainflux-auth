package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

// Server binds API endpoints to their handlers and initializes middleware.
func Server() http.Handler {
	mux := httprouter.New()
	mux.GET("/status", healthCheck)
	mux.POST("/users", registerUser)
	mux.POST("/users/:user_id/api-keys", addUserKey)
	mux.POST("/users/:user_id/devices/:device_id/api-keys", addDeviceKey)

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
