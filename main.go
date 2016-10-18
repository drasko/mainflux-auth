package main

import (
	"net/http"

	"github.com/mainflux/mainflux-auth-server/api"
)

const port string = ":8180"

func main() {
	http.ListenAndServe(port, api.Server())
}
