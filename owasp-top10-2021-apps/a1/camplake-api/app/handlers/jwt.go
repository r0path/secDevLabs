package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"camp-lake-api/types"

	"github.com/dgrijalva/jwt-go"
)

func getJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "camplake-secret-key"
	}
	return []byte(secret)
}

func CreateToken(creds types.Credentials) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &types.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		log.Fatal(err)
	}
	return tokenString, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func TokenValid(r *http.Request) (types.Claims, error) {
	claims := types.Claims{}

	tokenString := ExtractToken(r)
	if tokenString == "" {
		log.Println("No token!")
		return claims, errors.New("no token provided")
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method: only HS256 is accepted")
		}
		return getJWTKey(), nil
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, errors.New("invalid token")
	}

	return claims, nil
}
