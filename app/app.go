package app

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"example.com/app/services"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"log"

	"github.com/gofiber/fiber/v2"
)

func Start() {
	// wiring everything up
	ch := Handlers{userService: services.NewUserService(repo.NewUserRepoImpl()),
		authService: services.NewAuthService(repo.NewAuthRepoImpl())}

	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		u, err := ch.getAllUsers(cookie, c)
		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(201).JSON(u)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		id := c.Params("id")

		newId, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		u, err := ch.GetUserByID(cookie, newId, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(200).JSON(u)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(domain.User)

		err := c.BodyParser(user)

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		err = ch.CreateUser(*user, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(201).SendString("Success...")
	})

	app.Post("/users/login", func(c *fiber.Ctx) error {
		details := new(domain.LoginDetails)

		err := c.BodyParser(details)

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		_, err = ch.Login(details.Email, details.Password, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(200).SendString("Logged in")
	})

	app.Post("/users/:id", func(c *fiber.Ctx) error {
		user := new(domain.User)

		err := c.BodyParser(user)

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		cookie := c.Cookies("session")

		id := c.Params("id")

		newId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		err = ch.UpdateUser(newId, *user, cookie, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(200).SendString("Success...")
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		id , err := primitive.ObjectIDFromHex(c.Params("id"))

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		err = ch.DeleteByID(cookie, id, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(204).SendString("Success...")
	})

	log.Fatal(app.Listen(":8080"))
}


