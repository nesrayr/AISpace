package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/nesrayr/database"
)

func main() {
	database.ConnectDB()
	//database.DB.Db.Exec("DROP TABLE moderators")
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	setupRoutes(app)

	app.Static("/", "./public")

	app.Listen(":3000")
}
