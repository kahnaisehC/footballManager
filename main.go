package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Template struct {
	templates *template.Template
}

// NOTE: PUT FIRST LETTER AS MAYUS
type Miembro struct {
	Nombre           *string
	Apellido         *string
	Telefono         *string
	DNI              int
	Fecha_nacimiento *time.Time
	Numero_socio     int
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func CreateMember(db *sql.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		nombre := c.FormValue("nombre")

		return c.Render(http.StatusOK, "successfullUpload", nombre)
	}
}

func Index(db *sql.DB) func(c echo.Context) error {
	var flaquetes []Miembro
	rows, err := db.Query("SELECT nombre, apellido, telefono, dni, fecha_nacimiento, numero_socio FROM jugadores")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var flaquete Miembro
		err := rows.Scan(
			&flaquete.Nombre,
			&flaquete.Apellido,
			&flaquete.Telefono,
			&flaquete.DNI,
			&flaquete.Fecha_nacimiento,
			&flaquete.Numero_socio,
		)
		if err != nil {
			panic(err)
		}
		flaquetes = append(flaquetes, flaquete)
	}
	fmt.Print(&flaquetes[0])

	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", flaquetes)
	}
}

func mustOpenDB(db *sql.DB, err error) *sql.DB {
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	t := &Template{
		templates: template.Must(
			template.New("").Funcs(template.FuncMap{
				"Deref": func(i *string) string { return *i },
			}).ParseGlob("public/views/*.html")),
	}
	var (
		DB_USER     = os.Getenv("DB_USER")
		DB_PORT     = os.Getenv("DB_PORT")
		DB_HOST     = os.Getenv("DB_HOST")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_NAME     = os.Getenv("DB_NAME")
	)
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
	db := mustOpenDB(sql.Open("postgres", psqlInfo))
	defer db.Close()

	e := echo.New()
	e.Renderer = t

	e.GET("/jugadores", Index(db))
	e.POST("/addPlayer", CreateMember(db))
	e.Logger.Fatal(e.Start(":1323"))
}
