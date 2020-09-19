package store

import (
	"log"
	"pbk-main/pkg/db"
	"pbk-main/pkg/model/dto"
	"pbk-main/pkg/store/query"
)

// type tokenKey struct {
// 	user, uuid string
// }

// type authStoreMap = map[tokenKey]*common.AuthInfo // only store refresh token info

type CredStore struct {
	dbImpl *db.DBImpl
	// authenticatedStorage authStoreMap
}

func NewCredentialStore(dbImpl *db.DBImpl) *CredStore {
	cs := &CredStore{}
	cs.dbImpl = dbImpl
	// cs.authenticatedStorage = make(authStoreMap)
	return cs
}

func (c *CredStore) GetCred(user string) (*dto.Credential, error) {
	conn, err := c.dbImpl.Connect()
	if err != nil {
		log.Printf("[Store][Credential] Can't open db connection. %+v \n", err)
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.Prepare(query.GetCred)
	if err != nil {
		return nil, err
	}
	cr := &dto.Credential{}
	err = stmt.QueryRow(user).Scan(&cr.ID, &cr.User, &cr.Pass)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

// func (c *CredStore) AddAuth(user, uuid string, ai *common.AuthInfo) {
// 	c.authenticatedStorage[tokenKey{user, uuid}] = ai

// 	c.logStore()
// }

// func (c *CredStore) CheckAuth(user, uuid string) (*common.AuthInfo, error) {
// 	info, ok := c.authenticatedStorage[tokenKey{user, uuid}]
// 	if !ok {
// 		return nil, errors.New("auth info not found")
// 	}
// 	return info, nil
// }

// func (c *CredStore) DeleteAuth(user, uuid string) error {
// 	if _, ok := c.authenticatedStorage[tokenKey{user, uuid}]; !ok {
// 		return errors.New("auth info not found, cannot delete")
// 	}
// 	delete(c.authenticatedStorage, tokenKey{user, uuid})

// 	c.logStore()
// 	return nil
// }

// func (c *CredStore) logStore() {
// 	log.Println("auth store: ", c.authenticatedStorage)
// }
