package request

type LoginReq struct {
	User string `json:"un"`
	Pass string `json:"pw"`
}

type ChangePassReq struct {
	OldPass string `json:"op"`
	NewPass string `json:"np"`
}

type RefreshReq struct {
	RefreshToken string `json:"token"`
}
