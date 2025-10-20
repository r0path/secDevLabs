package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"crypto/subtle"

	"golang.org/x/crypto/pbkdf2"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateSalt(byteNumber int) ([]byte, error) {
	salt, err := generateRandomBytes(byteNumber)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// Hash returns the hashed password and the salt used.
func Hash(password string) (string, string, error) {
	salt, err := generateSalt(32)
	if err != nil {
		return "", "", err
	}

	hashedPassword := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	// Encode hash and salt using base64 to safely store as text
	hashB64 := base64.StdEncoding.EncodeToString(hashedPassword)
	saltB64 := base64.StdEncoding.EncodeToString(salt)

	return hashB64, saltB64, nil
}

// VerifyHash returns true if the received password matches the hash with the given salt.
func VerifyHash(password string, hashedPassword string, salt string) bool {
	// Decode salt and expected hash from base64. If decoding fails, fall back to raw bytes
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		// Fallback to raw bytes from older storage format
		saltBytes = []byte(salt)
	}
	expectedHash, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		// Fallback to raw bytes from older storage format
		expectedHash = []byte(hashedPassword)
	}

	newHash := pbkdf2.Key([]byte(password), saltBytes, 4096, 32, sha256.New)

	// Compare the computed hash with the expected hash in constant time
	if len(newHash) != len(expectedHash) {
		return false
	}
	if subtle.ConstantTimeCompare(newHash, expectedHash) == 1 {
		return true
	}
	return false
}
