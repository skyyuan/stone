package common

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// GenKeyPair 生成appkey及secretkey
func GenKeyPair() (string, string) {
	timestamp := string(time.Now().Unix())
	// appkey =  uuid +  timestamp
	keySeed := uuid.New().String() + timestamp
	keyOut := sha1.Sum([]byte(keySeed))
	appkey := hex.EncodeToString(keyOut[:len(keyOut)])

	// secretkey = appkey + timestamp + rand str
	secSeed := appkey + timestamp + RandString(16)
	secOut := sha1.Sum([]byte(secSeed))
	secret := hex.EncodeToString(secOut[:len(secOut)])

	return appkey, secret
}
