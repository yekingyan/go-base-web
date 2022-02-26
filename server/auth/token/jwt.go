package token

import (
	"crypto/rsa"
	"gService/share/id"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTToken generate a jwt token
type JWTToken struct {
	expire     int64
	nowFunc    func() time.Time
	privateKey *rsa.PrivateKey
}

// NewJWTToken return a jwt token creator
func NewJWTToken(expire int64, nowFunc func() time.Time, privateKey *rsa.PrivateKey) *JWTToken {
	return &JWTToken{
		expire:     expire,
		nowFunc:    nowFunc,
		privateKey: privateKey,
	}
}

// GetExpiresIn returns the expire time of the session
func (j *JWTToken) GetExpiresIn() int64 {
	return j.expire
}

// GenerateToken return a jwt token
func (j *JWTToken) GenerateToken(userID id.UserID) (string, int64, error) {
	now := j.nowFunc().Unix()
	expire := now + j.expire
	jtk := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		ExpiresAt: expire,
		IssuedAt:  now,
		Issuer:    "gService",
		Subject:   userID.String(),
	})
	tk, err := jtk.SignedString(j.privateKey)
	return tk, expire, err
}
