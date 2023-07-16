package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func padKey(key []byte, blockSize int) []byte {
	if len(key) >= blockSize {
		return key[:blockSize]
	}
	paddedKey := make([]byte, blockSize)
	copy(paddedKey, key)
	return paddedKey
}

// Encrypts a string using AES encryption.
func Encrypt(plaintext string) string {
	key := []byte("ThisIsARandomKey1234") // 16-byte key
	block, _ := aes.NewCipher(padKey(key, aes.BlockSize))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	result := append(nonce, ciphertext...)
	return base64.URLEncoding.EncodeToString(result)
}

// Decrypts a string using AES decryption.
func Decrypt(ciphertext string) string {
	key := []byte("ThisIsARandomKey1234") // 16-byte key
	ciphertextBytes, _ := base64.URLEncoding.DecodeString(ciphertext)
	block, _ := aes.NewCipher(padKey(key, aes.BlockSize))
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce := ciphertextBytes[:nonceSize]
	ciphertextBytes = ciphertextBytes[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertextBytes, nil)
	return string(plaintext)
}
