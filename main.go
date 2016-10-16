package main

import (
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth-server/controllers"
)

func main() {
	mux := httprouter.New()
	mux.GET("/", controllers.Greeter)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":8180")
}
