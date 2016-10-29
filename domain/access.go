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
	Actions string `json:"actions"`
	Type    string `json:"type"`
	Id      string `json:"id"`
}

// AccessSpec represents a collection of resource access scope. It will be
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

		s.Type = strings.ToLower(s.Type)
		if s.Type != "channel" && s.Type != "device" && s.Type != "user" {
			return false
		}
	}

	return true
}

// AccessRequest specifies a system request that needs to be authorized.
type AccessRequest struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Id     string `json:"id"`
	Owner  string `json:"owner"`
	Key    string `json:"key"`
}

// Validate will try to validate an access request. The structure will be
// tested against the following conditions: an action can be any value from
// "RWX", a type can be either "channel", "device" or "user", and none of the
// remaining fields can be empty.
func (a *AccessRequest) Validate() bool {
	if len(a.Action) != 1 || a.Type == "" || a.Id == "" || a.Owner == "" || a.Key == "" {
		return false
	}

	a.Action = strings.ToUpper(a.Action)
	if !strings.Contains("RWX", a.Action) {
		return false
	}

	a.Type = strings.ToLower(a.Type)
	if a.Type != "channel" && a.Type != "device" && a.Type != "user" {
		return false
	}

	return true
}
