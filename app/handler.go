package app

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Handlers struct {
	userService services.UserService
	authService services.AuthService
}

func (ch *Handlers) getAllUsers(c *fiber.Ctx) (*[]domain.User, error) {
	users, err := ch.userService.GetAllUsers()

	if err != nil {
		c.Status(500)
		return nil, fmt.Errorf("server error, can't get all users: %w", err)
	}

	return users, nil
}

func (ch *Handlers) CreateUser(user domain.User, c *fiber.Ctx) error {
	err := ch.userService.CreateUser(&user)

	if err != nil {
		c.Status(500)
		return fmt.Errorf("server error, can't create user: %w", err)
	}

	return nil
}

func (ch *Handlers) GetUserByID(id primitive.ObjectID, c *fiber.Ctx) (*domain.User, error){

	user, err := ch.userService.GetUserByID(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(404)
			return nil, fmt.Errorf("server error: %w", err)
		}
		c.Status(500)
		return nil, fmt.Errorf("server error: %w", err)
	}

	return user, err
}

func (ch *Handlers) UpdateUser(id primitive.ObjectID, user domain.User, c *fiber.Ctx) error {

	_, err := ch.userService.UpdateUser(id, &user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(404)
			return fmt.Errorf("server error: %w", err)
		}
		c.Status(500)
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func (ch *Handlers) DeleteByID(id primitive.ObjectID, c *fiber.Ctx) error {

	err := ch.userService.DeleteByID(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(404)
			return fmt.Errorf("server error: %w", err)
		}
		c.Status(500)
		return fmt.Errorf("server error: %w", err)
	}

	return err
}

func (ch *Handlers) Login(email string, password string, c *fiber.Ctx) (*domain.User, error) {
	var auth domain.Authentication

	u, token, err := ch.authService.Login(email, password)

	if err != nil {
		c.Status(500)
		return nil, fmt.Errorf("server error: %w", err)
	}

	signedToken := make([]byte, 0, 100)
	signedToken = append(signedToken, []byte("Bearer " + token + "|")...)
	t, err := auth.SignToken([]byte(token))

	if err != nil {
		c.Status(500)
		return nil, fmt.Errorf("server error: %w", err)
	}

	signedToken = append(signedToken, t...)

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = string(signedToken)
	cookie.Expires = time.Now().Add(24 * time.Hour)

	// Set cookie
	c.Cookie(cookie)

	return u, nil
}