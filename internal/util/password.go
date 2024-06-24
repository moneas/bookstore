package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func HashPassword(password string) (string, error) {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func CheckPasswordHash(password, hash string) bool {
	hashedPassword, _ := HashPassword(password)
	return hashedPassword == hash
}

func DecodeBase64(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
