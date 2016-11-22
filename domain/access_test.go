package domain_test

import (
	"fmt"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
)

func TestAccessRequestValidate(t *testing.T) {
	cases := []struct {
		domain.AccessRequest
		valid bool
	}{
		{domain.AccessRequest{"C", domain.UserType, "id"}, true},
		{domain.AccessRequest{"R", domain.UserType, "id"}, true},
		{domain.AccessRequest{"U", domain.UserType, "id"}, true},
		{domain.AccessRequest{"D", domain.UserType, "id"}, true},
		{domain.AccessRequest{"C", domain.ChanType, "id"}, true},
		{domain.AccessRequest{"R", domain.ChanType, "id"}, true},
		{domain.AccessRequest{"U", domain.ChanType, "id"}, true},
		{domain.AccessRequest{"D", domain.ChanType, "id"}, true},
		{domain.AccessRequest{"C", domain.DevType, "id"}, true},
		{domain.AccessRequest{"R", domain.DevType, "id"}, true},
		{domain.AccessRequest{"U", domain.DevType, "id"}, true},
		{domain.AccessRequest{"D", domain.DevType, "id"}, true},
		{domain.AccessRequest{"C", domain.ChanType, ""}, true},
		{domain.AccessRequest{"C", "bad", "id"}, false},
		{domain.AccessRequest{"bad", domain.UserType, "id"}, false},
	}

	for i, c := range cases {
		if valid := c.Validate(); valid != c.valid {
			t.Errorf("case %d: expected %t got %t", i+1, c.valid, valid)
		}
	}
}

func TestSetAction(t *testing.T) {
	cases := []struct {
		method   string
		expected string
	}{
		{"POST", "C"},
		{"GET", "R"},
		{"PUT", "U"},
		{"DELETE", "D"},
		{"bad", ""},
	}

	for i, c := range cases {
		req := domain.AccessRequest{}
		req.SetAction(c.method)

		if req.Action != c.expected {
			t.Errorf("case %d: expected %s got %s", i+1, c.expected, req.Action)
		}
	}
}

func TestSetIdentity(t *testing.T) {
	id := "id"

	cases := []struct {
		uri     string
		resType string
		id      string
	}{
		{fmt.Sprintf("<hostname>/%s/%s", domain.DevType, id), domain.DevType, id},
		{fmt.Sprintf("<hostname>/%s/%s", domain.ChanType, id), domain.ChanType, id},
		{fmt.Sprintf("<hostname>/%s/%s", domain.UserType, id), domain.UserType, id},
		{fmt.Sprintf("<hostname>/%s", domain.UserType), domain.UserType, ""},
		{fmt.Sprintf("http://<hostname>/%s/%s", domain.DevType, id), domain.DevType, id},
		{fmt.Sprintf("http://<hostname>/%s/%s", domain.ChanType, id), domain.ChanType, id},
		{fmt.Sprintf("http://<hostname>/%s/%s", domain.UserType, id), domain.UserType, id},
		{fmt.Sprintf("http://<hostname>/%s", domain.UserType), domain.UserType, ""},
		{"<hostname>", "", ""},
		{"", "", ""},
	}

	for i, c := range cases {
		req := domain.AccessRequest{}
		req.SetIdentity(c.uri)

		if req.Type != c.resType {
			t.Errorf("case %d: invalid type, expected %s got %s", i+1, c.resType, req.Type)
		}

		if req.Id != c.id {
			t.Errorf("case %d: invalid id, expected %s got %s", i+1, c.id, req.Id)
		}
	}
}
