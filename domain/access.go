package domain

import "strings"

// AccessRequest specifies a system request that needs to be authorized.
type AccessRequest struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Id     string `json:"id"`
}

// Validate will determine whether an access request is valid or not. A request
// is considered valid if it has a non-empty Id, an action specified as either
// R (read), W (write) or X (execute), and a type that can be either "channel",
// "device" or "user".
func (a *AccessRequest) Validate() bool {
	if a.Id == "" {
		return false
	}

	a.Action = strings.ToUpper(a.Action)
	if len(a.Action) != 1 || !strings.Contains("RWX", a.Action) {
		return false
	}

	a.Type = strings.ToLower(a.Type)
	return a.Type == ChanType || a.Type == DevType || a.Type == UserType
}
