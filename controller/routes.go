package controller

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type HtmlTemplate struct {
	templates *template.Template
}

func (t *HtmlTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func InitializeRoutes(e *echo.Echo) {
	//Initializing templates
	t := new(HtmlTemplate)
	t.templates = template.Must(template.ParseGlob("view/*.html"))
	e.Renderer = t

	//Initializing static content

	//Initializing routes
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/", HomeHandler)
	e.GET("/report", ReportHandler)
}
