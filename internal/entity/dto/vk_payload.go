package dto

type User struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Avatar     string `json:"avatar"`
	AvatarBase string `json:"avatar_base,omitempty"`
	Phone      string `json:"phone"`
}

type Payload struct {
	Type  string `json:"type"`
	Auth  int    `json:"auth"`
	User  User   `json:"user"`
	Token string `json:"token"`
	TTL   int    `json:"ttl"`
	UUID  string `json:"uuid"`
	Hash  string `json:"hash"`
}

type VkBodyRequest struct {
	Payload Payload `json:"payload"`
}

type VkAPIPayload struct {
	AccessToken              string  `json:"access_token"`
	AccessTokenID            int     `json:"access_token_id"`
	AdditionalSignupRequired bool    `json:"additional_signup_required"`
	Email                    string  `json:"email"`
	ExpiresIn                float64 `json:"expires_in"`
	IsPartial                bool    `json:"is_partial"`
	IsService                bool    `json:"is_service"`
	Source                   int     `json:"source"`
	SourceDescription        string  `json:"source_description"`
	UserID                   float64 `json:"user_id"`
}

type VkAPIBodyResponse struct {
	Response VkAPIPayload `json:"response"`
}
