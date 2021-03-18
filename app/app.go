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

		u, err := ch.getAllUsers(cookie)
		if err != nil {
			c.Status(400)
			return err
		}
		c.Status(200)
		return c.JSON(u)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		id := c.Params("id")

		newId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		u, err := ch.GetUserByID(cookie, newId)
		if err != nil {
			c.Status(400)
			return fmt.Errorf("error...%w",  err)
		}
		c.Status(200)
		return c.JSON(u)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(domain.User)

		if err := c.BodyParser(user); err != nil {
			return err
		}

		err := ch.CreateUser(*user)
		for err != nil {
			c.Status(400)
			return fmt.Errorf("error...%w",  err)
		}

		c.Status(201)
		return c.SendString("Success...")
	})

	app.Post("/users/login", func(c *fiber.Ctx) error {
		details := new(domain.LoginDetails)

		if err := c.BodyParser(details); err != nil {
			return err
		}

		_, err := ch.Login(details.Email, details.Password, c)
		for err != nil {
			c.Status(401)
			return fmt.Errorf("error...%w",  err)
		}

		c.Status(200)
		return c.SendString("Logged in")
	})

	app.Post("/users/:id", func(c *fiber.Ctx) error {
		user := new(domain.User)

		if err := c.BodyParser(user); err != nil {
			return err
		}

		cookie := c.Cookies("session")

		id := c.Params("id")

		newId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		err = ch.UpdateUser(newId, *user, cookie)
		for err != nil {
			c.Status(500)
			return fmt.Errorf("error...")
		}

		c.Status(200)
		return c.SendString("Success...")
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		cookie := c.Cookies("session")

		id := c.Params("id")

		newId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		err = ch.DeleteByID(cookie, newId)
		if err != nil {
			c.Status(500)
			return err
		}

		c.Status(204)
		return c.SendString("Success...")
	})

	log.Fatal(app.Listen(":8080"))
}


