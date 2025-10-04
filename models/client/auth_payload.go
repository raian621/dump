package client

type AuthPayload struct {
	AccessToken  string `json:"access,omitempty"`
	RefreshToken string `json:"refresh,omitempty"`
}
