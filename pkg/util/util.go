package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/logger"
	"io"
)

func Hash(value string) string {
	hash := sha1.New()
	hash.Write([]byte(value))
	hashByte := hash.Sum(nil)
	return hex.EncodeToString(hashByte)
}

func Encrypt(plainText, key string) string {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.Err("%v", err)
		return ""
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger.Err("%v", err)
		return ""
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Err("%v", err)
		return ""
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.RawStdEncoding.EncodeToString(cipherText)
}

func Decrypt(encryptedText, key string) string {
	cipherText, _ := base64.RawStdEncoding.DecodeString(encryptedText)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.Err("%v", err)
		return ""
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger.Err("%v", err)
		return ""
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return ""
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		logger.Err("%v", err)
		return ""
	}
	return string(plaintext)
}

func InterfacePtrToInt(input interface{}) int {
	if val, ok := input.(*int); ok {
		return *val
	}
	return 0
}

func InterfaceToInt(input interface{}) int {
	if val, ok := input.(int); ok {
		return val
	}
	return 0
}

func InterfacePtrToString(input interface{}) string {
	if val, ok := input.(*string); ok {
		return *val
	}
	return ""
}

func InterfaceToString(input interface{}) string {
	if val, ok := input.(string); ok {
		return val
	}
	return ""
}
