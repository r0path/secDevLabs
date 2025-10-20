package services

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Recovery bool   `json:"admin"`
	jwt.StandardClaims
}

func GenerateJwt(login string, recovery bool) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	// We include a unique jti (JWT ID) for single-use tracking and a short expiration for recovery tokens.
	// Recovery tokens are sensitive, so set a short lifetime (e.g., 15 minutes) and include a random jti.
	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", err
	}
	jti := hex.EncodeToString(jtiBytes)

	claims := &JwtCustomClaims{
		Name:     login,
		Recovery: recovery,
		StandardClaims: jwt.StandardClaims{
			Id:        jti,
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	// Note: Server-side single-use enforcement is required to fully mitigate reuse; this change
	// adds a jti so the application can track and revoke tokens if backed by persistent storage.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
