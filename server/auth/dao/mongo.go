package dao

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuthMongo is anth data access object.
type AuthMongo struct {
	userCol    *mongo.Collection
	sessionCol *mongo.Collection
}

// UserField is user collection field.
type UserField string

// user col field.
const (
	IDField       UserField = "_id"
	UsernameField UserField = "username"
	PasswordField UserField = "password"
)

// change UserField to String.
func (f UserField) String() string {
	return string(f)
}

// UserRow is a row of user collection.
type UserRow struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// NewMongo returns a AuthMongo.
func NewMongo(db *mongo.Database) *AuthMongo {
	return &AuthMongo{
		userCol:    db.Collection("user"),
		sessionCol: db.Collection("session"),
	}
}

// CreateUser insert vaild user to user collection.
func (m *AuthMongo) CreateUser(username, hashPassword string) (bool, UserRow, error) {
	c := context.Background()
	res := m.userCol.FindOneAndUpdate(
		c,
		&bson.M{
			UsernameField.String(): username,
		},
		&bson.M{
			"$setOnInsert": bson.M{
				UsernameField.String(): username,
				PasswordField.String(): hashPassword,
			},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	row := UserRow{}
	if err := res.Err(); err != nil {
		return false, row, err
	}
	err := res.Decode(&row)
	if err != nil {
		fmt.Println("Decode err:", err)
		return false, row, err
	}
	if row.Password != hashPassword {
		return false, row, fmt.Errorf("already exist user: %s", username)
	}
	return true, row, nil
}

// GetUserByName returns a user by username.
func (m *AuthMongo) GetUserByName(username string) (UserRow, error) {
	c := context.Background()
	res := m.userCol.FindOne(
		c,
		&bson.M{
			UsernameField.String(): username,
		},
	)
	row := UserRow{}
	if err := res.Err(); err != nil {
		return row, err
	}
	err := res.Decode(&row)
	if err != nil {
		return row, err
	}
	return row, nil
}
