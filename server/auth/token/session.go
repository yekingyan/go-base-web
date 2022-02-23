package token

import (
	"gService/auth/dao"

	"go.uber.org/zap"
)

// SessionToken is a token for authentication
type SessionToken struct {
	Expire int64
	Mongo  *dao.AuthMongo
	Logger *zap.Logger
}

// NewSessionToken return a session creator
func NewSessionToken(expire int64, mongo *dao.AuthMongo, logger *zap.Logger) *SessionToken {
	return &SessionToken{
		Expire: expire,
		Mongo:  mongo,
		Logger: logger,
	}
}

// GenerateToken returns a session for userID
func (s *SessionToken) GenerateToken(userID string) (string, int64, error) {
	row, err := s.Mongo.CreateSession(userID, s.Expire)
	return row.ID, row.ExpireTime, err
}

// GetExpiresIn returns the expire time of the session
func (s *SessionToken) GetExpiresIn() int64 {
	return s.Expire
}