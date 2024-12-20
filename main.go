package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "ianfli")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	var (
		DB_USER     = os.Getenv("DB_USER")
		DB_PORT     = os.Getenv("DB_PORT")
		DB_HOST     = os.Getenv("DB_HOST")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_NAME     = os.Getenv("DB_NAME")
	)
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	rows, err := db.Query("SELECT nombre FROM jugadores;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var jug string
		if err := rows.Scan(&jug); err != nil {
			panic(err)
		}
		fmt.Print(jug)
	}
	e := echo.New()
	e.Renderer = t
	e.GET("/", Index)

	e.Logger.Fatal(e.Start(":1323"))
}
