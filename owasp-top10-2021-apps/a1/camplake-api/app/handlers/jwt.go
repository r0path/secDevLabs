package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"camp-lake-api/types"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

func CreateToken(creds types.Credentials) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &types.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
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
		return claims, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing algorithm is HMAC (reject "alg": "none" and other algorithms)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}

	return claims, nil
}
