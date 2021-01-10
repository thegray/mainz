package request

type SBByIdReq struct {
	ID *int `json:"id" validate:"required"`
}

type GetSecretById struct {
	ID  *int    `json:"id" validate:"required"`
	Key *string `json:"k" validate:"required"`
}
type UpdateSBReq struct {
	ID       *int    `json:"id" validate:"required"`
	Username string  `json:"un"`
	Email    string  `json:"em" validate:"omitempty,email"`
	Platform *string `json:"pf" validate:"required"`
	Details  string  `json:"dt"`
	CredID   *int    `json:"cid" validate:"required"`
}

type UpdateSecretReq struct {
	ID         *int    `json:"id" validate:"required"`
	SecretPass string  `json:"sp"`
	Key        *string `json:"k" validate:"required"`
	SecretInfo string  `json:"si"`
}

type RevealSecretReq struct {
	ID  *int    `json:"id" validate:"required"`
	Key *string `json:"k" validate:"required"`
}

type TestValidationReq struct {
	Test string `json:"test" validate:"excludes= "` //input string should not contain ' '
}
