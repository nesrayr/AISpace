package handlers

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/nesrayr/database"
	"github.com/nesrayr/models"
	"time"
)

func Home(c *fiber.Ctx) error {
	laboratories := []models.Laboratory{}
	articles := []models.Article{}
	logos := []models.Logo{}
	database.DB.Db.Find(&laboratories)
	database.DB.Db.Find(&articles)
	database.DB.Db.Find(&logos)

	return c.Render("index", fiber.Map{
		"Laboratories": laboratories,
		"Articles":     articles,
		"Logos":        logos,
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
	photos := []models.Photo{}
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).First(&article)
	if result.Error != nil {
		return NotFound(c)
	}
	images := database.DB.Db.Where("article_id=?", id).Find(&photos)
	if images.Error != nil {
		return NotFound(c)
	}

	return c.Render("article/show", fiber.Map{
		"Title":   "Article",
		"Article": article,
		"Photos":  photos,
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

func CreateImage(c *fiber.Ctx) error {
	image := new(models.Photo)
	if err := c.BodyParser(image); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	image.StrData = base64.StdEncoding.EncodeToString(image.Data)
	database.DB.Db.Create(&image)

	return c.Status(200).JSON(image)
}

func GetImage(c *fiber.Ctx) error {
	images := []models.Photo{}
	database.DB.Db.Find(&images)

	return c.Status(200).JSON(images)
}

func CreateLogo(c *fiber.Ctx) error {
	logo := new(models.Logo)
	if err := c.BodyParser(logo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	logo.StrData = base64.StdEncoding.EncodeToString(logo.Data)
	database.DB.Db.Create(&logo)
	return c.Status(200).JSON(logo)
}

func GetLogo(c *fiber.Ctx) error {
	logos := []models.Logo{}
	database.DB.Db.Find(&logos)

	return c.Status(200).JSON(logos)
}

func DeleteLogo(c *fiber.Ctx) error {
	logo := new(models.Logo)
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).Delete(&logo)
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
	})
}

func DeleteImage(c *fiber.Ctx) error {
	image := new(models.Photo)
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).Delete(&image)
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
	})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendFile("views/404.html")
}
