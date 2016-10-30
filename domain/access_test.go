package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth-server/domain"
)

func TestAccessSpecValidate(t *testing.T) {
	cases := []struct {
		domain.AccessSpec
		valid bool
	}{
		{domain.AccessSpec{[]domain.Scope{{"R", domain.UserType, "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{{"RW", domain.UserType, "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{{"RWX", domain.UserType, "id"}}}, true},

		{domain.AccessSpec{[]domain.Scope{{"W", domain.DevType, "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{{"WR", domain.DevType, "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{{"XWR", domain.DevType, "id"}}}, true},

		{domain.AccessSpec{[]domain.Scope{{"X", domain.ChanType, "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{{"XR", domain.ChanType, "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{{"XRW", domain.ChanType, "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{{"R", domain.UserType, "id"}, {"RW", domain.DevType, "id"}}}, true},

		{domain.AccessSpec{}, false},
		{domain.AccessSpec{[]domain.Scope{{"bad", domain.UserType, "id"}}}, false},
		{domain.AccessSpec{[]domain.Scope{{"R", "bad", "id"}}}, false},
		{domain.AccessSpec{[]domain.Scope{{"R", domain.UserType, ""}}}, false},
		{domain.AccessSpec{[]domain.Scope{{"R", domain.UserType, "id"}, {"R", "bad", ""}}}, false},
	}

	for i, c := range cases {
		valid := c.Validate()
		if valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}

func TestAccessRequestValidate(t *testing.T) {
	cases := []struct {
		domain.AccessRequest
		valid bool
	}{
		{domain.AccessRequest{"R", domain.UserType, "id", "owner", "key"}, true},
		{domain.AccessRequest{"W", domain.UserType, "id", "owner", "key"}, true},
		{domain.AccessRequest{"X", domain.UserType, "id", "owner", "key"}, true},
		{domain.AccessRequest{"R", domain.DevType, "id", "owner", "key"}, true},
		{domain.AccessRequest{"W", domain.DevType, "id", "owner", "key"}, true},
		{domain.AccessRequest{"X", domain.DevType, "id", "owner", "key"}, true},
		{domain.AccessRequest{"R", domain.ChanType, "id", "dev", "key"}, true},
		{domain.AccessRequest{"W", domain.ChanType, "id", "dev", "key"}, true},
		{domain.AccessRequest{"X", domain.ChanType, "id", "dev", "key"}, true},

		{domain.AccessRequest{"R", domain.ChanType, "id", "", "key"}, false},
		{domain.AccessRequest{"X", domain.ChanType, "", "dev", "key"}, false},
		{domain.AccessRequest{"W", "bad", "id", "dev", "key"}, false},
		{domain.AccessRequest{"bad", domain.UserType, "id", "owner", "key"}, false},
	}

	for i, c := range cases {
		valid := c.Validate()
		if valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}
