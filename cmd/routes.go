package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesrayr/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/laboratory/new", handlers.NewLaboratoryView)
	app.Get("/article/new", handlers.NewArticleView)
	app.Post("/laboratory", handlers.CreateLaboratory)
	app.Post("/article", handlers.CreateArticle)
	app.Get("/article/:id", handlers.ShowArticle)
	app.Get("/laboratory/:id", handlers.ShowLaboratory)
	app.Get("/article/:id/edit", handlers.EditArticle)
	app.Patch("/article/:id", handlers.UpdateArticle)
	app.Get("/laboratory/:id/edit", handlers.EditLaboratory)
	app.Patch("/laboratory/:id", handlers.UpdateLaboratory)
	app.Delete("/article/:id", handlers.DeleteArticle)
}
