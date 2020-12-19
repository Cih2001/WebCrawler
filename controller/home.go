package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

// HomeHandler handles the requests for home page
func HomeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
