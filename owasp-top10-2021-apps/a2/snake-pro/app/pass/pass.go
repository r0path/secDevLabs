package pass

import (
	"crypto/subtle"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// CheckPass checks a password.
// If truePassword appears to be a bcrypt hash it will use bcrypt.CompareHashAndPassword.
// Otherwise it falls back to a constant-time comparison to preserve backward compatibility.
func CheckPass(truePassword, attemptPassword string) bool {
	// Detect bcrypt hash prefixes ($2a$, $2b$, $2y$ are common)
	if strings.HasPrefix(truePassword, "$2a$") || strings.HasPrefix(truePassword, "$2b$") || strings.HasPrefix(truePassword, "$2y$") {
		err := bcrypt.CompareHashAndPassword([]byte(truePassword), []byte(attemptPassword))
		return err == nil
	}

	a := []byte(truePassword)
	b := []byte(attemptPassword)

	if len(a) != len(b) {
		// Pad both to the same length to ensure constant-time comparison and avoid leaking length information
		max := len(a)
		if len(b) > max {
			max = len(b)
		}
		pa := make([]byte, max)
		pb := make([]byte, max)
		copy(pa, a)
		copy(pb, b)
		return subtle.ConstantTimeCompare(pa, pb) == 1 && len(a) == len(b)
	}

	return subtle.ConstantTimeCompare(a, b) == 1
}
