package application

import (
	"pbk-main/pkg/lib/core"
	"pbk-main/pkg/model/common"
	"pbk-main/pkg/model/dto"
	"pbk-main/pkg/model/request"
	"pbk-main/pkg/store"
)

type AccountApp struct {
	store *store.AccountStore
}

func NewAccountApp(acStore *store.AccountStore) *AccountApp {
	acApp := AccountApp{store: acStore}
	return &acApp
}

func (acc *AccountApp) GetAccountById(id int) (*dto.Account, error) {
	data, err := acc.store.GetAccountById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (acc *AccountApp) DeleteAccount(id int) (int64, error) {
	rows, err := acc.store.DeleteAccount(id)
	if err != nil {
		return -1, err
	}
	return rows, nil
}

func (acc *AccountApp) CreateAccount(req *request.CreateAccSBReq, shouldCreateSecret bool) (int64, error) {
	newAcc := &dto.Account{}
	accCrypt := &common.AccountCrypto{}
	newAcc.Username = req.Username
	newAcc.EmailPhonenum = req.EmailPhonenum
	newAcc.Platform = req.Platform
	newAcc.Details = req.Details
	newAcc.CredID = req.CredID //TODO: change source from param not req

	secrets := make(map[string]string)

	if shouldCreateSecret {
		accCrypt.Username = req.Username
		accCrypt.Platform = req.Platform
		accCrypt.SecretPass = req.SecretPass
		accCrypt.SecretInfo = req.SecretInfo

		var err error
		secrets, err = core.ConstructSecretPass(accCrypt, &req.Key)
		if err != nil {
			return -1, err
		}
	}

	catIds := []int{1} // TODO: change this maybe
	if req.CategoryID != nil {
		catIds = req.CategoryID
	}

	insertedId, err := acc.store.InsertNewAccount(newAcc, catIds)
	if err != nil {
		return -1, err
	}

	if shouldCreateSecret {
		newSb := &dto.SafetyBox{}
		newSb.SecretPass = secrets["sp"]
		newSb.SecretInfo = secrets["si"]
		newSb.AccID = int(insertedId)

		err := acc.store.InsertNewSafetyBox(newSb)
		if err != nil {
			return -1, err
		}
	}

	return insertedId, nil
}

func (acc *AccountApp) UpdateAccount(id int, req *request.UpdateSBReq) error {

	_, err := acc.store.GetAccountById(id)
	if err != nil {
		return err
	}

	newSb := &dto.Account{}
	newSb.Username = req.Username
	newSb.EmailPhonenum = req.Email
	newSb.Platform = *req.Platform
	newSb.Details = req.Details
	newSb.CredID = *req.CredID

	err = acc.store.UpdateAccount(id, newSb)
	if err != nil {
		return err
	}

	return nil
}
