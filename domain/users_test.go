package domain_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/mainflux/mainflux-auth/domain"
)

func TestCreateUser(t *testing.T) {
	scope := domain.Scope{"RWX", "*", "*"}

	// TODO: design negative cases
	cases := []struct {
		username string
		password string
	}{
		{"x", "x"},
		{"y", "y"},
	}

	for i, c := range cases {
		user, err := domain.CreateUser(c.username, c.username)
		if err != nil {
			_, ok := err.(*domain.ServiceError)
			if !ok {
				t.Errorf("case %d: all errors must be ServiceError", i+1)
			}
		}

		if user.Username != c.username {
			t.Errorf("case %d: expected %s got %s", i+1, c.username, user.Username)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.password))
		if err != nil {
			t.Errorf("case %d: invalid password", i+1)
		}

		p, err := domain.ContentOf(user.MasterKey)
		if err != nil {
			t.Errorf("case %d: invalid master key", i+1)
		}

		if len(p.Scopes) != 1 || p.Scopes[0] != scope {
			t.Errorf("case %d: incompatible master key scope", i+1)
		}

		if user.Id == "" {
			t.Errorf("case %d: empty user ID", i+1)
		}
	}
}
