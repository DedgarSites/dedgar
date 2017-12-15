package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	host     = os.Getenv("POSTGRESQL_SERVICE_HOST")
	port     = os.Getenv("POSTGRESQL_SERVICE_PORT")
	user     = os.Getenv("POSTGRESQL_USER")
	password = os.Getenv("POSTGRESQL_PASSWORD")
	dbname   = os.Getenv("POSTGRESQL_DATABASE")
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// only return true if the url maps to a file in our specific hierarchy
func availableVids(show string, season string, episode string) bool {
	if _, err := os.Stat("./static/vid/" + show + "/" + season + "/" + episode + ".mp4"); err == nil {
		return true
	}
	return false
}

func getMain(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", "main")
}

func getContainer(c echo.Context) error {
	return c.Render(http.StatusOK, "container.html", "container")
}

// GET /watch/:show/:season/:episode
func getShow(c echo.Context) error {
	show := c.Param("show")
	season := c.Param("season")
	episode := c.Param("episode")

	vid_list := availableVids(show, season, episode)
	if vid_list {

		return c.Render(http.StatusOK, "episode_view.html", map[string]interface{}{
			"show":    show,
			"season":  season,
			"episode": episode,
		})
	}
	return c.Render(http.StatusNotFound, "404.html", "404 Video not found")
}

func getJapanese(c echo.Context) error {
	return c.Render(http.StatusOK, "level_selection.html", "level_selection")
}

// GET /kanji/:selection/:level
func getLevel(c echo.Context) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	var sqlQuery string

	switch c.Param("selection") {
	case "grade":
		sqlQuery = "SELECT kanj, von, vkun, transl, roma, rememb, jlpt, school FROM info WHERE school = $1"
	case "jlpt":
		sqlQuery = "SELECT kanj, von, vkun, transl, roma, rememb, jlpt, school FROM info WHERE jlpt = $1"
	}
	rows, err := db.Query(sqlQuery, c.Param("level"))

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	//entry := make(map[string]map[string]string)
	var entry []string

	for rows.Next() {
		var kanj string
		var von string
		var vkun string
		var transl string
		var roma string
		var rememb string
		var jlpt string
		var school string

		if err := rows.Scan(&kanj, &von, &vkun, &transl, &roma, &rememb, &jlpt, &school); err != nil {
			log.Fatal(err)
		}
		//		entry[kanj] = map[string]string{"kanj": kanj}
		//entry[kanj] = map[string]string{"kanj": kanj, "von": von, "vkun": vkun, "transl": transl, "roma": roma, "rememb": rememb, "jlpt": jlpt, "school": school}
		//		entry = append(entry, Kanji{kanj, von, vkun, transl, roma, rememb, jlpt, school})
		entry = append(entry, kanj)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)

	}

	selection := c.Param("selection")
	level := c.Param("level")
	entrymap := map[string]interface{}{"entry": entry, "selection": selection, "level": level}

	return c.Render(http.StatusOK, "kanji_list.html", entrymap) //map[string]interface{}{
	//	"entry":     entry,
	//	"selection": selection,
	//	"level":     level,
	//})

}

