package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
)

func TestKeySpecValidate(t *testing.T) {
	owner := "test-owner"

	cases := []struct {
		spec  domain.KeySpec
		valid bool
	}{
		{domain.KeySpec{Owner: owner}, true},
		{domain.KeySpec{owner, []domain.Scope{{"C", domain.UserType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"CR", domain.UserType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"CRUD", domain.UserType, "id"}}}, true},

		{domain.KeySpec{owner, []domain.Scope{{"C", domain.DevType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"CR", domain.DevType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"CRUD", domain.DevType, "id"}}}, true},

		{domain.KeySpec{owner, []domain.Scope{{"C", domain.ChanType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"CR", domain.ChanType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"CRU", domain.ChanType, "id"}}}, true},
		{domain.KeySpec{owner, []domain.Scope{{"CRUD", domain.UserType, "id"}, {"CR", domain.DevType, "id"}}}, true},

		{domain.KeySpec{"", []domain.Scope{{"C", domain.UserType, "id"}}}, false},
		{domain.KeySpec{owner, []domain.Scope{{"bad", domain.UserType, "id"}}}, false},
		{domain.KeySpec{owner, []domain.Scope{{"C", "bad", "id"}}}, false},
	}

	for i, c := range cases {
		valid := c.spec.Validate()
		if valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}
