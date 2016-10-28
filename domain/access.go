package domain

import (
	"regexp"
	"sort"
	"strings"
)

// Scope represents a resource(s) access scope definition. Each definition
// consists of allowed actions, resource type and its identifier. Keep in
// mind that the '*' is treated as wild card - it can be used to indicate
// "all" available resource types and "all" resource identifiers.
type Scope struct {
	Actions  string `json:"actions"`
	Resource string `json:"resource"`
	Id       string `json:"id"`
}

// Payload represents a collection of resource access scope. It will be
// embedded into the generated API key.
type AccessSpec struct {
	Scopes []Scope `json:"scopes"`
}

// Validate will try to validate an access specification. The structure will be
// tested against the following conditions: an action can be any permutation of
// "RWX", a resource can be either "channel", "device" or "user", and an id
// cannot be empty.
func (a *AccessSpec) Validate() bool {
	if len(a.Scopes) == 0 {
		return false
	}

	for _, s := range a.Scopes {
		if s.Id == "" {
			return false
		}

		if len(s.Actions) == 0 || len(s.Actions) > 3 {
			return false
		}

		items := strings.Split(s.Actions, "")
		sort.Strings(items)
		s.Actions = strings.ToUpper(strings.Join(items, ""))

		if ok, _ := regexp.MatchString("^(R)?(W)?(X)?$", s.Actions); !ok {
			return false
		}

		s.Resource = strings.ToLower(s.Resource)
		if s.Resource != "channel" && s.Resource != "device" && s.Resource != "user" {
			return false
		}
	}

	return true
}

// AccessRequest specifies a system request that needs to be authorized.
type AccessRequest struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Id       string `json:"id"`
	Device   string `json:"device"`
	Key      string `json:"key"`
}

// Validate will try to validate an access request. The structure will be
// tested against the following conditions: an action can be any value from
// "RWX", a resource can be either "channel", "device" or "user", and an id
// cannot be empty.
func (a *AccessRequest) Validate() bool {
	if a.Id == "" {
		return false
	}

	if len(a.Action) != 1 {
		return false
	}

	a.Action = strings.ToUpper(a.Action)
	if !strings.Contains("RWX", a.Action) {
		return false
	}

	a.Resource = strings.ToLower(a.Resource)
	if a.Resource != "channel" && a.Resource != "device" && a.Resource != "user" {
		return false
	}

	if a.Resource == "channel" && a.Device == "" {
		return false
	}

	return true
}
