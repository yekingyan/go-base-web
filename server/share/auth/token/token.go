package sharetoken

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// JWTTokenVerifier verify a jwt token
type JWTTokenVerifier struct {
	publicKey *rsa.PublicKey
}

func (j *JWTTokenVerifier) getKeyFunc() jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	}
}

// Verify jwt token
func (j *JWTTokenVerifier) Verify(token string) (string, error) {
	jtk, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, j.getKeyFunc())
	if err != nil {
		return "", fmt.Errorf("token parse error: %v", err)
	}
	if !jtk.Valid {
		return "", fmt.Errorf("token is invalid")
	}
	clm, ok := jtk.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claims is not StandardClaims")
	}
	if err := clm.Valid(); err != nil {
		return "", fmt.Errorf("token claims is invalid: %v", err)
	}
	return clm.Subject, nil
}
