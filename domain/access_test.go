package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
)

func TestAccessRequestValidate(t *testing.T) {
	cases := []struct {
		domain.AccessRequest
		valid bool
	}{
		{domain.AccessRequest{"R", domain.UserType, "id"}, true},
		{domain.AccessRequest{"W", domain.UserType, "id"}, true},
		{domain.AccessRequest{"X", domain.UserType, "id"}, true},
		{domain.AccessRequest{"R", domain.DevType, "id"}, true},
		{domain.AccessRequest{"W", domain.DevType, "id"}, true},
		{domain.AccessRequest{"X", domain.DevType, "id"}, true},
		{domain.AccessRequest{"R", domain.ChanType, "id"}, true},
		{domain.AccessRequest{"W", domain.ChanType, "id"}, true},
		{domain.AccessRequest{"X", domain.ChanType, "id"}, true},
		{domain.AccessRequest{"R", domain.ChanType, ""}, false},
		{domain.AccessRequest{"W", "bad", "id"}, false},
		{domain.AccessRequest{"bad", domain.UserType, "id"}, false},
	}

	for i, c := range cases {
		if valid := c.Validate(); valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}
