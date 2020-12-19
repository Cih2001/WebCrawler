package main

import (
	"Cih2001/WebCrawler/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// we use labstack echo framework.
	e := echo.New()

	controller.InitializeRoutes(e)

	// print logs on stdout
	e.Use(middleware.Logger())

	// we start server on port 1323. this should be put in a config file.
	// however, for our simple project, we just hardcode it and we don't use
	// a config file at all.
	e.Start(":1323")
}
