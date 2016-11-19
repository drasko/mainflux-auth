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

	mux.GET("/api-keys", fetchKeys)
	mux.POST("/api-keys", addKey)

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
