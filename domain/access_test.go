package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
)

func TestAccessSpecValidate(t *testing.T) {
	owner := "test-owner"

	cases := []struct {
		domain.AccessSpec
		valid bool
	}{
		{domain.AccessSpec{Owner: owner}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"R", domain.UserType, "id"}}}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"RW", domain.UserType, "id"}}}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"RWX", domain.UserType, "id"}}}, true},

		{domain.AccessSpec{owner, []domain.Scope{{"W", domain.DevType, "id"}}}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"WR", domain.DevType, "id"}}}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"XWR", domain.DevType, "id"}}}, true},

		{domain.AccessSpec{owner, []domain.Scope{{"X", domain.ChanType, "id"}}}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"XR", domain.ChanType, "id"}}}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"XRW", domain.ChanType, "id"}}}, true},
		{domain.AccessSpec{owner, []domain.Scope{{"R", domain.UserType, "id"}, {"RW", domain.DevType, "id"}}}, true},

		{domain.AccessSpec{"", []domain.Scope{{"R", domain.UserType, "id"}}}, false},
		{domain.AccessSpec{owner, []domain.Scope{{"bad", domain.UserType, "id"}}}, false},
		{domain.AccessSpec{owner, []domain.Scope{{"R", "bad", "id"}}}, false},
		{domain.AccessSpec{owner, []domain.Scope{{"R", domain.UserType, ""}}}, false},
		{domain.AccessSpec{owner, []domain.Scope{{"R", domain.UserType, "id"}, {"R", "bad", ""}}}, false},
	}

	for i, c := range cases {
		valid := c.Validate()
		if valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}
