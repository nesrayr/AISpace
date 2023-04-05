package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
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

func AuthMain(c *fiber.Ctx) error {
	url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func AuthHome(c *fiber.Ctx) error {
	laboratories := []models.Laboratory{}
	articles := []models.Article{}
	logos := []models.Logo{}
	database.DB.Db.Find(&laboratories)
	database.DB.Db.Find(&articles)
	database.DB.Db.Find(&logos)

	return c.Render("auth_index", fiber.Map{
		"Laboratories": laboratories,
		"Articles":     articles,
		"Logos":        logos,
	})
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
		return fiber.NewError(http.StatusBadRequest, "Failed to save refresh token to Redis")
	}
	jwtToken, err := generateJWT(token)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"token": jwtToken,
	})
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

func AuthShowArticle(c *fiber.Ctx) error {
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

	return c.Render("article/auth_show", fiber.Map{
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

func GetAdmin(c *fiber.Ctx) error {
	users := []models.User{}
	database.DB.Db.Find(&users)

	return c.Status(200).JSON(users)
}

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

func NotModerator(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendFile("views/not_moderator.html")
}
