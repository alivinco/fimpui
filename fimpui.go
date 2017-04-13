package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/alivinco/fimpui/integr/mqtt"
	"fmt"
)


func main() {
	wsUpgrader := mqtt.WsUpgrader{"localhost:1883"}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	index := "static/fimpui/dist/index.html"
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://ocalhost:4200"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.GET("/mqtt",wsUpgrader.Upgrade)
	e.File("/fimp", index)
	e.File("/fimp/zwave-man", index)
	e.File("/fimp/settings", index)
	e.File("/fimp/thing-view/*", index)
	e.Static("/fimp/static", "static/fimpui/dist/")
	e.Logger.Debug(e.Start(":8081"))
	fmt.Println("Exiting the app")
}
