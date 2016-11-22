package domain

import (
	"net/http"
	"strings"
)

const (
	actions  = "CRUD"
	protocol = "http://"

	// UserType defines a canonical user type name.
	UserType string = "users"

	// DevType defines a canonical device type name.
	DevType string = "devices"

	// ChanType defines a canonical channel type name.
	ChanType string = "channels"
)

// AccessRequest specifies a system request that needs to be authorized.
type AccessRequest struct {
	Action string
	Type   string
	Id     string
}

// Validate will determine whether an access request is valid or not. A request
// is considered valid if it has an action specified as either C (create), R
// (retrieve), U (update) or D (delete), and a type that can be either
// 'UserType', 'ChanType' or 'DevType'.
func (a *AccessRequest) Validate() bool {
	if len(a.Action) != 1 || !strings.Contains(actions, a.Action) {
		return false
	}

	return a.Type == ChanType || a.Type == DevType || a.Type == UserType
}

// SetAction sets an action mnemonic to one of C (create), R (retrieve), U
// (update) or D (delete) based on input HTTP method.
func (a *AccessRequest) SetAction(method string) error {
	switch strings.ToUpper(method) {
	case "GET":
		a.Action = "R"
	case "POST":
		a.Action = "C"
	case "PUT":
		a.Action = "U"
	case "DELETE":
		a.Action = "D"
	default:
		return &AuthError{Code: http.StatusBadRequest}
	}

	return nil
}

// SetIdentity sets a resource type and ID, as extracted from the provided URI.
// It is expected URI to have a form of http://<hostname>/<type>/<id>/<other>,
// with only the <hostname> and <type> being required parts.
func (a *AccessRequest) SetIdentity(uri string) error {
	uri = strings.TrimPrefix(uri, protocol)
	parts := strings.Split(uri, "/")

	if len(parts) < 2 {
		return &AuthError{Code: http.StatusBadRequest}
	}

	switch strings.ToLower(parts[1]) {
	case DevType:
		a.Type = DevType
	case ChanType:
		a.Type = ChanType
	case UserType:
		a.Type = UserType
	default:
		return &AuthError{Code: http.StatusBadRequest}
	}

	if len(parts) > 2 {
		a.Id = parts[2]
	}

	return nil
}
