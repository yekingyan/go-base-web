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

// CreateSession returns a session for userID
func (s *SessionToken) CreateSession(userID string) (dao.SessionRow, error) {
	return s.Mongo.CreateSession(userID, s.Expire)
}
