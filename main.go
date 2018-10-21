package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"./models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	asession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/middleware"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	Subject = "dedgar contact form submission"
	CharSet = "UTF-8"
)

var (
	defaultCost, _    = strconv.Atoi(os.Getenv("DEFAULT_COST"))
	sender            = os.Getenv("ADMIN_EMAIL")
	recipient         = os.Getenv("ADMIN_EMAIL")
	host              = os.Getenv("POSTGRESQL_SERVICE_HOST")
	port              = os.Getenv("POSTGRESQL_SERVICE_PORT")
	dbUser            = os.Getenv("POSTGRESQL_USER")
	dbPass            = os.Getenv("POSTGRESQL_PASSWORD")
	dbName            = os.Getenv("POSTGRESQL_DATABASE")
	certAcc           = os.Getenv("CERT_ACC")
	postMap           = make(map[string]string)
	psqlInfo          = fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, dbUser, dbPass, dbName)
	db, _             = gorm.Open("postgres", psqlInfo)
	g_id, g_key       = getOauth("/home/codemaya/ansible/google_auth_creds")
	oauthStateString  = "random"
	googleOauthConfig = &oauth2.Config{
		ClientID:     g_id,
		ClientSecret: g_key,
		RedirectURL:  "http://127.0.0.1:8080/oauth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

type Contact struct {
	Name    string //`json:"name" form:"name"`
	Email   string //`json:"email" form:"email"`
	Message string //`json:"message" form:"message"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func handleGoogleCallback(c echo.Context) error {
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

	fmt.Println("accessToken", token.AccessToken)

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("error getting response")
		fmt.Println(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(contents))
	return c.String(200, string(contents)+`
*** validate token: https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=`+token.AccessToken)
}

func handleGoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// only return true if the url maps to a file in our specific hierarchy
// can be replaced with a
func availableVids(show string, season string, episode string) bool {
	if _, err := os.Stat("./static/vid/" + show + "/" + season + "/" + episode + ".mp4"); err == nil {
		return true
	}
	return false
}

// GET /
func getMain(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", postMap)
}

// handle any error by attempting to render a custom page for it
func custom404Handler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.Render(code, errorPage, code); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func getCert(c echo.Context) error {
	response := c.Param("response")
	return c.String(http.StatusOK, response+"."+certAcc)
}

// GET /about
func getAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", nil)
}

// GET /contact
func getContact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

// GET /login
func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// GET /login
func getRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

// GET /privacy
func getPrivacy(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy.html", nil)
}

// GET /dev
func getDev(c echo.Context) error {
	return c.Render(http.StatusOK, "dev.html", nil)
}

// GET /graph
func getGraph(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if _, ok := sess.Values["current_user"].(string); ok {
		return c.Render(http.StatusOK, "graph_a.html", nil)
	}
	return c.Redirect(http.StatusPermanentRedirect, "/login")
}

// POST /post-contact
func postContact(c echo.Context) error {

	if strings.Contains(c.FormValue("message"), "http") && strings.Contains(c.FormValue("message"), "dedgar.com/") == false {
		return c.String(http.StatusOK, "Form submitted")
	}

	TextBody := c.FormValue("name") + "\n" + c.FormValue("email") + "\n" + c.FormValue("message")

	sess, err := asession.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(sender),
	}

	result, err := svc.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

	}
	fmt.Println(c.FormValue("name"))
	fmt.Println(c.FormValue("email"))
	fmt.Println(c.FormValue("message"))
	fmt.Println("Email Sent to address: " + recipient)
	fmt.Println(result)
	return c.String(http.StatusOK, "Form submitted")
}

// GET /post/:postname
func getPost(c echo.Context) error {
	post := c.Param("postname")
	if _, ok := postMap[post]; ok {
		return c.Render(http.StatusOK, post+".html", post)
	}
	return c.Render(http.StatusNotFound, "e04.html", "404 Post not found")
}

// GET /post
func getPostView(c echo.Context) error {
	return c.Render(http.StatusOK, "post_view.html", postMap)
}

func findSummary(fpath string) string {
	file, err := os.Open(fpath + "_summary")
	if err != nil {
		return "No summary"
	}
	defer file.Close()

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		buffer.WriteString(line)
		//		if line == "<!--more-->" {
		//			break
		//		}
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}

// Populates a map of postnames that gets checked every call to GET /post/:postname.
// We're running in a container, so populating this on startup works fine as we won't be adding
// any new posts while the container is running.
func findPosts(dirpath string, extension string) map[string]string {
	if err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		}
		if strings.HasSuffix(path, extension) {
			postname := strings.Split(path, extension)[0]
			summary := findSummary(postname)
			//fmt.Println(summary)
			//fmt.Println(fmt.Sprintf("%T", summary))
			postMap[filepath.Base(postname)] = summary
		}
		return err
	}); err != nil {
		panic(err)
	}
	return postMap
}

func getOauth(filepath string) (id, key string) {
	filebytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	file_str := string(filebytes)

	id, key = strings.Split(file_str, "\n")[0], strings.Split(file_str, "\n")[1]

	return id, key
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
	db.NewRecord(new_user)
	db.Create(&new_user)
}

