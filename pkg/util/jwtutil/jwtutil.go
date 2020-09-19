package jwtutil

import (
	"errors"
	"pbk-main/pkg/model/common"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_ACCESS_KEY = []byte("need_to_change_this")    //todo: move to config
var JWT_REFRESH_KEY = []byte("need_to_change_this_2") //todo: move to config

func CreateToken(claims interface{}, tokenType string, expiry int64) (string, error) {

	var tokenKey []byte
	var c jwt.Claims
	switch tokenType {
	case "access":
		tokenKey = JWT_ACCESS_KEY
		c = claims.(common.AccessClaim)
	case "refresh":
		tokenKey = JWT_REFRESH_KEY
		c = claims.(common.RefreshClaim)
	}
	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, c)
	signedToken, err := token.SignedString(tokenKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func TokenValid(tokenStr string, tokenType string) error {
	token, err := VerifyToken(tokenStr, tokenType)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func VerifyToken(tokenStr string, tokenType string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check if token signed with "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signing method")
		}
		// return []byte(os.Getenv("ACCESS_SECRET")), nil
		if tokenType == "access" {
			return JWT_ACCESS_KEY, nil
		} else {
			return JWT_REFRESH_KEY, nil
		}

	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenToAccessClaim(token *jwt.Token) (*common.AccessClaim, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}

	// get uuid
	refreshUuid, ok := claims["uuid"].(string) //convert the interface to string
	if !ok {
		return nil, errors.New("missing param 'uid'")
	}
	// get username
	username, ok := claims["name"].(string)
	if !ok {
		return nil, errors.New("missing param 'un'")
	}
	//get userid
	// userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userid"]), 10, 0)
	// if err != nil {
	// 	return nil, err
	// }
	rc := &common.AccessClaim{
		// UserID:   int(userId),
		Name: username,
		UUID: refreshUuid}
	return rc, nil
}

func ExtractTokenToRefreshClaim(token *jwt.Token) (*common.RefreshClaim, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}

	// get refresh uuid
	refreshUuid, ok := claims["uuidr"].(string) //convert the interface to string
	if !ok {
		return nil, errors.New("missing param 'uidr'")
	}
	// get access uuid
	accessUuid, ok := claims["uuida"].(string) //convert the interface to string
	if !ok {
		return nil, errors.New("missing param 'uida'")
	}
	// get username
	username, ok := claims["name"].(string)
	if !ok {
		return nil, errors.New("missing param 'un'")
	}
	// userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userid"]), 10, 0)
	// if err != nil {
	// 	return nil, err
	// }
	rc := &common.RefreshClaim{
		// UserID: int(userId),
		Name:  username,
		UUIDR: refreshUuid,
		UUIDA: accessUuid}
	return rc, nil
}
