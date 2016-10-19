package main

import "fmt"
import "time"
import "github.com/dgrijalva/jwt-go"

var hmacSecret = []byte("54f398ac6538dab757b6ce2bc206b1e947cd0d63")

func NewToken(username string) string {
	var now int64 = time.Now().Unix()
	const oneDay int64 = 1000 * 60 * 60 * 2 // Accept for 2 hours
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// https://en.wikipedia.org/wiki/JSON_Web_Token#Standard_fields
	"iss": "Silo",
    "sub": "Authorization",
    "aud": username,
    "exp": now + oneDay,
    "nbf": now,
    "iat": now,
    // "jti": not needed, as we don't have different issuers
	})

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return tokenString
}

// https://godoc.org/github.com/dgrijalva/jwt-go#example-New--Hmac
func VerifyToken(tokenString string) (map[string]interface{}, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return nil, false
	}
}