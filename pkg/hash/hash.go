package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
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

//XorDecrypt is a simple XOR decryption function
func XorDecrypt(hexStr string) string {
	result := ""
	key := "this is a secret key."
	for i := 0; i < len(hexStr); i += 2 {
		byteStr := hexStr[i : i+2]
		byteVal, _ := hex.DecodeString(byteStr)
		charCode := int(byteVal[0]) ^ int(key[i/2%len(key)])
		result += string(charCode)
	}
	return result
}
