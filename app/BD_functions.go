package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

func getUser(id uuid.UUID) (User, error) {
	if rand.Intn(100) > 50 {
		return User{}, errors.New("unlucky this time")
	} else {
		return User{name: "Lucky"}, nil
	}
}

func addUser(user User) error {
	fmt.Printf("[addUser] %v\n", user)
	return nil
}
