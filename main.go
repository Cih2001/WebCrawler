package main

import (
	"Cih2001/WebCrawler/controller"
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <server_listen_addr>\n", os.Args[0])
		fmt.Printf("e.g.: %s :1323 \n", os.Args[0])
		return
	}
	// we use labstack echo framework.
	e := echo.New()

	controller.InitializeRoutes(e)

	// print logs on stdout
	e.Use(middleware.Logger())

	// we start server on port 1323. this should be put in a config file.
	// however, for our simple project, we just hardcode it and we don't use
	// a config file at all.
	if err := e.Start(os.Args[1]); err != nil {
		fmt.Printf("Error listening on %s with error:%s\n", os.Args[1], err.Error())
	}
}
