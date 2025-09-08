package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"go.uber.org/zap"
	"math/big"
	"strings"
	"sync"
)

var (
	EncryptImplVar *EncryptImpl
	once           sync.Once
)

type EncryptImpl struct {
}

func NewEncryptImpl() *EncryptImpl {
	once.Do(func() {
		EncryptImplVar = &EncryptImpl{}
	})
	return EncryptImplVar
}

func (e *EncryptImpl) RandomString(length int) (string, error) {
	chars := "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"
	var b strings.Builder
	for i := 0; i < length; i++ {
		nBig, err := crand.Int(crand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b.WriteByte(chars[nBig.Int64()])
	}
	return b.String(), nil
}

func (e *EncryptImpl) Pkcs7Pad(src []byte, blockSize int) []byte {
	padLen := blockSize - (len(src) % blockSize)
	return append(src, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
}

func (e *EncryptImpl) AesEncryptCBCBase64(plaintext, key, iv string) (string, error) {
	keyBytes := []byte(key)
	ivBytes := []byte(iv)
	plainBytes := []byte(plaintext)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	if len(ivBytes) != block.BlockSize() {
		return "", fmt.Errorf("invalid IV size: %d", len(ivBytes))
	}

	padded := e.Pkcs7Pad(plainBytes, block.BlockSize())
	encrypted := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, ivBytes)
	mode.CryptBlocks(encrypted, padded)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (e *EncryptImpl) EncryptPassword(password, salt string) string {
	salt = strings.TrimSpace(salt)
	if salt == "" {
		return password
	}

	randomPrefix, err := e.RandomString(64)
	if err != nil {
		zap.L().Error("生成随机前缀失败", zap.Error(err))
		return password
	}
	iv, err := e.RandomString(16)
	if err != nil {
		zap.L().Error("生成IV失败", zap.Error(err))
		return password
	}

	combined := randomPrefix + password
	enc, err := e.AesEncryptCBCBase64(combined, salt, iv)
	if err != nil {
		zap.L().Error("AES加密失败，回退为明文密码", zap.Error(err))
		return password
	}
	return enc
}
