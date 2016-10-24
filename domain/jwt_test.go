package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth-server/domain"
)

const key string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlt7ImFjdGlvbnMiOiJSIiwicmVzb3VyY2UiOiJjaGFubmVsIiwiaWQiOiJ0ZXN0LWlkIn1dLCJpc3MiOiJtYWluZmx1eCIsInN1YiI6InRlc3QtaWQifQ.QaAdnzbG17SVNb870sj0XKHhO8rPSu_xEvXbeb9PEp4"

var payload = domain.Payload{
	[]domain.Scope{
		domain.Scope{
			Actions:  "R",
			Resource: "channel",
			Id:       "test-id",
		},
	},
}

func TestCreateKey(t *testing.T) {
	subject := "test-id"
	actual, err := domain.CreateKey(subject, &payload)
	if err != nil {
		t.Error("failed to create JWT:", err)
	}

	if actual != key {
		t.Errorf("expected %s got %s", key, actual)
	}
}

func TestScopesOf(t *testing.T) {
	actual, err := domain.ScopesOf(key)
	if err != nil {
		t.Error("failed to extract scopes:", err)
	}

	if len(actual.Scopes) != len(payload.Scopes) {
		t.Error("scopes are not of the same length")
	}

	for i, s := range actual.Scopes {
		out := payload.Scopes[i]

		if s.Actions != out.Actions {
			t.Errorf("expected %s got %s", out.Actions, s.Actions)
		}

		if s.Resource != out.Resource {
			t.Errorf("expected %s got %s", out.Resource, s.Resource)
		}

		if s.Id != out.Id {
			t.Errorf("expected %s got %s", out.Id, s.Id)
		}
	}
}
