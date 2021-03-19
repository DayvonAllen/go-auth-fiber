package app

import (
	"example.com/app/app/helpers"
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

		err := helpers.IsLoggedIn(cookie, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		u, err := ch.getAllUsers(c)
		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(201).JSON(u)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		err := helpers.IsLoggedIn(cookie, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		id, err := primitive.ObjectIDFromHex(c.Params("id"))

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		u, err := ch.GetUserByID(id, c)

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

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		err := helpers.IsLoggedIn(cookie, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		id , err := primitive.ObjectIDFromHex(c.Params("id"))

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		user := new(domain.User)

		err = c.BodyParser(user)

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		err = ch.UpdateUser(id, *user, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(200).SendString("Success...")
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		err := helpers.IsLoggedIn(cookie, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		id , err := primitive.ObjectIDFromHex(c.Params("id"))

		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("%v", err))
		}

		err = ch.DeleteByID(id, c)

		if err != nil {
			return c.SendString(fmt.Sprintf("%v", err))
		}

		return c.Status(204).SendString("Success...")
	})

	log.Fatal(app.Listen(":8080"))
}


