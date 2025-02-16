package entity

type AuthType string

const (
	AuthTypeBearer AuthType = "bearer"
	AuthTypeAPIKey AuthType = "api"
)

type Auth struct {
	AuthType       AuthType              `json:"auth_type"`
	BearerToken    *AuthValueBearerToken `json:"bearer_token"`
	AuthTypeAPIKey *AuthValueTypeAPIKey  `json:"api_key"`
}

type AuthValueBearerToken struct {
	Token string `json:"token"`
}

type AuthValueTypeAPIKey struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
