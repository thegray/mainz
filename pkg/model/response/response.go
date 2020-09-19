package response

type Revealed struct {
	SecretPass string `json:"sp"`
	SecretInfo string `json:"si"`
}
