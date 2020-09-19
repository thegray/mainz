package common

import "github.com/dgrijalva/jwt-go"

type AccessClaim struct {
	jwt.StandardClaims
	// UserID   int    `json:"userid"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type RefreshClaim struct {
	jwt.StandardClaims
	// UserID   int    `json:"userid"`
	Name  string `json:"name"`
	UUIDR string `json:"uuidr"`
	UUIDA string `json:"uuida"`
}
