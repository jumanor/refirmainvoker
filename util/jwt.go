package util

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var SECRET_KEY_JWT string

// Creamos Token
func GenerarJWT() (string, error) {

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		Issuer:    "jumanor",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(SECRET_KEY_JWT))
}

// Verificamos Token
func VerificarJWT(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY_JWT), nil
	})

	if token.Valid {
		return nil
	} else {
		return err
	}

}
