package oauth

type TokenInfo struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id"`
	Expires int64  `json:"expires"`
}

type UserInfo struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Company string `json:"company"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
