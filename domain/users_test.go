package domain_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/mainflux/mainflux-auth-server/domain"
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

	for _, c := range cases {
		user, err := domain.CreateUser(c.username, c.username)
		if err != nil {
			_, ok := err.(*domain.ServiceError)
			if !ok {
				t.Errorf("all errors must be ServiceError")
			}
		}

		if user.Username != c.username {
			t.Errorf("expected %s got %s", c.username, user.Username)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.password))
		if err != nil {
			t.Errorf("invalid password")
		}

		p, err := domain.ScopesOf(user.MasterKey)
		if err != nil {
			t.Errorf("invalid master key")
		}

		if len(p.Scopes) != 1 || p.Scopes[0] != scope {
			t.Errorf("incompatible master key payload")
		}

		if user.Id == "" {
			t.Errorf("empty user ID")
		}
	}
}
