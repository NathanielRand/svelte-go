package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func mainPage(c *fiber.Ctx) error {
	return c.Render("main", nil)
}

func main() {

	//template render engine
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine, //set as render engine
	})
	app.Static("/public", "./public")
	app.Get("/", mainPage)
	app.Listen(":8080")
}
