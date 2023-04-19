package main

import (
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	_ "github.com/nesrayr/cmd/docs"
	"github.com/nesrayr/database"
	"log"
)

// @title AISpace
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
// @schemes https
func main() {
	database.ConnectDB()
	//database.DB.Db.Exec("DROP TABLE moderators")
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Static("/", "./public")
	app.Get("/swagger/*", swagger.HandlerDefault)
	SetupRoutes(app)

	err := app.ListenTLS(":3000", "server.crt", "server.key")
	if err != nil {
		fmt.Println(err)
		log.Fatal(2)
	}
}