// GET /:selection/:level/:kanji
func getKanji(c echo.Context) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	// ensure :kanji isn't used as an escaped query like "%e9%9b%a8"
	uni_kanj, err := url.QueryUnescape(c.Param("kanji"))

	//	fmt.Println(c.Param("selection"), c.Param("level"))
	// start list of all in level get

	var sqlQuery string

	switch c.Param("selection") {
	case "grade":
		sqlQuery = "SELECT kanj, von, vkun, transl, roma, rememb, jlpt, school FROM info WHERE school = $1"
	case "jlpt":
		sqlQuery = "SELECT kanj, von, vkun, transl, roma, rememb, jlpt, school FROM info WHERE jlpt = $1"
	}
	rows, err := db.Query(sqlQuery, c.Param("level"))
	//	rows, err := db.Query(sqlQuery, c.Param("selection"), c.Param("level"))

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	other_kanj := make(map[string]int)
	kanj_index := make(map[int]string)

	k_index := 0

	for rows.Next() {
		var kanj string
		var von string
		var vkun string
		var transl string
		var roma string
		var rememb string
		var jlpt string
		var school string

		switch err := rows.Scan(&kanj, &von, &vkun, &transl, &roma, &rememb, &jlpt, &school); err {
		case sql.ErrNoRows:
			return c.Render(http.StatusNotFound, "404.html", "No rows were found")
		case nil:
			//fmt.Println(kanj, von, vkun, transl, roma, rememb, jlpt, school)
		default:
			log.Fatal(err)
		}

		other_kanj[kanj] = k_index
		kanj_index[k_index] = kanj
		//otherkanj = append(otherkanj, kanj)
		k_index++
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// start single kanji definition get

	if err != nil {
		log.Fatal(err)
	}

	singleQuery := "SELECT kanj, von, vkun, transl, roma, rememb, jlpt, school FROM info WHERE kanj = $1"
	row := db.QueryRow(singleQuery, uni_kanj)

	if err != nil {
		log.Fatal(err)
	}

	var kanj string
	var von string
	var vkun string
	var transl string
	var roma string
	var rememb string
	var jlpt string
	var school string
	var p_index int
	var n_index int
	var p_kanj string
	var n_kanj string
	var u_level string
	var u_selection string

	switch err := row.Scan(&kanj, &von, &vkun, &transl, &roma, &rememb, &jlpt, &school); err {
	case sql.ErrNoRows:
		// use a 404 here
		fmt.Println("No rows were returned!")
	case nil:
		//		fmt.Println(kanj, von, vkun, transl, roma, rememb, jlpt, school)
	default:
		log.Fatal(err)
	}

	num_items := len(other_kanj)

	p_index = other_kanj[uni_kanj] - 1
	n_index = other_kanj[uni_kanj] + 1

	// if we're at the beginning of the map, previous should be the last item
	if p_index < 0 {
		p_kanj = kanj_index[num_items-1]
	} else {
		p_kanj = kanj_index[p_index]
	}

	// if we reach the end of the map, next should cycle back to the beginning
	if n_index == num_items {
		n_kanj = kanj_index[0]
	} else {
		n_kanj = kanj_index[n_index]
	}

	u_level = c.Param("level")
	u_selection = c.Param("selection")

	entry := map[string]string{
		"kanj":        kanj,
		"von":         von,
		"vkun":        vkun,
		"transl":      transl,
		"roma":        roma,
		"rememb":      rememb,
		"jlpt":        jlpt,
		"school":      school,
		"p_kanj":      p_kanj,
		"n_kanj":      n_kanj,
		"u_level":     u_level,
		"u_selection": u_selection,
	}

	// TODO regex checking on values of :level and :selection
	return c.Render(http.StatusOK, "flashcard.html", entry)
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
	return c.String(http.StatusOK, response)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseFiles("tmpl/map.html",
			"tmpl/kanji_list.html",
			"tmpl/flashcard.html",
			"tmpl/container.html",
			"tmpl/header.html",
			"tmpl/404.html",
			"tmpl/episode_view.html",
			"tmpl/level_selection.html",
			"tmpl/main.html",
			"tmpl/footer.html",
		)),
	}
	e := echo.New()
	e.Static("/", "static")
	e.Renderer = t
	e.HTTPErrorHandler = custom404Handler
	//	e.Pre(middleware.HTTPSWWWRedirect())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", getMain)
	e.GET("/watch/:show/:season/:episode", getShow)
	//	e.GET("/grade/:level", getLevel)
	e.GET("/kanji", getJapanese)
	e.GET("/kanji/", getJapanese)
	e.GET("/kanjitainer", getContainer)
	e.GET("/kanjitainer/", getContainer)
	e.GET("/kanji/:selection/:level", getLevel)
	e.GET("/kanji/:selection/:level/:kanji", getKanji)
	e.GET("/.well-known/acme-challenge/test", getCert)
	e.GET("/.well-known/acme-challenge/test/", getCert)
	e.GET("/.well-known/acme-challenge/:response", getCert)
	e.GET("/.well-known/acme-challenge/:response/", getCert)
	e.GET("/well-known/acme-challenge/:response", getCert)
	e.GET("/well-known/acme-challenge/:response/", getCert)
	e.Logger.Info(e.Start(":8080"))
	//	e.Logger.Info(e.StartAutoTLS(":443"))
}