// POST /register
func postRegister(c echo.Context) error {
	TextBody := c.FormValue("login") + "\n" + c.FormValue("password")
	fmt.Println(TextBody)

	if userFound(c.FormValue("username")) || emailFound(c.FormValue("email")) {
		return c.String(http.StatusOK, "Email address or username already taken, try again!")
	}

	createUser(c.FormValue("email"), c.FormValue("username"), c.FormValue("password"))

	return c.Redirect(http.StatusPermanentRedirect, "/login")
}

func emailFound(eName string) bool {
	var user models.User
	var found_e models.User

	db.Where(&models.User{Email: eName}).First(&user).Scan(&found_e)

	if found_e.Email != "" {
		fmt.Printf("%s already taken!", found_e.Email)
		return true
	}

	fmt.Printf("%s not taken!", found_e.Email)
	return false
}

func userFound(uName string) bool {
	var user models.User
	var found_u models.User

	db.Where(&models.User{UName: uName}).First(&user).Scan(&found_u)

	if found_u.UName != "" {
		fmt.Println("Username found.")
		return true
	}

	fmt.Println("Username not found.")
	return false
}

func compareLogin(uName, pWord string) bool {
	var user models.User
	var found_u models.User

	db.Where(&models.User{UName: uName}).First(&user).Scan(&found_u)

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

// POST /login
func postLogin(c echo.Context) error {
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

func checkDB() {
	if !db.HasTable(&models.User{}) {
		fmt.Println("Creating users table")
		db.CreateTable(&models.User{})
	}
}

func ServerHeader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
			mainCookie(c)
			fmt.Println("serverheader /admin")
			sess, _ := session.Get("session", c)

			if sess.Values["authenticated"] == "true" {
				fmt.Println("in if block")
				fmt.Println(sess.Values)
				return next(c)
			}

			return next(c)
		}
	}
}

func getTrial(c echo.Context) error {
	sess, _ := session.Get("session", c)
	logged_in_dude := sess.Values["current_user"].(string)
	return c.String(http.StatusOK, logged_in_dude)
}

func ServerTet(http.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("servertest /auth")
			//c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
			return next(c)
		}
	}
}

func mainCookie(c echo.Context) { //error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["foo"] = "bar"
	sess.Values["authenticated"] = "true"
	sess.Save(c.Request(), c.Response())
}

func main() {
	t := &Template{
		templates: func() *template.Template {
			tmpl := template.New("")
			if err := filepath.Walk("./tmpl", func(path string, info os.FileInfo, err error) error {
				if strings.HasSuffix(path, ".html") {
					_, err = tmpl.ParseFiles(path)
					if err != nil {
						log.Println(err)
					}
				}
				return err
			}); err != nil {
				panic(err)
			}
			return tmpl
		}(),
	}
	e := echo.New()
	e.Static("/", "static")
	e.Renderer = t
	//e.HTTPErrorHandler = custom404Handler
	//	e.Pre(middleware.HTTPSWWWRedirect())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//admin_group := e.Group("/posts", ServerHeader())
	//admin_group.Use(ServerHeader())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET")))))
	e.Use(ServerHeader())

	go checkDB()

	//		return req.RemoteAddr == "127.0.0.1" || (currentUser.(*models.User) != nil && currentUser.(*models.User).Role == "admin")

	findPosts("./tmpl/posts", ".html")
	//fmt.Println(findPosts("./tmpl/posts", ".html"))
	e.GET("/", getMain)
	e.POST("/", getMain)
	e.GET("/about", getAbout)
	e.GET("/register", getRegister)
	e.POST("/register", postRegister)
	e.GET("/login", getLogin)
	e.POST("/login", postLogin)
	e.GET("/login/google", handleGoogleLogin)
	e.GET("/oauth/callback", handleGoogleCallback)
	e.GET("/about-us", getAbout)
	e.GET("/trial", getTrial)
	e.GET("/graph", getGraph)
	e.GET("/contact", getContact)
	e.GET("/contact-us", getContact)
	e.GET("/privacy-policy", getPrivacy)
	e.GET("/privacy", getPrivacy)
	e.GET("/dev", getDev)
	e.POST("/post-contact", postContact)
	e.GET("/post", getPostView)
	e.GET("/post/", getPostView)
	e.GET("/posts", getPostView)
	e.GET("/posts/", getPostView)
	e.GET("/post/:postname", getPost)
	e.GET("/posts/:postname", getPost)
	e.GET("/.well-known/acme-challenge/test", getCert)
	e.GET("/.well-known/acme-challenge/test/", getCert)
	e.GET("/.well-known/acme-challenge/:response", getCert)
	e.GET("/.well-known/acme-challenge/:response/", getCert)
	e.GET("/well-known/acme-challenge/:response", getCert)
	e.GET("/well-known/acme-challenge/:response/", getCert)
	e.File("/robots.txt", "static/public/robots.txt")
	e.File("/sitemap.xml", "static/public/sitemap.xml")
	e.Logger.Info(e.Start(":8080"))
	//	e.Logger.Info(e.StartAutoTLS(":443"))
}
