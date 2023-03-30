package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesrayr/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/", handlers.AuthHome)
	app.Get("/laboratory/new", handlers.NewLaboratoryView)
	app.Get("/article/new", handlers.NewArticleView)
	app.Post("/laboratory", handlers.CreateLaboratory)
	app.Post("/article", handlers.CreateArticle)
	app.Get("/article/:id", handlers.ShowArticle)
	app.Get("/article/auth/:id", handlers.AuthShowArticle)
	app.Get("/laboratory/:id", handlers.ShowLaboratory)
	app.Get("/article/:id/edit", handlers.EditArticle)
	app.Patch("/article/:id", handlers.UpdateArticle)
	app.Get("/laboratory/:id/edit", handlers.EditLaboratory)
	app.Patch("/laboratory/:id", handlers.UpdateLaboratory)
	app.Delete("/article/:id", handlers.DeleteArticle)

	app.Post("/image/new", handlers.CreateImage)
	app.Get("/image", handlers.GetImage)
	app.Post("/logo/new", handlers.CreateLogo)
	app.Get("/logo", handlers.GetLogo)
	app.Delete("/logo/:id", handlers.DeleteLogo)
	app.Delete("/image/:id", handlers.DeleteImage)
	app.Post("/user/new", handlers.CreateUser)
	app.Get("/user", handlers.GetUser)
	app.Delete("/user/:id", handlers.DeleteUser)

	app.Get("/auth", handlers.AuthMain)
	app.Get("/auth/google/callback", handlers.AuthCallback)
}
