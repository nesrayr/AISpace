package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nesrayr/database"
	"github.com/nesrayr/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
	"time"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

var (
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "https://localhost:3000/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

func Home(c *fiber.Ctx) error {
	var role string
	if r, ok := c.Locals("user_role").(string); ok {
		role = r
	} else {
		role = "User"
	}
	laboratories := []models.Laboratory{}
	articles := []models.Article{}
	logos := []models.Logo{}
	database.DB.Db.Find(&laboratories)
	database.DB.Db.Find(&articles)
	database.DB.Db.Find(&logos)
	fmt.Println(role)
	return c.Render("index", fiber.Map{
		"Laboratories": laboratories,
		"Articles":     articles,
		"Logos":        logos,
		"Role":         role,
	})
}

func AuthMain(c *fiber.Ctx) error {
	url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Println(url)
	err := c.Redirect(url)
	if err != nil {
		fmt.Println("Error while redirecting:", err)
	}
	return nil
}

func AuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return fiber.NewError(http.StatusBadRequest, "Missing authorization code")
	}
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Failed to exchange code for token")
	}
	if err := database.DB.RedisClient.Set(context.Background(), "refresh_token", token.RefreshToken,
		time.Duration(time.Now().Add(time.Hour*24*3).Unix())).Err(); err != nil {
		fmt.Println("Error saving refresh token to Redis:", err)
		return fiber.NewError(http.StatusBadRequest, "Failed to save refresh token to Redis")
	}
	jwtToken, err := generateJWT(token)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
	})
	//fmt.Println("c.Locals('jwt'):", c.Locals("jwt"))
	return c.Redirect("/")
}

