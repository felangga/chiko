package entity_test

import (
	"testing"

	"chiko/pkg/entity"

	"github.com/google/go-cmp/cmp"
)

func TestParseMetadata(t *testing.T) {
	test := []struct {
		name     string
		session  entity.Session
		expected []string
	}{
		{
			name: "success, normal request",
			session: entity.Session{
				Name: "Session 1",
				Metadata: map[string]string{
					"name": "testing-name",
				},
				Authorization: &entity.Auth{
					AuthType: entity.AuthTypeBearer,
					BearerToken: &entity.AuthValueBearerToken{
						Token: "testing-token",
					},
				},
				ServerURL:      "localhost:50051",
				RequestPayload: "testing-payload",
			},
			expected: []string{
				"name: testing-name",
				"Authorization: Bearer testing-token",
			},
		},
		{
			name: "success, no authorization",
			session: entity.Session{
				Name: "Session 1",
				Metadata: map[string]string{
					"name": "testing-name",
				},
				ServerURL:      "localhost:50051",
				RequestPayload: "testing-payload",
			},
			expected: []string{
				"name: testing-name",
			},
		},
		{
			name: "success, empty metadata",
			session: entity.Session{
				Name:           "Session 1",
				ServerURL:      "localhost:50051",
				RequestPayload: "testing-payload",
			},
			expected: []string{},
		},
	}

	for _, tc := range test {
		parsed := tc.session.ParseMetadata()

		if diff := cmp.Diff(tc.expected, parsed); diff != "" {
			t.Fatal(diff)
		}
	}

}
