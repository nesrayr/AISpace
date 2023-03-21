package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesrayr/database"
	"github.com/nesrayr/models"
	"time"
)

func Home(c *fiber.Ctx) error {
	laboratories := []models.Laboratory{}
	articles := []models.Article{}
	database.DB.Db.Find(&laboratories)
	database.DB.Db.Find(&articles)

	return c.Render("index", fiber.Map{
		"Title":        "AISpace",
		"Subtitle1":    "Лаборатория",
		"Subtitle2":    "Статьи",
		"Laboratories": laboratories,
		"Articles":     articles,
	})
}

func NewLaboratoryView(c *fiber.Ctx) error {
	return c.Render("laboratory/new", fiber.Map{
		"Title": "New info",
	})
}

func NewArticleView(c *fiber.Ctx) error {
	return c.Render("article/new", fiber.Map{
		"Title": "New article",
	})
}

func CreateLaboratory(c *fiber.Ctx) error {
	laboratory := new(models.Laboratory)
	if err := c.BodyParser(laboratory); err != nil {
		return NewLaboratoryView(c)
	}
	result := database.DB.Db.Create(&laboratory)
	if result.Error != nil {
		return NewLaboratoryView(c)
	}

	return Home(c)
}

func CreateArticle(c *fiber.Ctx) error {
	article := new(models.Article)
	if err := c.BodyParser(article); err != nil {
		return NewArticleView(c)
	}
	article.Edited = time.Now()
	result := database.DB.Db.Create(&article)
	if result.Error != nil {
		return NewArticleView(c)
	}
	return Home(c)
}

func ShowArticle(c *fiber.Ctx) error {
	article := models.Article{}
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).First(&article)
	if result.Error != nil {
		return NotFound(c)
	}

	return c.Render("article/show", fiber.Map{
		"Title":   "Article",
		"Article": article,
	})
}

func ShowLaboratory(c *fiber.Ctx) error {
	laboratory := models.Laboratory{}
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).First(&laboratory)
	if result.Error != nil {
		return NotFound(c)
	}

	return c.Render("laboratory/show", fiber.Map{
		"Title":      "Laboratory",
		"Laboratory": laboratory,
	})
}

func EditArticle(c *fiber.Ctx) error {
	article := models.Article{}
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).First(&article)
	if result.Error != nil {
		return NotFound(c)
	}

	return c.Render("article/edit", fiber.Map{
		"Title":   "Edit article",
		"Article": article,
	})
}

func UpdateArticle(c *fiber.Ctx) error {
	article := new(models.Article)
	id := c.Params("id")

	if err := c.BodyParser(article); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}
	article.Edited = time.Now()
	result := database.DB.Db.Model(&article).Where("id=?", id).Updates(article)
	if result.Error != nil {
		return EditArticle(c)
	}

	return c.Render("article/show", fiber.Map{
		"Title":   "Article",
		"Article": article,
	})
}

func EditLaboratory(c *fiber.Ctx) error {
	laboratory := models.Laboratory{}
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).First(&laboratory)
	if result.Error != nil {
		return NotFound(c)
	}

	return c.Render("laboratory/edit", fiber.Map{
		"Title":      "Edit laboratory",
		"Laboratory": laboratory,
	})
}

func UpdateLaboratory(c *fiber.Ctx) error {
	laboratory := new(models.Laboratory)
	id := c.Params("id")

	if err := c.BodyParser(laboratory); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}

	result := database.DB.Db.Model(&laboratory).Where("id=?", id).Updates(laboratory)
	if result.Error != nil {
		return EditLaboratory(c)
	}

	return c.Render("laboratory/show", fiber.Map{
		"Title":      "Laboratory",
		"Laboratory": laboratory,
	})
}

func DeleteArticle(c *fiber.Ctx) error {
	article := models.Article{}
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).Delete(&article)
	if result.Error != nil {
		return NotFound(c)
	}

	return Home(c)
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendFile("views/404.html")
}
