package core

import (
	"strings"
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
	scopes models.Scopes = models.Scopes{[]models.Scope{scope}}
)

func TestCreate(t *testing.T) {
	actual, err := Create(scopes)
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

		if strings.Compare(item.Resource, s.Resource) != 0 {
			t.Errorf("%d: expected res %s got %s", i, item.Resource, s.Resource)
		}

		if strings.Compare(item.Id, s.Id) != 0 {
			t.Errorf("%d: expected ID %s got %s", i, item.Id, s.Id)
		}

		if strings.Compare(item.Actions, s.Actions) != 0 {
			t.Errorf("%d: expected actions %s got %s", i, item.Actions, s.Actions)
		}
	}
}
