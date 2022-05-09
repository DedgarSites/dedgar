package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dedgarsites/dedgar/controllers"
	"github.com/dedgarsites/dedgar/datastores"
)

var (
	// Routers supplies an instance of echo to be used in the main function.
	Routers  *echo.Echo
	sitePath = os.Getenv("SITE_PATH")
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func init() {
	if sitePath == "" {
		sitePath = "."
	}
	t := &Template{
		templates: func() *template.Template {
			tmpl := template.New("")
			if err := filepath.Walk(sitePath+"/tmpl", func(path string, info os.FileInfo, err error) error {
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

	Routers = echo.New()
	Routers.Static("/", sitePath+"/static")
	Routers.Renderer = t

	Routers.Use(middleware.Logger())
	Routers.Use(middleware.Recover())
	Routers.Use(middleware.CORS())

	Routers.GET("/", controllers.GetMain)
	Routers.POST("/", controllers.GetMain)

	datastores.FindPosts(sitePath+"/tmpl/posts", ".html")

	Routers.GET("/", controllers.GetMain)
	Routers.POST("/", controllers.GetMain)
	Routers.GET("/about", controllers.GetAbout)
	Routers.GET("/about-us", controllers.GetAbout)
	Routers.GET("/contact", controllers.GetContact)
	Routers.GET("/contact-us", controllers.GetContact)
	Routers.GET("/privacy-policy", controllers.GetPrivacy)
	Routers.GET("/privacy", controllers.GetPrivacy)
	Routers.POST("/post-contact", controllers.PostContact)
	Routers.GET("/post", controllers.GetPostView)
	Routers.GET("/post/", controllers.GetPostView)
	Routers.GET("/posts", controllers.GetPostView)
	Routers.GET("/posts/", controllers.GetPostView)
	Routers.GET("/post/:postname", controllers.GetPost)
	Routers.GET("/posts/:postname", controllers.GetPost)
	Routers.File("/ads.txt", sitePath+"/static/public/ads.txt")
	Routers.File("/robots.txt", sitePath+"/static/public/robots.txt")
	Routers.File("/sitemap.xml", sitePath+"/static/public/sitemap.xml")
}
