package auth

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := []string{
		"12345678",
		"asdfghjkl.*?",
		"12323asdf",
		"123asdf%^&%",
		"2  3  4  5  6",
		"花木成畦手自栽",
		"@工a1.asdf.0",
	}
	for _, p := range password {
		t.Run(p, func(t *testing.T) {
			h, err := HashPassword(p)
			if err != nil {
				t.Error(err)
			}
			pass := CheckPasswordHash(p, h)
			fmt.Println(p, h, pass)
			fmt.Println("cost:", GetHashingCost([]byte(h)))
			if !pass {
				t.Errorf("password hash failed")
			}
		})
	}
}
