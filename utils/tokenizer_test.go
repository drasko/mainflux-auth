package utils

import (
	"testing"

	"github.com/mainflux/mainflux-auth-server/models"
)

const token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlt7ImFjdGlvbnMiOiJSIiwicmVzb3VyY2UiOiJjaGFubmVsIiwiaWQiOiJ0ZXN0LWlkIn1dLCJpc3MiOiJtYWluZmx1eCJ9.XumLIqHS4YBhhPmS1wBcf915s48zRnMMwQiIXnSzcFU"

var (
	scope models.Scope = models.Scope{
		Actions:  "R",
		Resource: "channel",
		Id:       "test-id",
	}
	scopes models.Scopes = models.Scopes{Items: []models.Scope{scope}}
)

func TestSetKey(t *testing.T) {
	oldKey := key
	newKey := "dummy-key"
	SetKey(newKey)

	if key != newKey {
		t.Errorf("expected %s got %s", newKey, key)
	}

	SetKey(oldKey)
}

func TestCreateToken(t *testing.T) {
	actual, err := CreateToken(scopes)
	if err != nil {
		t.Error("failed to create JWT:", err)
	}

	if actual != token {
		t.Errorf("expected %s got %s", token, actual)
	}
}

func TestScopesOf(t *testing.T) {
	actual, err := ScopesOf(token)
	if err != nil {
		t.Error("couldn't extract scopes:", err)
	}

	if len(scopes.Items) != len(actual.Items) {
		t.Errorf("different number of containing scopes")
	}

	for i, s := range scopes.Items {
		item := actual.Items[i]

		if item.Resource != s.Resource {
			t.Errorf("%d: expected res %s got %s", i, item.Resource, s.Resource)
		}

		if item.Id != s.Id {
			t.Errorf("%d: expected ID %s got %s", i, item.Id, s.Id)
		}

		if item.Actions != s.Actions {
			t.Errorf("%d: expected actions %s got %s", i, item.Actions, s.Actions)
		}
	}
}
