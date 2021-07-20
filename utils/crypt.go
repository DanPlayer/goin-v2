package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func MD5(file []byte) string {
	crypto := md5.New()
	crypto.Write(file)
	return hex.EncodeToString(crypto.Sum(nil))
}

func Md5hex(str string) string {
	ret := md5.Sum([]byte(str))
	return hex.EncodeToString(ret[:])
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// AesEncrypt AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// ContentEncrypt 加密原始数据（Base64+AES）
func ContentEncrypt(origData, key string) string {
	token, err := AesEncrypt([]byte(origData), []byte(key))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(token)
}

// ContentDecrypt 解密数据（Base64+AES）
func ContentDecrypt(token, key string) (string, error) {
	tokenStr, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	content, err := AesDecrypt([]byte(tokenStr), []byte(key))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
