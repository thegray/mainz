package response

type Tokens struct {
	AccessToken  string `json:"access_token"`
	Exp          int64  `json:"expiry"`
	RefreshToken string `json:"refresh_token"`
}
