package entity_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/felangga/chiko/pkg/entity"
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
				Metadata: []*entity.Metadata{
					{
						Active: true,
						Key:    "Authorization",
						Value:  "Bearer testing-token",
					},
				},
				Authorization:  &entity.Auth{},
				ServerURL:      "localhost:50051",
				RequestPayload: "testing-payload",
			},
			expected: []string{
				"Authorization: Bearer testing-token",
			},
		},
		{
			name: "success, no authorization",
			session: entity.Session{
				Name: "Session 1",
				Metadata: []*entity.Metadata{
					{
						Active: true,
						Key:    "name",
						Value:  "testing-name",
					},
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
		t.Run(tc.name, func(t *testing.T) {
			parsed := tc.session.ParseMetadata()

			if diff := cmp.Diff(tc.expected, parsed); diff != "" {
				t.Fatal(diff)
			}
		})
	}

}