type Claims struct {
	Sub       string `json:"sub"`
	Username  string `json:"username"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Role      string `json:"role"`
}

func (c Claims) Valid() error {
	if c.ExpiresAt < time.Now().Unix() {
		return jwt.ErrSignatureInvalid
	}
	return nil
}

func generateJWT(token *oauth2.Token) (string, error) {
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}
	users := []models.User{}
	var role string
	res := database.DB.Db.Where("email=?", userInfo.Email).First(&users)
	if res.Error != nil {
		role = "User"
	} else {
		role = users[0].Role
	}

	claims := &Claims{
		Sub:       userInfo.Sub,
		Username:  userInfo.Name,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(),
		Role:      role,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
		Secure:   true,
	})

	c.Locals("user_role", nil)

	return c.Redirect("/")
}

// NewLaboratoryView func get the new laboratory form
// @Description Get the new laboratory form
// @Summary get the new laboratory form
// @Tags Laboratory
// @Accept json
// @Produce html
// @Success 200
// @Router /laboratory/new [get]
func NewLaboratoryView(c *fiber.Ctx) error {
	return c.Render("laboratory/new", fiber.Map{
		"Title": "New info",
	})
}

// NewArticleView func get the new article form
// @Description Get the new article form
// @Summary get the new article form
// @Tags Article
// @Accept json
// @Produce html
// @Success 200
// @Router /article/new [get]
func NewArticleView(c *fiber.Ctx) error {
	return c.Render("article/new", fiber.Map{
		"Title": "New article",
	})
}

// CreateLaboratory func for creates new laboratory
// @Summary Create a new laboratory
// @Description Create a new laboratory
// @Tags Laboratory
// @Accept  json
// @Produce  html
// @Param info body string true "Info"
// @Success 200
// @Router /laboratory [post]
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

// CreateArticle func for creates new article
// @Summary Create a new article
// @Description Create a new article
// @Tags Article
// @Accept  json
// @Produce  html
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Success 200
// @Router /article [post]
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

// ShowArticle func gets article by given ID or 404 error.
// @Description Get article by given ID.
// @Summary get article by given ID
// @Tags Article
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200
// @Router /article/{id} [get]
func ShowArticle(c *fiber.Ctx) error {
	var role string
	if r, ok := c.Locals("user_role").(string); ok {
		role = r
	} else {
		role = "User"
	}
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
		"Role":    role,
	})
}

// ShowLaboratory func gets laboratory by given ID or 404 error.
// @Description Get laboratory by given ID.
// @Summary get laboratory by given ID
// @Tags Laboratory
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200
// @Router /laboratory/{id} [get]
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

// EditArticle func returns for of updating article by given ID.
// @Description Form of updating article.
// @Summary Form of updating article
// @Tags Article
// @Accept json
// @Produce json
// @Param id body string true "Article ID"
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Success 200 {string} status "ok"
// @Router /article/{id}/edit [patch]
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

// UpdateArticle func for updates article by given ID.
// @Description Update article.
// @Summary update article
// @Tags Article
// @Accept json
// @Produce json
// @Param id body string true "Article ID"
// @Param title body string true "Title"
// @Param description body string true "Description"
// @Success 201 {string} status "ok"
// @Router /article/{id} [patch]
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

// EditLaboratory func returns for of updating laboratory by given ID.
// @Description Form of updating laboratory.
// @Summary Form of updating article
// @Tags Laboratory
// @Accept json
// @Produce json
// @Param id body string true "Article ID"
// @Param info body string true "Info"
// @Success 200 {string} status "ok"
// @Router /laboratory/{id}/edit [patch]
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

// UpdateLaboratory func for updates laboratory by given ID.
// @Description Update laboratory.
// @Summary update laboratory
// @Tags Laboratory
// @Accept json
// @Produce json
// @Param id body string true "Laboratory ID"
// @Param info body string true "Info"
// @Success 201 {string} status "ok"
// @Router /laboratory/{id} [patch]
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

// DeleteArticle func for deletes article by given ID.
// @Description Delete article by given ID.
// @Summary delete article by given ID
// @Tags Article
// @Accept json
// @Produce json
// @Param id body string true "Article ID"
// @Success 204 {string} status "ok"
// @Router /article/{id} [delete]
func DeleteArticle(c *fiber.Ctx) error {
	article := models.Article{}
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).Delete(&article)
	if result.Error != nil {
		return NotFound(c)
	}

	return Home(c)
}

// CreateImage func for creates a new image.
// @Description Create a new image.
// @Summary create a new image
// @Tags Image
// @Accept json
// @Produce json
// @Param data body string true "Data"
// @Success 200
// @Router /image/new [post]
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

// GetImages func gets all exists images.
// @Description Get all exists images.
// @Summary get all exists images
// @Tags Image
// @Accept json
// @Produce json
// @Success 200
// @Router /images [get]
func GetImages(c *fiber.Ctx) error {
	images := []models.Photo{}
	database.DB.Db.Find(&images)

	return c.Status(200).JSON(images)
}

// CreateLogo func for creates a new logo.
// @Description Create a new logo.
// @Summary create a new logo
// @Tags Logo
// @Accept json
// @Produce json
// @Param data body string true "Data"
// @Success 200
// @Router /logo/new [post]
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

// GetLogos func gets all exists logos.
// @Description Get all exists logos.
// @Summary get all exists logos
// @Tags Logo
// @Accept json
// @Produce json
// @Success 200
// @Router /logos [get]
func GetLogos(c *fiber.Ctx) error {
	logos := []models.Logo{}
	database.DB.Db.Find(&logos)

	return c.Status(200).JSON(logos)
}

// DeleteLogo func for deletes logo by given ID.
// @Description Delete logo by given ID.
// @Summary delete logo by given ID
// @Tags Logo
// @Accept json
// @Produce json
// @Param id body string true "Logo ID"
// @Success 204 {string} status "ok"
// @Router /logo/{id} [delete]
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

// DeleteImage func for deletes image by given ID.
// @Description Delete image by given ID.
// @Summary delete image by given ID
// @Tags Image
// @Accept json
// @Produce json
// @Param id body string true "Image ID"
// @Success 204 {string} status "ok"
// @Router /image/{id} [delete]
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

// CreateAdmin func for creates a new admin.
// @Description Create a new admin.
// @Summary create a new admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Success 200
// @Router /user/new [post]
func CreateAdmin(c *fiber.Ctx) error {
	admin := new(models.User)
	if err := c.BodyParser(admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	database.DB.Db.Create(&admin)
	return c.Status(200).JSON(admin)
}

// GetAdmins func gets all exists users.
// @Description Get all exists users.
// @Summary get all exists users
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
// @Router /users [get]
func GetAdmins(c *fiber.Ctx) error {
	users := []models.User{}
	database.DB.Db.Find(&users)

	return c.Status(200).JSON(users)
}

// DeleteAdmin func for deletes admin by given ID.
// @Description Delete admin by given ID.
// @Summary delete admin by given ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param id body string true "Admin ID"
// @Success 204 {string} status "ok"
// @Router /user/{id} [delete]
func DeleteAdmin(c *fiber.Ctx) error {
	user := new(models.User)
	id := c.Params("id")

	result := database.DB.Db.Where("id=?", id).Delete(&user)
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

//func NotModerator(c *fiber.Ctx) error {
//	return c.Status(fiber.StatusNotFound).SendFile("views/not_moderator.html")
//}
