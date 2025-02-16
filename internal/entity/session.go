package entity

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Session struct {
	ID                 uuid.UUID                `json:"id"`
	Name               string                   `json:"name"`
	ServerURL          string                   `json:"server_url"`
	Authorization      *Auth                    `json:"authorization"`
	ActiveConnection   *grpc.ClientConn         `json:"-"`
	AvailableServices  []string                 `json:"-"`
	SelectedMethod     *string                  `json:"selected_method"`
	AvailableMethods   []string                 `json:"-"`
	RequestPayload     string                   `json:"request_payload"`
	DescriptorSource   grpcurl.DescriptorSource `json:"-"`
	Metadata           []*Metadata              `json:"metadata"`
	AllowUnknownFields bool                     `json:"allow_unknown_fields"`

	// SSL Certificates
	EnableTLS          bool  `json:"enable_tls"`
	SSLCert            *Cert `json:"ssl_cert"`
	InsecureSkipVerify bool  `json:"insecure_skip_verify"`

	// GRPC Call Options
	MaxMsgSz       int     `json:"max_msg_sz"`
	ConnectTimeout float64 `json:"connect_timeout"`
	MaxTimeOut     float64 `json:"max_timeout"`
	KeepAliveTime  float64 `json:"keepalive_time"`
}

// ParseMetadata used to convert the metadata and authorization parameters to array of strings
func (s *Session) ParseMetadata() []string {
	// Convert metadata to array of strings
	result := make([]string, 0, len(s.Metadata))
	for _, v := range s.Metadata {
		if v.Active {
			result = append(result, v.Key+": "+v.Value)
		}
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
