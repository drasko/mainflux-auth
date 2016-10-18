package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

func Server() http.Handler {
	mux := httprouter.New()
	mux.GET("/status", HealthCheck)

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
