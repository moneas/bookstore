package util

import (
	"crypto/md5"
	"encoding/hex"
)

func HashPassword(password string) (string, error) {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
