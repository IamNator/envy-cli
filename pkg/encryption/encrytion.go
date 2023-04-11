package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func MakeEncrypter(secret string) func(string) (string, error) {
	return func(value string) (string, error) {
		return encryptionAdapter(value, secret)
	}
}

func MakeDecrypter(secret string) func(string) (string, error) {
	return func(value string) (string, error) {
		return decryptionAdapter(value, secret)
	}
}

const (
	SaltLength   = 16
	KeyLength    = 32 // 256-bit key length
	IterationNum = 100000
)

// GenerateKey generates a key by using a password and a random salt
func generateKey(password string, salt []byte) ([]byte, error) {
	key := pbkdf2.Key([]byte(password), salt, IterationNum, KeyLength, sha256.New)
	return key, nil
}

// Encrypt encrypts the value using the secret key
func encryptionAdapter(value string, secret string) (string, error) {

	// Generate a random salt
	salt := make([]byte, SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key, err := generateKey(secret, salt)
	if err != nil {
		return "", err
	}

	result, err := encrypt([]byte(value), key)
	if err != nil {
		return "", err
	}

	// Prepend the salt to the ciphertext
	result = append(salt, result...)

	encoded := base64.StdEncoding.EncodeToString(result)

	return encoded, nil
}

// Decrypt decrypts the cipher text using the secret key
func decryptionAdapter(cipherText string, secret string) (string, error) {

	ciphertext, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < SaltLength {
		return "", errors.New("ciphertext too short")
	}

	// Extract the salt from the ciphertext
	salt := ciphertext[:SaltLength]
	ciphertext = ciphertext[SaltLength:]

	key, err := generateKey(secret, salt)
	if err != nil {
		return "", err
	}

	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		return "", err
	}

	// Output: some plaintext data
	return string(plaintext), nil
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
