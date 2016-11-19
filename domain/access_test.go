package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
)

func TestKeySpecValidate(t *testing.T) {
	owner := "test-owner"

	cases := []struct {
		domain.KeySpec
		valid bool
	}{
		{domain.KeySpec{Owner: owner}, true},
		{domain.KeySpec{owner, []domain.Scope{{"R", domain.UserType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"RW", domain.UserType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"RWX", domain.UserType, "id"}}}, true},

		{domain.KeySpec{owner, []domain.Scope{{"W", domain.DevType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"WR", domain.DevType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"XWR", domain.DevType, "id"}}}, true},

		{domain.KeySpec{owner, []domain.Scope{{"X", domain.ChanType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"XR", domain.ChanType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"XRW", domain.ChanType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"R", domain.UserType, "id"}, {"RW", domain.DevType, "id"}}}, true},

		{domain.KeySpec{"", []domain.Scope{{"R", domain.UserType, "id"}}}, false},
		{domain.KeySpec{owner, []domain.Scope{{"bad", domain.UserType, "id"}}}, false},
		{domain.KeySpec{owner, []domain.Scope{{"R", "bad", "id"}}}, false},
		{domain.KeySpec{owner, []domain.Scope{{"R", domain.UserType, ""}}}, false},
		{domain.KeySpec{owner, []domain.Scope{{"R", domain.UserType, "id"}, {"R", "bad", ""}}}, false},
	}

	for i, c := range cases {
		valid := c.Validate()
		if valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}
