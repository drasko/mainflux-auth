package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth-server/domain"
)

const key string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlt7ImFjdGlvbnMiOiJSIiwicmVzb3VyY2UiOiJjaGFubmVsIiwiaWQiOiJ0ZXN0LWlkIn1dLCJpc3MiOiJtYWluZmx1eCIsInN1YiI6InRlc3QtaWQifQ.QaAdnzbG17SVNb870sj0XKHhO8rPSu_xEvXbeb9PEp4"

var access = domain.AccessSpec{
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
	actual, err := domain.CreateKey(subject, &access)
	if err != nil {
		t.Error("failed to create JWT:", err)
	}

	if actual != key {
		t.Errorf("expected %s got %s", key, actual)
	}
}

func TestContentOf(t *testing.T) {
	cases := []struct {
		in     string
		hasErr bool
	}{
		{key, false},
		{"bad", true},
	}

	for i, c := range cases {
		content, err := domain.ContentOf(c.in)
		if err != nil {
			if !c.hasErr {
				t.Errorf("case %d: didn't expect an error", i+1)
			}

			continue
		}

		if c.hasErr {
			t.Errorf("case %d: expected error to be thrown", i+1)
		}

		if len(content.Scopes) != len(access.Scopes) {
			t.Errorf("case %d: scopes are not of the same length", i+1)
		}

		for j, s := range content.Scopes {
			out := access.Scopes[j]
			if s.Actions != out.Actions {
				t.Errorf("case %d: expected %s got %s", i+1, out.Actions, s.Actions)
			}

			if s.Resource != out.Resource {
				t.Errorf("case %d: expected %s got %s", i+1, out.Resource, s.Resource)
			}

			if s.Id != out.Id {
				t.Errorf("case %d: expected %s got %s", i+1, out.Id, s.Id)
			}
		}
	}
}
