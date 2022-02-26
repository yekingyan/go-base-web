package dao

import (
	"context"
	"fmt"
	mgotesting "gService/share/mongo/testing"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var mongoURI string

func getAuthMongo(ctx context.Context, t *testing.T) *AuthMongo {
	fmt.Println("mc:", mongoURI)
	mc, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoURI),
		options.Client().SetConnectTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatal("failed to connect to mongo:", err)
	}
	logger, _ := zap.NewDevelopment()
	m := NewMongo(mc.Database("gservice"), logger)
	return m
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	m := getAuthMongo(ctx, t)

	users := []*UserRow{
		{
			Username: "test1",
			Password: "test1",
		},
		{
			Username: "test2",
			Password: "test2",
		},
	}

	for _, user := range users {
		t.Run(user.Username, func(t *testing.T) {
			hp := user.Password // orginal password for test
			ok, row, err := m.CreateUser(user.Username, hp)
			if !ok || err != nil {
				t.Fatal("failed to create user:", err)
			}
			if row.Username != user.Username {
				t.Errorf("username mismatch: %s != %s", row.Username, user.Username)
			}
		})
	}

	fmt.Println("repeat create user test")
	user := users[0]
	t.Run(user.Username, func(t *testing.T) {
		hp := user.Password + "ee"
		ok, row, err := m.CreateUser(user.Username, hp)
		if ok {
			fmt.Println("row:", row)
			t.Error("rea create user:", err)
		}
	})
}

func TestMain(m *testing.M) {
	code := func() int {
		var clear mgotesting.Clear
		mongoURI, clear = mgotesting.RunDockerMongo()
		defer clear()
		return m.Run()
	}()
	os.Exit(code)
}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	m := getAuthMongo(ctx, t)
	ok, u1, err := m.CreateUser("test1", "test1")
	if u1 == (UserRow{}) {
		t.Fatal("failed to create user:", err)
	}
	if !ok {
		fmt.Println("use the old user")
	}
	u2, err := m.GetUserByID(u1.ID)
	if err != nil {
		t.Fatal("failed to get user:", err)
	}
	if u1 != u2 {
		t.Error("user mismatch:", u1, u2)
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	mc, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:27017"),
		options.Client().SetConnectTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatal("failed to connect to mongo:", err)
	}
	col := mc.Database("gservice").Collection("users")
	// r, err := col.InsertOne(ctx, bson.M{
	// 	"username": "update",
	// 	"d1": bson.M{
	// 		"d11": "d11",
	// 	},
	// 	"d2": bson.M{
	// 		"d21": "d21",
	// 	},
	// })
	// fmt.Println("insert result:", r, err)
	// r, err := col.UpdateOne(ctx,
	// 	bson.M{"username": "update"},
	// 	bson.M{
	// 		"$set": bson.M{
	// 			"d1.d11": "d11 2 d111",
	// 		}})
	// fmt.Println("update result:", r, err)
	r, err := col.UpdateOne(ctx,
		bson.M{"username": "update"},
		bson.M{
			"$unset": bson.M{
				"d1.d11": 1,
			}})
	fmt.Println("update result:", r, err)
}
