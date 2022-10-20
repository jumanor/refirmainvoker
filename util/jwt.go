package util

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var sampleSecretKey = []byte("SecretYouShouldHide")

// Creamos Token
func GenerarJWT() (string, error) {

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		Issuer:    "jumanor",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(sampleSecretKey)
}

// Verificamos Token
func VerificarJWT(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if token.Valid {
		return nil
	} else {
		return err
	}

}
