package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/nesrayr/database"
)

//type GoogleUserInfo struct {
//	Sub           string `json:"sub"`
//	Name          string `json:"name"`
//	GivenName     string `json:"given_name"`
//	FamilyName    string `json:"family_name"`
//	Picture       string `json:"picture"`
//	Email         string `json:"email"`
//	EmailVerified bool   `json:"email_verified"`
//}

func main() {
	database.ConnectDB()
	//database.DB.Db.Exec("TRUNCATE TABLE photos RESTART IDENTITY ")
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	//oauth2Config := &oauth2.Config{
	//	ClientID:     os.Getenv("CLIENT_ID"),
	//	ClientSecret: os.Getenv("CLIENT_SECRET"),
	//	RedirectURL:  "http://localhost:3000/auth/google/callback",
	//	Scopes: []string{
	//		"https://www.googleapis.com/auth/userinfo.profile",
	//		"https://www.googleapis.com/auth/userinfo.email",
	//	},
	//	Endpoint: google.Endpoint,
	//}

	setupRoutes(app)
	//app.Get("/auth", func(c *fiber.Ctx) error {
	//	url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	//	return c.Redirect(url)
	//})
	//
	//app.Get("/auth/google/callback", func(c *fiber.Ctx) error {
	//	code := c.Query("code")
	//	token, err := oauth2Config.Exchange(context.Background(), code)
	//	if err != nil {
	//		return err
	//	}
	//	client := oauth2Config.Client(context.Background(), token)
	//	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	//	if err != nil {
	//		return err
	//	}
	//	defer resp.Body.Close()
	//
	//	var userInfo GoogleUserInfo
	//	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
	//		return err
	//	}
	//
	//	fmt.Fprintf(c, "Имя пользователя: %s\n", userInfo.Name)
	//	fmt.Fprintf(c, "Email: %s\n", userInfo.Email)
	//	return nil
	//})

	app.Static("/", "./public")

	app.Listen(":3000")
}
