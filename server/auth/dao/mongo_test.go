package dao

import (
	"context"
	"fmt"
	mgotesting "gService/share/mongo/testing"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var mongoURI string

func TestRegister(t *testing.T) {
	ctx := context.Background()
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
