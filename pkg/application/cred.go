package application

import (
	"errors"
	"log"
	"pbk-main/pkg/model/common"
	"pbk-main/pkg/model/response"
	"pbk-main/pkg/store"
	"pbk-main/pkg/store/memstore"
	"pbk-main/pkg/util/jwtutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CredApp struct {
	store *store.CredStore
}

type tokenDetails struct {
	username     string
	userid       int
	accessToken  string
	refreshToken string
	accessUuid   string
	refreshUuid  string
	atExpires    int64
	rtExpires    int64
}

var APPLICATION_NAME = "pbk-main" //todo: move to config

func NewCredApp(credStore *store.CredStore) *CredApp {
	credapp := CredApp{store: credStore}
	return &credapp
}

// hashed, err := hashPw([]byte(pass))
// if err != nil {
// 	return err
// }
// log.Println("hashed: ", hashed)

func (ca *CredApp) CheckLogin(user, pass, useragent, userip string) (*response.Tokens, error) {

	stored, err := ca.store.GetCred(user)
	if err != nil {
		// log.Printf("err: %s \n", err)
		return nil, err
	}
	// check if password correct
	err = compareHash(stored.Pass, []byte(pass))
	if err != nil {
		return nil, err
	}

	// then, create access token and refresh token
	td, err := createToken(user)
	if err != nil {
		return nil, err
	}

	// store tokens info to memory/redis/cache
	ai := &common.AuthInfo{}
	ai.UserAgent = useragent
	ai.IP = userip
	ai.AuthTime = time.Now()
	ai.ExpTime = td.atExpires

	memstore.AddAccessAuth(user, td.accessUuid, ai)
	memstore.AddRefreshAuth(user, td.refreshUuid, td.rtExpires)

	result := &response.Tokens{}
	result.AccessToken = td.accessToken
	result.Exp = td.atExpires
	result.RefreshToken = td.refreshToken

	return result, nil
}

func (ca *CredApp) CreateRefreshToken(claim *common.RefreshClaim, userAgent, userIp string) (*response.Tokens, error) {

	if _, err := memstore.CheckAccs(claim.Name, claim.UUIDA); err != nil {
		return nil, errors.New("cannot refresh, user already logged out")
	}

	// then, create new access token and refresh token
	td, err := createToken(claim.Name)
	if err != nil {
		return nil, err
	}

	// store tokens to memory, delete the old one
	// ai := createAuth(td, userAgent, userIp)
	ai := &common.AuthInfo{}
	ai.UserAgent = userAgent
	ai.IP = userIp
	ai.AuthTime = time.Now()
	ai.ExpTime = td.atExpires

	if err = memstore.DeleteAccs(claim.Name, claim.UUIDA); err != nil {
		log.Print("[Cred][App][CreateRefreshToken] Err: ", err.Error())
	}
	if err = memstore.DeleteRefr(claim.Name, claim.UUIDR); err != nil {
		log.Print("[Cred][App][CreateRefreshToken] Err: ", err.Error())
	}
	memstore.AddAccessAuth(claim.Name, td.accessUuid, ai)
	memstore.AddRefreshAuth(claim.Name, td.refreshUuid, td.rtExpires)

	result := &response.Tokens{}
	result.AccessToken = td.accessToken
	result.Exp = td.atExpires
	result.RefreshToken = td.refreshToken
	return result, nil
}

func (ca *CredApp) Logout(username, uuid string) error {
	if err := memstore.DeleteAccs(username, uuid); err != nil {
		log.Print("[Cred][App][CreateRefreshToken] Err: ", err.Error())
		return err
	}
	return nil
}

func hashPw(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func compareHash(hashed string, plain []byte) error {
	byteHash := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(byteHash, plain)
	if err != nil {
		return err
	}

	return nil
}

func createToken(username string) (*tokenDetails, error) {
	result_td := &tokenDetails{}
	// td.user = user
	result_td.atExpires = time.Now().Add(time.Minute * 10).Unix() //todo: move to config
	// td.rtExpires = time.Now().Add(time.Hour * 12).Unix()
	result_td.rtExpires = time.Now().Add(time.Minute * 30).Unix() //todo: move to config

	result_td.accessUuid = uuid.New().String()
	result_td.refreshUuid = uuid.New().String()

	aclaims := common.AccessClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: result_td.atExpires,
		},
		// UserID:   id,
		Name: username,
		UUID: result_td.accessUuid,
	}
	authToken, err := jwtutil.CreateToken(aclaims, "access", result_td.atExpires)
	if err != nil {
		return nil, err
	}

	rclaims := common.RefreshClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: result_td.rtExpires,
		},
		// UserID:   id,
		Name:  username,
		UUIDR: result_td.refreshUuid,
		UUIDA: result_td.accessUuid,
	}

	refToken, err := jwtutil.CreateToken(rclaims, "refresh", result_td.rtExpires)
	if err != nil {
		return nil, err
	}

	result_td.accessToken = authToken
	result_td.refreshToken = refToken

	return result_td, nil
}

// func createAuth(td *tokenDetails, agent, ip string) *common.AuthInfo {
// 	ai := &common.AuthInfo{}
// 	ai.UserAgent = agent
// 	ai.IP = ip
// 	ai.AuthTime = time.Now()
// 	ai.ExpTime = td.atExpires
// 	return ai
// }
