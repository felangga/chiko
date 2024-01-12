package entity

type AuthType string

const (
	AuthTypeBearer AuthType = "bearer"
)

type Auth struct {
	AuthType    AuthType         `json:"auth_type"`
	BearerToken *AuthBearerToken `json:"bearer_token"`
}

type AuthBearerToken struct {
	Token string `json:"token"`
}
