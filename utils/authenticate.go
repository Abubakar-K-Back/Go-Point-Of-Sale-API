package utils

import (
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Passcode string `json:"passcode"`
	jwt.StandardClaims
}

func SplitToken(headerToken string) string {
	parsToken := strings.SplitAfter(headerToken, " ")
	tokenString := parsToken[1]
	return tokenString
}

func AuthenticateJWT(passcode string) string {
	expirationTime := time.Now().Add(5000 * time.Minute)
	claims := &Claims{
		Passcode: passcode,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatalf("Error in Creating JWT")
	}
	return tokenString

}

func AuthToken(str string) error {
	_, err := jwt.ParseWithClaims(str, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})
	if err != nil {
		return err
	}
	return nil
}
