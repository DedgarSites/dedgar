package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dedgarsites/dedgar/datastores"
	"github.com/dedgarsites/dedgar/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	defaultCost, _    = strconv.Atoi(os.Getenv("DEFAULT_COST"))
	oauthStateString  = "random" // TODO randomize
	googleOauthConfig = &oauth2.Config{
		ClientID:     datastores.OAuthID,
		ClientSecret: datastores.OAuthKey,
		RedirectURL:  "https://dedgar.com/oauth/callback", //"http://127.0.0.1:8080/oauth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

// HandleGoogleCallback listens on
// GET /oauth/callback
func HandleGoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	code := c.QueryParam("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("error getting response")
		fmt.Println(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
		fmt.Println(err)
	}

	var gUser models.GoogleUser
	err = json.Unmarshal(contents, &gUser)
	if err != nil {
		fmt.Println("Error Unmarshaling google user json: ", err)
	}

	if ok := datastores.AuthMap[gUser.Email]; ok {
		sess, _ := session.Get("session", c)
		sess.Values["authenticated"] = "true"
		sess.Values["google_logged_in"] = gUser.Email
		sess.Save(c.Request(), c.Response())

		return c.Render(http.StatusOK, "dashboard.html", nil)
	}
	return c.String(200, string(contents)+`https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=`+token.AccessToken)
}

// GET /login/google
func HandleGoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// POST /login
func PostLogin(c echo.Context) error {
	if !userFound(c.FormValue("username")) {
		return c.String(http.StatusOK, "Username not found!")
	}

	if compareLogin(c.FormValue("username"), c.FormValue("password")) {
		sess, _ := session.Get("session", c)
		sess.Values["current_user"] = c.FormValue("username")
		sess.Values["logged_in"] = "true"
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusPermanentRedirect, "/")
	}

	return c.Render(http.StatusUnauthorized, "404.html", "401 not authenticated")
}

// POST /register
func PostRegister(c echo.Context) error {
	TextBody := c.FormValue("login") + "\n" + c.FormValue("password")
	fmt.Println(TextBody)

	if userFound(c.FormValue("username")) || emailFound(c.FormValue("email")) {
		return c.String(http.StatusOK, "Email address or username already taken, try again!")
	}

	createUser(c.FormValue("email"), c.FormValue("username"), c.FormValue("password"))

	return c.Redirect(http.StatusPermanentRedirect, "/login")
}

func HashPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	return string(bytes), err
}

func createUser(eName, uName, pWord string) {
	hashed_pw, err := HashPass(pWord)

	if err != nil {
		log.Fatal(err)
	}

	new_user := models.User{Email: eName, UName: uName, Password: hashed_pw}
	datastores.DB.NewRecord(new_user)
	datastores.DB.Create(&new_user)
}

func compareLogin(uName, pWord string) bool {
	var user models.User
	var found_u models.User

	datastores.DB.Where(&models.User{UName: uName}).First(&user).Scan(&found_u)

	if found_u.UName == "" {
		fmt.Println("Invalid username or password!")
		return false
	}

	hashedPW := found_u.Password

	err := bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(pWord))

	if err != nil {
		fmt.Println("Invalid username or password!")
		fmt.Println(err)
		return false
	}

	fmt.Println("Found login combo matched!")
	return true
}

func userFound(uName string) bool {
	var user models.User
	var found_u models.User

	datastores.DB.Where(&models.User{UName: uName}).First(&user).Scan(&found_u)

	if found_u.UName != "" {
		fmt.Println("Username found.")
		return true
	}

	fmt.Println("Username not found.")
	return false
}

func emailFound(eName string) bool {
	var user models.User
	var found_e models.User

	datastores.DB.Where(&models.User{Email: eName}).First(&user).Scan(&found_e)

	if found_e.Email != "" {
		fmt.Printf("%s already taken!", found_e.Email)
		return true
	}

	fmt.Printf("%s not taken!", found_e.Email)
	return false
}
