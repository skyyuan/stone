package auth

import (
	"crypto/sha1"
	"encoding/hex"

	"stone/common"
)

// AppAuth application authentication
type AppAuth struct {
	common.Model
	AppKey      string `json:"app_key" gorm:"unique_index"`
	Secret      string `json:"secret"`
	AppName     string `json:"app_name" gorm:"index"`
	Description string `json:"description"`
}

// Verify verify checksum
func (auth *AppAuth) Verify(nonce, timestamp, checksum string) bool {
	input := auth.Secret + nonce + timestamp
	out := sha1.Sum([]byte(input))
	hexout := hex.EncodeToString(out[:len(out)])

	return hexout == checksum
}

// AppVerify verify checksum with appkey
func AppVerify(appKey, nonce, timestamp, checksum string) (bool, error) {
	db := common.DBBegin()
	defer db.DBRollback()
	var auth AppAuth
	if err := db.Where(&AppAuth{AppKey: appKey}).First(&auth).Error; err != nil {
		return false, err
	}
	return auth.Verify(nonce, timestamp, checksum), nil
}

// GenerateAppAuth 生成一个App的AppAuth
func GenerateAppAuth(appName, desc string) (*AppAuth, error) {
	appkey, secret := common.GenKeyPair()
	auth := &AppAuth{
		AppKey:      appkey,
		Secret:      secret,
		AppName:     appName,
		Description: desc,
	}
	db := common.DBBegin()
	defer db.DBRollback()
	if err := db.Create(auth).Error; err != nil {
		return nil, err
	}
	db.DBCommit()
	return auth, nil
}
