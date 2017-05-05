package main

import (
	"net/http"

	"fmt"

	"github.com/alivinco/fimpui/integr/mqtt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
)


type SystemInfo struct {
	Version string
}

func main() {
	sysInfo := SystemInfo{}
	versionFile,err := ioutil.ReadFile("VERSION")
	if err == nil {
		sysInfo.Version = string(versionFile)
	}

	wsUpgrader := mqtt.WsUpgrader{"localhost:1883"}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/fimp/system-info", func(c echo.Context) error {

		return c.JSON(http.StatusOK,sysInfo)
	})
	index := "static/fimpui/dist/index.html"
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.GET("/mqtt", wsUpgrader.Upgrade)
	e.File("/fimp", index)
	e.File("/fimp/zwave-man", index)
	e.File("/fimp/settings", index)
	e.File("/fimp/timeline", index)
	e.File("/fimp/ikea-man", index)
	e.File("/fimp/thing-view/*", index)
	e.Static("/fimp/static", "static/fimpui/dist/")
	e.Logger.Debug(e.Start(":8081"))
	fmt.Println("Exiting the app")
}
