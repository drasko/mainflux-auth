package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth-server/api/handlers"
)

func Server() http.Handler {
	mux := httprouter.New()
	mux.GET("/status", handlers.HealthCheck)

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
