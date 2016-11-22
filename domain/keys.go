package domain

import (
	"regexp"
	"sort"
	"strings"
)

// KeyList represents API keys created by user.
type KeyList struct {
	Keys []string `json:"keys"`
}

// Scope represents a resource(s) access scope definition. Each definition
// consists of allowed actions, resource type and its identifier.
type Scope struct {
	Actions string `json:"actions"`
	Type    string `json:"type"`
	Id      string `json:"id"`
}

// KeySpec represents a collection of actions that the key owner can perform.
type KeySpec struct {
	Owner  string  `json:"owner"`
	Scopes []Scope `json:"scopes"`
}

// Validate will try to validate an access specification. The structure will be
// tested against the following conditions: an action can be any permutation of
// "CRUD", a resource can be either "channels", "devices" or "users", and an id
// cannot be empty.
func (a *KeySpec) Validate() bool {
	if a.Owner == "" {
		return false
	}

	for _, s := range a.Scopes {
		if s.Id == "" {
			return false
		}

		if len(s.Actions) == 0 || len(s.Actions) > 4 {
			return false
		}

		items := strings.Split(s.Actions, "")
		sort.Strings(items)
		s.Actions = strings.ToUpper(strings.Join(items, ""))

		if ok, _ := regexp.MatchString("^(C)?(D)?(R)?(U)?$", s.Actions); !ok {
			return false
		}

		s.Type = strings.ToLower(s.Type)
		if s.Type != ChanType && s.Type != DevType && s.Type != UserType {
			return false
		}
	}

	return true
}
