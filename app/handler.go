package app

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type Handlers struct {
	userService services.UserService
	authService services.AuthService
}

func (ch *Handlers) getAllUsers(cookie string) (*[]domain.User, error) {
	var auth domain.Authentication
	_, err := auth.IsLoggedIn(cookie)

	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	users, err := ch.userService.GetAllUsers()
	if err != nil {
		log.Panicf("error: %v", err)
	}
	return users, nil
}

func (ch *Handlers) CreateUser(user domain.User) error {
	err := ch.userService.CreateUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (ch *Handlers) GetUserByID(cookie string, id primitive.ObjectID) (*domain.User, error){
	var auth domain.Authentication
	_, err := auth.IsLoggedIn(cookie)

	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := ch.userService.GetUserByID(id)

	if err != nil {
		return nil, fmt.Errorf("error...")
	}

	return user, err
}

func (ch *Handlers) UpdateUser(id primitive.ObjectID, user domain.User, cookie string) error {
	var auth domain.Authentication
	_, err := auth.IsLoggedIn(cookie)

	if err != nil {
		return err
	}

	_, err = ch.userService.UpdateUser(id, &user)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

func (ch *Handlers) DeleteByID(cookie string, id primitive.ObjectID) error {
	var auth domain.Authentication
	_, err := auth.IsLoggedIn(cookie)
	if err != nil {
		return err
	}

	err = ch.userService.DeleteByID(id)

	if err != nil {
		return err
	}

	return err
}

func (ch *Handlers) Login( email string, password string, c *fiber.Ctx) (*domain.User, error) {
	var auth domain.Authentication

	u, token, err := ch.authService.Login(email, password)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	signedToken := make([]byte, 0, 100)
	signedToken = append(signedToken, []byte("Bearer " + token + "|")...)
	t, err := auth.SignToken([]byte(token))

	signedToken = append(signedToken, t...)

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = string(signedToken)
	cookie.Expires = time.Now().Add(24 * time.Hour)

	// Set cookie
	c.Cookie(cookie)

	return u, nil
}