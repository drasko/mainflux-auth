package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func checkAccess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) != 2 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	req := domain.AccessRequest{}
	req.SetAction(r.Header.Get("X-Action"))
	req.SetIdentity(r.Header.Get("X-Resource"))

	if err := services.CheckPermissions(token[1], req); err != nil {
		authErr := err.(*domain.AuthError)
		w.WriteHeader(authErr.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
}
