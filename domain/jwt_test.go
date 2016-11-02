package domain_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
)

const key string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlt7ImFjdGlvbnMiOiJSIiwidHlwZSI6ImNoYW5uZWwiLCJpZCI6InRlc3QtaWQifV0sImlzcyI6Im1haW5mbHV4Iiwic3ViIjoidGVzdC1pZCJ9.lvEUcdxg2TX9lpsaCblXs7L7xUaq5nosEgez4vzlQMo"

var access = domain.AccessSpec{
	[]domain.Scope{
		{"R", "channel", "test-id"},
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

			if s.Type != out.Type {
				t.Errorf("case %d: expected %s got %s", i+1, out.Type, s.Type)
			}

			if s.Id != out.Id {
				t.Errorf("case %d: expected %s got %s", i+1, out.Id, s.Id)
			}
		}
	}
}
