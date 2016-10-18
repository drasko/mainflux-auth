package api

import (
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth-server/api/handlers"
)

const port string = ":8180"

func NewServer() {
	mux := httprouter.New()
	mux.GET("/status", handlers.HealthCheck)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(port)
}
