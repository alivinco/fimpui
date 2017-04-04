package main

import (
"net/http"

"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Static("/", "static/fimpui/dist")
	e.Logger.Fatal(e.Start(":8081"))
}
