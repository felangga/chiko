package entity

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Session struct {
	ID                uuid.UUID                `json:"id"`
	Name              string                   `json:"name"`
	ServerURL         string                   `json:"server_url"`
	Authorization     *Auth                    `json:"authorization"`
	ActiveConnection  *grpc.ClientConn         `json:"-"`
	AvailableServices []string                 `json:"-"`
	SelectedMethod    *string                  `json:"selected_method"`
	AvailableMethods  []string                 `json:"-"`
	RequestPayload    string                   `json:"request_payload"`
	DescriptorSource  grpcurl.DescriptorSource `json:"-"`
	Metadata          map[string]string        `json:"metadata"`
}

// ParseMetadata used to convert the metadata and authorization parameters to array of strings
func (s *Session) ParseMetadata() []string {
	result := make([]string, 0, len(s.Metadata))
	for k, v := range s.Metadata {
		result = append(result, k+": "+v)
	}

	// Add the authorization to header
	if s.Authorization != nil {
		switch s.Authorization.AuthType {
		case AuthTypeBearer:
			result = append(result, "Authorization: Bearer "+s.Authorization.BearerToken.Token)
		}
	}

	return result
}
