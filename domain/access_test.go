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
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"R", "user", "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"RW", "user", "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"RWX", "user", "id"}}}, true},

		{domain.AccessSpec{[]domain.Scope{domain.Scope{"W", "device", "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"WR", "device", "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"XWR", "device", "id"}}}, true},

		{domain.AccessSpec{[]domain.Scope{domain.Scope{"X", "channel", "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"XR", "channel", "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"XRW", "channel", "id"}}}, true},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"R", "user", "id"}, domain.Scope{"RW", "device", "id"}}}, true},

		{domain.AccessSpec{}, false},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"bad", "user", "id"}}}, false},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"R", "bad", "id"}}}, false},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"R", "user", ""}}}, false},
		{domain.AccessSpec{[]domain.Scope{domain.Scope{"R", "user", "id"}, domain.Scope{"R", "bad", ""}}}, false},
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
		{domain.AccessRequest{"R", "user", "id", "owner", "key"}, true},
		{domain.AccessRequest{"W", "user", "id", "owner", "key"}, true},
		{domain.AccessRequest{"X", "user", "id", "owner", "key"}, true},
		{domain.AccessRequest{"R", "device", "id", "owner", "key"}, true},
		{domain.AccessRequest{"W", "device", "id", "owner", "key"}, true},
		{domain.AccessRequest{"X", "device", "id", "owner", "key"}, true},
		{domain.AccessRequest{"R", "channel", "id", "dev", "key"}, true},
		{domain.AccessRequest{"W", "channel", "id", "dev", "key"}, true},
		{domain.AccessRequest{"X", "channel", "id", "dev", "key"}, true},

		{domain.AccessRequest{"R", "channel", "id", "", "key"}, false},
		{domain.AccessRequest{"X", "channel", "", "dev", "key"}, false},
		{domain.AccessRequest{"W", "bad", "id", "dev", "key"}, false},
		{domain.AccessRequest{"bad", "user", "id", "owner", "key"}, false},
	}

	for i, c := range cases {
		valid := c.Validate()
		if valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}
