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
	mux.POST("/sessions", login)
	mux.POST("/users", registerUser)

	mux.POST("/api-keys", addKey)
	mux.PUT("/api-keys/:key", updateKey)

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
