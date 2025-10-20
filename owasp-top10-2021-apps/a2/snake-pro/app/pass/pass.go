package pass

import "golang.org/x/crypto/bcrypt"

// CheckPass checks a password using bcrypt comparison
func CheckPass(storedHash, attemptPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(attemptPassword))
	return err == nil
}
