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

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
