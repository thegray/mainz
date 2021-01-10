package application

import (
	"encoding/json"
	"log"
	"pbk-main/pkg/model/dto"
	"pbk-main/pkg/model/request"
	"pbk-main/pkg/store"
	"pbk-main/pkg/util/cryptoutil"
	"reflect"
)

type sptype struct {
	Un string
	Pw string
	Pf string
}

type SafetyBoxApp struct {
	store *store.SafetyBoxStore
}

func NewSafetyBoxApp(sbStore *store.SafetyBoxStore) *SafetyBoxApp {
	sbapp := SafetyBoxApp{store: sbStore}
	return &sbapp
}

func (sba *SafetyBoxApp) GetSafetyBoxById(id int) (*dto.SafetyBox, error) {
	data, err := sba.store.GetSafetyBoxById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (sba *SafetyBoxApp) DeleteSafetyBox(id int) (int64, error) {
	rows, err := sba.store.DeleteSafetyBox(id)
	if err != nil {
		return -1, err
	}
	return rows, nil
}

func (sba *SafetyBoxApp) CreateSafetyBox(req *request.CreateSBReq) (int64, error) {
	newSb := &dto.SafetyBox{}
	newSb.Username = req.Username
	newSb.SecretPass = req.SecretPass //this one should be required
	key := *req.Key
	newSb.Email = req.Email
	newSb.Platform = *req.Platform
	newSb.Details = req.Details
	newSb.SecretInfo = req.SecretInfo
	newSb.CredID = *req.CredID
	catIds := []int{1}
	if req.CategoryID != nil {
		catIds = req.CategoryID
	}
	spToEnc, err := constructSecretPass(newSb.Username, newSb.SecretPass, newSb.Platform)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	cipherSP, err := cryptoutil.EncryptAES(key, spToEnc)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	cipherSI, err := cryptoutil.EncryptAES(key, newSb.SecretInfo)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	newSb.SecretPass = cipherSP
	newSb.SecretInfo = cipherSI

	// plainSP, err := cryptoutil.DecryptAES(key, cipherSP)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	// log.Println("plainSP: ", plainSP)
	// plainSI, err := cryptoutil.DecryptAES(key, cipherSI)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	// log.Println("plainSI: ", plainSI)

	id, err := sba.store.InsertNewSafetyBox(newSb, catIds)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (sba *SafetyBoxApp) UpdateSafetyBox(id int, req *request.UpdateSBReq) error {

	_, err := sba.store.GetSafetyBoxById(id)
	if err != nil {
		return err
	}

	newSb := &dto.SafetyBox{}
	newSb.Username = req.Username
	newSb.Email = req.Email
	newSb.Platform = *req.Platform
	newSb.Details = req.Details
	newSb.CredID = *req.CredID

	err = sba.store.UpdateSafetyBox(id, newSb)
	if err != nil {
		return err
	}

	return nil
}

func (sba *SafetyBoxApp) UpdateSecret(id int, req *request.UpdateSecretReq) error {

	foundItem, err := sba.store.GetSafetyBoxById(id)
	if err != nil {
		return err
	}

	spToEnc, err := constructSecretPass(foundItem.Username, req.SecretPass, foundItem.Platform)
	if err != nil {
		log.Println(err)
		return err
	}
	cipherSP, err := cryptoutil.EncryptAES(*req.Key, spToEnc)
	if err != nil {
		log.Println(err)
		return err
	}
	cipherSI, err := cryptoutil.EncryptAES(*req.Key, req.SecretInfo)
	if err != nil {
		log.Println(err)
		return err
	}

	err = sba.store.UpdateSecret(id, cipherSP, cipherSI)
	if err != nil {
		return err
	}

	return nil
}

func (sba *SafetyBoxApp) RevealSecret(req *request.RevealSecretReq) (string, string, error) {
	foundItem, err := sba.store.GetSafetyBoxById(*req.ID)
	if err != nil {
		return "", "", err
	}

	plainSP, err := cryptoutil.DecryptAES(*req.Key, foundItem.SecretPass)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	sp := &sptype{}
	err = json.Unmarshal([]byte(plainSP), sp)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	plainSI, err := cryptoutil.DecryptAES(*req.Key, foundItem.SecretInfo)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return sp.Pw, plainSI, nil
}

func constructSecretPass(username, secretpass, platform string) (string, error) {
	// format := "\"un\"=\"%s\"&\"pw\"=\"%s\"&\"pf\"=\"%s\""
	// format := "un=%s&pw=%s&pf=%s"
	// format := "{\"un\":\"%s\",\"pw\":\"%s\",\"pf\":\"%s\"}"
	sp := &sptype{Un: username, Pw: secretpass, Pf: platform}
	out, err := json.Marshal(sp)
	if err != nil {
		return "", err
	}

	full := string(out)
	// full := fmt.Sprintf(format, username, secretpass, platform)
	return full, nil
}

func getStringValue(p interface{}) string {
	if reflect.ValueOf(p).Kind() == reflect.Ptr {
		r := p.(*string)
		return *r
	}
	return p.(string)
}
