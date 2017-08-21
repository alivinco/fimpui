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
	flowmodel "github.com/alivinco/fimpui/flow/model"
	log "github.com/Sirupsen/logrus"
	//"time"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"strconv"
)

type SystemInfo struct {
	Version string
}

// SetupLog configures default logger
// Supported levels : info , degug , warn , error
func SetupLog(logfile string,level string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})
	logLevel , err := log.ParseLevel(level)
	if err == nil {
		log.SetLevel(logLevel)
	}else {
		log.SetLevel(log.DebugLevel)
	}

	if logfile != "" {
		l := lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    5, // megabytes
			MaxBackups: 2,
		}
		log.SetOutput(&l)
	}

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

	SetupLog(configs.LogFile,configs.LogLevel)
	log.Info("--------------Starting FIMPUI----------------")
	//---------FLOW------------------------
	log.Info("<main> Starting Flow manager")
	flowManager,err := flow.NewManager(configs)
	if err != nil {
		log.Error("Can't Init Flow manager . Error :",err)
	}
	flowManager.InitMessagingTransport()
	err = flowManager.LoadAllFlowsFromStorage()
	if err != nil {
		log.Error("Can't load Flows from storage . Error :",err)
	}
	log.Info("<main> Started")
	//-------------------------------------
	//---------THINGS REGISTRY-------------
	log.Info("<main> Starting Things registry ")
	thingRegistryStore := registry.NewThingRegistryStore(configs.RegistryDbFile)
	log.Info("<main> Started ")
	//-------------------------------------
	//---------REGISTRY INTEGRATION--------
	log.Info("<main> Starting MqttIntegration ")
	mqttRegInt := registry.NewMqttIntegration(configs,thingRegistryStore)
	mqttRegInt.InitMessagingTransport()
	log.Info("<main> Started ")
	//-------------------------------------
	log.Info("<main> Starting VinculumClient ")
	vinculumClient := fhcore.NewVinculumClient(configs.VinculumAddress)
	err = vinculumClient.Connect()
	if err != nil {
		log.Error("<main> Can't connect to Vinculum")
	}else {
		log.Info("<main> Started ")
	}

	//---------GOOGLE OBJECT STORE---------
	log.Info("<main> Initializing Google Object Store ")
	objectStorage, _ := logexport.NewGcpObjectStorage("fh-cube-log")
	log.Info("<main> Done ")
	//-------------------------------------
	sysInfo := SystemInfo{}
	versionFile, err := ioutil.ReadFile("VERSION")
	if err == nil {
		sysInfo.Version = string(versionFile)
	}
	//--------VINCULUM PROXY----------------
	coreUrl := "ws://" + configs.VinculumAddress
	go startWsCoreProxy(coreUrl)
	//--------------------------------------

	wsUpgrader := mqtt.WsUpgrader{strings.Replace(configs.MqttServerURI,"tcp://","",-1)}
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
		log.Info(hostAlias)
		if hostAlias == "" {
			hostAlias = "unknown"
		}
		uploadStatus := objectStorage.UploadLogSnapshot(configs.ReportLogFiles, hostAlias, configs.ReportLogSizeLimit)
		return c.JSON(http.StatusOK, uploadStatus)
	})
	e.GET("/fimp/api/registry/things", func(c echo.Context) error {
		things , err := thingRegistryStore.GetAllThings()
		if err == nil {
			return c.JSON(http.StatusOK, things)
		}else {
			return c.JSON(http.StatusInternalServerError, err)
		}

	})

	e.GET("/fimp/api/registry/services", func(c echo.Context) error {
		services,err := thingRegistryStore.GetAllServices()
		if err == nil {
			return c.JSON(http.StatusOK, services)
		}else {
			return c.JSON(http.StatusInternalServerError, err)
		}

	})

	e.GET("/fimp/api/registry/interfaces", func(c echo.Context) error {
		thingAddr := c.QueryParam("thingAddr")
		thingTech := c.QueryParam("thingTech")
		serviceName := c.QueryParam("serviceName")
		intfMsgType := c.QueryParam("intfMsgType")
		locationIdStr := c.QueryParam("locationId")
		var locationId int
		locationId,_ = strconv.Atoi(locationIdStr)

		services,err := thingRegistryStore.GetFlatInterfaces(thingAddr,thingTech,serviceName,intfMsgType,registry.ID(locationId))
		if err == nil {
			return c.JSON(http.StatusOK, services)
		}else {
			return c.JSON(http.StatusInternalServerError, err)
		}

	})

	e.GET("/fimp/api/registry/locations", func(c echo.Context) error {
		locations,err := thingRegistryStore.GetAllLocations()
		if err == nil {
			return c.JSON(http.StatusOK, locations)
		}else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.GET("/fimp/api/registry/thing/:tech/:address", func(c echo.Context) error {
		things, _ := thingRegistryStore.GetThingByAddress(c.Param("tech"), c.Param("address"))
		return c.JSON(http.StatusOK, things)
	})
	e.GET("/fimp/api/registry/clear_all", func(c echo.Context) error {
		thingRegistryStore.ClearAll()
		return c.NoContent(http.StatusOK)
	})

	e.PUT("/fimp/api/registry/thing", func(c echo.Context) error {
		thing := registry.Thing{}
		err := c.Bind(&thing)
		fmt.Println(err)
		thingRegistryStore.UpsertThing(&thing)
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
		var resp *flowmodel.FlowMeta
		if id == "-" {
			flow := flowManager.GenerateNewFlow()
			resp = &flow;
		}else {
			resp = flowManager.GetFlowById(id).FlowMeta
		}

		return c.JSON(http.StatusOK, resp)
	})

	e.GET("/fimp/flow/context/:flowid", func(c echo.Context) error {
		id := c.Param("flowid")
		ctx := flowManager.GetGlobalContext().GetRecords(id)
		return c.JSON(http.StatusOK, ctx)
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
	log.Info("Exiting the app")
}
