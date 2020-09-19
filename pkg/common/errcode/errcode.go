package errcode

// dsta
// d = domain
// s = service
// t = error type
// a = additional info
const (
	GUnknownError = "0000"
	GBadRequest   = "0010"
	GAuthError    = "0020"

	CREDLoginError                     = "1010"
	CREDRefreshError_FailCreateToken   = "1110"
	CREDRefreshError_MismatchTokenType = "1111"

	SBUnknownError     = "2000"
	SBNotExist         = "2100"
	SBFailInsert       = "2110"
	SBFailUpdate       = "2120"
	SBFailUpdateSecret = "2130"
	SBFailReveal       = "2140"
)
