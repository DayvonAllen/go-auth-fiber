package repo

import (
	"context"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
	*domain.User
}

func(a AuthRepoImpl) Login(email, password string) (*domain.User, string, error) {
	var login domain.Authentication
	var user domain.User
	opts := options.FindOne()
	_ = dbConnection.Collection.FindOne(context.TODO(), bson.D{{"email", email}},opts).Decode(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err == nil {
		token, _ := login.GenerateJWT(user)
		return &user, token, nil
	}

	return nil, "", fmt.Errorf("error logging in")
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}