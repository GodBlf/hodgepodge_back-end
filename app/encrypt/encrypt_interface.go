package encrypt

type Encrypt interface {
	RandomString(length int) (string, error)
	Pkcs7Pad(src []byte, blockSize int) []byte
	AesEncryptCBCBase64(plaintext, key, iv string) (string, error)
	EncryptPassword(password, salt string) string
}
