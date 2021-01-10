package request

type ByIdReq struct {
	ID *int `json:"id" validate:"required"`
}

type CreateAccSBReq struct {
	Username      string `json:"un" validate:"required"`
	EmailPhonenum string `json:"ep" validate:"omitempty,email"`
	Platform      string `json:"pf" validate:"required"`
	Details       string `json:"dt"`
	CredID        int    `json:"cid" validate:"required"` //TODO: maybe not needed, can use auth middleware
	CategoryID    []int  `json:"catid"`                   //TODO: need more works

	SecretPass string `json:"sp"`
	SecretInfo string `json:"si"`
	Key        string `json:"k"`
}
