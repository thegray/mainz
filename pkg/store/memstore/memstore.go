package memstore

import (
	"errors"
	"log"
	"pbk-main/pkg/model/common"
)

// type tokenKey struct {
// 	user, uuid string
// }

type accsStoreMap = map[string]uuidStoreMap
type uuidStoreMap = map[string]*common.AuthInfo

type refrStoreMap = map[string]map[string]int64

var accsStorage accsStoreMap
var refrStorage refrStoreMap

func init() {
	accsStorage = make(accsStoreMap)
	refrStorage = make(refrStoreMap)
}

func GetAccsStore() *accsStoreMap {
	return &accsStorage
}

func GetRefrStore() *refrStoreMap {
	return &refrStorage
}

func AddAccessAuth(user, uuid string, ai *common.AuthInfo) {
	uuidMap, ok := accsStorage[user]
	if ok {
		uuidMap[uuid] = ai
		return
	}

	newUuidMap := make(map[string]*common.AuthInfo)
	newUuidMap[uuid] = ai
	accsStorage[user] = newUuidMap
}

func AddRefreshAuth(user, uuid string, expTime int64) {
	inner, ok := refrStorage[user]
	if !ok {
		inner = make(map[string]int64)
		refrStorage[user] = inner
	}
	inner[uuid] = expTime
}

func CheckAccs(user, uuid string) (*common.AuthInfo, error) {
	info, ok := accsStorage[user][uuid]
	if !ok {
		return nil, errors.New("access token not found")
	}
	return info, nil
}

func CheckRefr(user, uuid string) (int64, error) {
	expTime, ok := refrStorage[user][uuid]
	if !ok {
		return -1, errors.New("refresh token not found")
	}
	return expTime, nil
}

func DeleteAccs(user, uuid string) error {
	uuidMap, ok := accsStorage[user]
	if !ok {
		return errors.New("access token not found, cannot delete")
	}
	delete(uuidMap, uuid)

	// log.Println("delete auth store: ", accsStorage)
	return nil
}

func DeleteRefr(user, uuid string) error {
	inner, ok := refrStorage[user]
	if !ok {
		return errors.New("refresh token not found, cannot delete")
	}
	delete(inner, uuid)
	return nil
}

func logStore() {
	log.Println("auth store: ", accsStorage)
}
