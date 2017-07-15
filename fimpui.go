package main

import (
	"net/http"

	"fmt"

	"encoding/json"
	"flag"
	"github.com/alivinco/fimpui/integr/fhcore"
	"github.com/alivinco/fimpui/integr/logexport"
	"github.com/alivinco/fimpui/integr/mqtt"
	"github.com/alivinco/fimpui/model"
	"github.com/alivinco/fimpui/process"
	"github.com/alivinco/fimpui/registry"
	"github.com/koding/websocketproxy"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"net/url"
	"github.com/alivinco/fimpui/flow"
	log "github.com/Sirupsen/logrus"
)

type SystemInfo struct {
	Version string
}

func startWsCoreProxy(backendUrl string) {
	u, _ := url.Parse(backendUrl)
	http.Handle("/", http.FileServer(http.Dir("static/fhcore")))
	http.Handle("/ws", websocketproxy.ProxyHandler(u))
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Print(err)
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	configs := &model.FimpUiConfigs{}
	var configFile string
	flag.StringVar(&configFile, "c", "", "Config file")
	flag.Parse()
	if configFile == "" {
		configFile = "/opt/fimpui/config.json"
	} else {
		fmt.Println("Loading configs from file ", configFile)
	}
	configFileBody, err := ioutil.ReadFile(configFile)
	err = json.Unmarshal(configFileBody, configs)
	if err != nil {
		panic("Can't load config file.")
	}
	//---------FLOW------------------------
	flowManager := flow.NewManager(configs)
	flowManager.InitMessagingTransport()
	err = flowManager.LoadAllFlowsFromStorage()
	if err != nil {
		log.Error("Can't load Flows from storage . Error :",err)
	}
	//-------------------------------------
	//---------THINGS REGISTRY-------------
	thingRegistryStore := registry.NewThingRegistryStore("thingsStore.json")
	//-------------------------------------
	vinculumClient := fhcore.NewVinculumClient(configs.VinculumAddress)
	err = vinculumClient.Connect()
	if err != nil {
		fmt.Println("Vinculum is not connected")
	}

	objectStorage, _ := logexport.NewGcpObjectStorage("fh-cube-log")
	sysInfo := SystemInfo{}
	versionFile, err := ioutil.ReadFile("VERSION")
	if err == nil {
		sysInfo.Version = string(versionFile)
	}
	//--------VINCULUM PROXY----------------
	coreUrl := "ws://" + configs.VinculumAddress
	go startWsCoreProxy(coreUrl)
	//--------------------------------------
	wsUpgrader := mqtt.WsUpgrader{"localhost:1883"}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/fimp/system-info", func(c echo.Context) error {

		return c.JSON(http.StatusOK, sysInfo)
	})
	e.GET("/fimp/configs", func(c echo.Context) error {
		return c.JSON(http.StatusOK, configs)
	})
	e.GET("/fimp/fr/upload-log-snapshot", func(c echo.Context) error {

		//logexport.UploadLogToGcp()
		//files := []string {"/var/log/daily.out"}
		hostAlias := c.QueryParam("hostAlias")
		fmt.Println(hostAlias)
		if hostAlias == "" {
			hostAlias = "unknown"
		}
		uploadStatus := objectStorage.UploadLogSnapshot(configs.ReportLogFiles, hostAlias, configs.ReportLogSizeLimit)
		return c.JSON(http.StatusOK, uploadStatus)
	})
	e.GET("/fimp/registry/things", func(c echo.Context) error {
		things := thingRegistryStore.GetAllThings()
		return c.JSON(http.StatusOK, things)
	})

	e.GET("/fimp/registry/thing/:tech/:address", func(c echo.Context) error {
		things, _ := thingRegistryStore.GetThingByAddress(c.Param("tech"), c.Param("address"))
		return c.JSON(http.StatusOK, things)
	})
	e.GET("/fimp/registry/clear_all", func(c echo.Context) error {
		thingRegistryStore.ClearAll()
		return c.NoContent(http.StatusOK)
	})

	e.PUT("/fimp/registry/thing", func(c echo.Context) error {
		thing := registry.Thing{}
		err := c.Bind(&thing)
		fmt.Println(err)
		thingRegistryStore.UpsertThing(thing)
		return c.NoContent(http.StatusOK)
	})

	e.GET("/fimp/vinculum/devices", func(c echo.Context) error {
		resp, _ := vinculumClient.GetMessage([]string{"device"})
		return c.JSON(http.StatusOK, resp.Msg.Data.Param.Device)
	})

	e.GET("/fimp/vinculum/import_to_registry", func(c echo.Context) error {
		process.LoadVinculumDeviceInfoToStore(thingRegistryStore, vinculumClient)
		return c.NoContent(http.StatusOK)
	})

	e.GET("/fimp/flow/list", func(c echo.Context) error {
		resp := flowManager.GetFlowList()
		return c.JSON(http.StatusOK, resp)
	})
	e.GET("/fimp/flow/definition/:id", func(c echo.Context) error {
		id := c.Param("id")
		var resp *flow.Flow
		if id == "-" {
			flow := flowManager.GenerateNewFlow()
			resp = &flow;
		}else {
			resp = flowManager.GetFlowById(id)
		}

		return c.JSON(http.StatusOK, resp)
	})

	e.POST("/fimp/flow/definition/:id", func(c echo.Context) error {
		id := c.Param("id")
		body ,err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		flowManager.UpdateFlowFromJsonAndSaveToStorage(id,body)
		return c.NoContent(http.StatusOK)
	})
	e.DELETE("/fimp/flow/definition/:id", func(c echo.Context) error {
		id := c.Param("id")
		flowManager.DeleteFlow(id)
		return c.NoContent(http.StatusOK)
	})

	index := "static/fimpui/dist/index.html"
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200", "http:://localhost:8082"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.GET("/mqtt", wsUpgrader.Upgrade)
	e.File("/fimp", index)
	//e.File("/fhcore", "static/fhcore.html")
	e.File("/fimp/zwave-man", index)
	e.File("/fimp/settings", index)
	e.File("/fimp/timeline", index)
	e.File("/fimp/ikea-man", index)
	e.File("/fimp/flow", index)
	e.File("/fimp/flow/flow-editor/*", index)
	e.File("/fimp/flight-recorder", index)
	e.File("/fimp/thing-view/*", index)
	e.Static("/fimp/static", "static/fimpui/dist/")
	e.Logger.Debug(e.Start(":8081"))
	fmt.Println("Exiting the app")
}
