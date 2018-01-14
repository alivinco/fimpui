package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpui/flow"
	flowmodel "github.com/alivinco/fimpui/flow/model"
	"github.com/alivinco/fimpui/integr/fhcore"
	"github.com/alivinco/fimpui/integr/logexport"
	"github.com/alivinco/fimpui/integr/mqtt"
	"github.com/alivinco/fimpui/model"
	"github.com/alivinco/fimpui/process"
	"github.com/alivinco/fimpui/registry"
	"github.com/koding/websocketproxy"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"time"
	"strconv"
	"strings"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"github.com/alivinco/fimpui/integr/zwave"
	"github.com/alivinco/fimpui/statsdb"
	//_ "net/http/pprof"
)

type SystemInfo struct {
	Version string
}

// SetupLog configures default logger
// Supported levels : info , degug , warn , error
func SetupLog(logfile string, level string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})
	logLevel, err := log.ParseLevel(level)
	if err == nil {
		log.SetLevel(logLevel)
	} else {
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
	// pprof server
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
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

	SetupLog(configs.LogFile, configs.LogLevel)
	log.Info("--------------Starting FIMPUI----------------")
	//---------FLOW------------------------
	log.Info("<main> Starting Flow manager")
	flowManager, err := flow.NewManager(configs)
	if err != nil {
		log.Error("Can't Init Flow manager . Error :", err)
	}
	flowManager.InitMessagingTransport()
	err = flowManager.LoadAllFlowsFromStorage()
	if err != nil {
		log.Error("Can't load Flows from storage . Error :", err)
	}
	log.Info("<main> Started")
	//-------------------------------------
	//---------THINGS REGISTRY-------------
	log.Info("<main>-------------- Starting Things registry ")
	thingRegistryStore := registry.NewThingRegistryStore(configs.RegistryDbFile)
	log.Info("<main> Started ")
	//-------------------------------------
	//---------REGISTRY INTEGRATION--------
	log.Info("<main>-------------- Starting MqttIntegration ")
	mqttRegInt := registry.NewMqttIntegration(configs, thingRegistryStore)
	mqttRegInt.InitMessagingTransport()
	log.Info("<main> Started ")
	//---------STATS STORE-----------------
	log.Info("<main>-------------- Stats store ")
	statsStore := statsdb.NewStatsStore("stats.db")
	streamProcessor := statsdb.NewStreamProcessor(configs,statsStore)
	streamProcessor.InitMessagingTransport()
	log.Info("<main> Started ")
	//----------VINCULUM CLIENT------------
	log.Info("<main>-------------- Starting VinculumClient ")
	vinculumClient := fhcore.NewVinculumClient(configs.VinculumAddress)
	err = vinculumClient.Connect()
	if err != nil {
		log.Error("<main> Can't connect to Vinculum")
	} else {
		log.Info("<main> Started ")
	}
    // --------VINCULUM ADAPTER------------
	log.Info("<main>-------------- Starting VinculumAdapter ")
	vinculumAd := fhcore.NewVinculumAdapter(configs,vinculumClient)
	vinculumAd.InitMessagingTransport()

	//---------GOOGLE OBJECT STORE---------
	log.Info("<main>-------------- Initializing Google Object Store ")
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
	var brokerAddress string
	var isSSL bool
	if strings.Contains(configs.MqttServerURI,"ssl") {
		brokerAddress = strings.Replace(configs.MqttServerURI, "ssl://", "", -1)
		isSSL = true
	}else {
		brokerAddress = strings.Replace(configs.MqttServerURI, "tcp://", "", -1)
		isSSL = false
	}
	wsUpgrader := mqtt.WsUpgrader{brokerAddress,isSSL}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/fimp/system-info", func(c echo.Context) error {

		return c.JSON(http.StatusOK, sysInfo)
	})
	e.GET("/fimp/api/configs", func(c echo.Context) error {
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
	e.POST("/fimp/api/zwave/products/upload-to-cloud", func(c echo.Context) error {
		cloud,err  := zwave.NewProductCloudStore( configs.ZwaveProductTemplates,"fh-products")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		err = cloud.UploadProductCacheToCloud()
		if err == nil {
			return c.NoContent(http.StatusOK)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.GET("/fimp/api/zwave/products/list-local-templates", func(c echo.Context) error {
		templateType := c.QueryParam("type")
		returnStable := true
		if templateType == "cache" {
			returnStable = false
		}

		cloud,err  := zwave.NewProductCloudStore( configs.ZwaveProductTemplates,"fh-products")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		templates,err := cloud.ListTemplates(returnStable)
		if err == nil {
			return c.JSON(http.StatusOK,templates)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.GET("/fimp/api/zwave/products/template", func(c echo.Context) error {
		templateType := c.QueryParam("type")
		fileName := c.QueryParam("name")
		returnStable := true
		if templateType == "cache" {
			returnStable = false
		}

		store,err  := zwave.NewProductCloudStore( configs.ZwaveProductTemplates,"fh-products")
		template , err :=store.GetTemplate(returnStable,fileName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		if err == nil {
			return c.Blob(http.StatusOK,"application/json",template)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.POST("/fimp/api/zwave/products/template-op/:operation/:name", func(c echo.Context) error {
		operation := c.Param("operation")
		name := c.Param("name")
		store,_  := zwave.NewProductCloudStore( configs.ZwaveProductTemplates,"fh-products")
		var err error
		switch operation {
		case "move":
			err = store.MoveToStable(name)
		case "upload":
			err = store.UploadSingleProductToStageCloud(name)
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError,err)
		}
		return c.NoContent(http.StatusOK)
	})

	e.DELETE("/fimp/api/zwave/products/template/:type/:name", func(c echo.Context) error {
		templateType := c.Param("type")
		templateName := c.Param("name")
		var isStable bool
		switch templateType {
		case "cache":
			isStable = false
		case "stable":
			isStable = true
		default:
			return c.NoContent(http.StatusInternalServerError)

		}
		store,err  := zwave.NewProductCloudStore( configs.ZwaveProductTemplates,"fh-products")
		err = store.DeleteTemplate(isStable,templateName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,err)
		}
		return c.NoContent(http.StatusOK)
	})

	e.POST("/fimp/api/zwave/products/template/:type/:name", func(c echo.Context) error {
		templateType := c.Param("type")
		templateName := c.Param("name")
		var isStable bool
		switch templateType {
		case "cache":
			isStable = false
		case "stable":
			isStable = true
		default:
			return c.NoContent(http.StatusInternalServerError)

		}
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		store,err  := zwave.NewProductCloudStore( configs.ZwaveProductTemplates,"fh-products")
		err = store.UpdateTemplate(isStable,templateName,body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,err)
		}
		return c.NoContent(http.StatusOK)
	})

	e.POST("/fimp/api/zwave/products/download-from-cloud", func(c echo.Context) error {
		cloud,err  := zwave.NewProductCloudStore( configs.ZwaveProductTemplates,"fh-products")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		prodNames,err := cloud.DownloadProductsFromCloud()
		if err == nil {
			return c.JSON(http.StatusOK,prodNames)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.GET("/fimp/api/stats/event-log", func(c echo.Context) error {

		pageSize := 1000
		page := 0
		pageSize, _ = strconv.Atoi(c.QueryParam("pageSize"))
		page, _ = strconv.Atoi(c.QueryParam("page"))
		statsErrors, err := statsStore.GetEvents(pageSize,page)

		if err == nil {
			return c.JSON(http.StatusOK, statsErrors)
		} else {
			log.Error("Faild to fetch errors ",err)
			return c.JSON(http.StatusInternalServerError, err)
		}

	})

	e.GET("/fimp/api/stats/metrics/counters", func(c echo.Context) error {

		result := make(map[string]interface{})
		result["restart_time"] = statsStore.GetResetTime()
		result["metrics"] = statsStore.GetCounterMetrics()

		if err == nil {
			return c.JSON(http.StatusOK, result)
		} else {
			log.Error("Faild to fetch errors ",err)
			return c.JSON(http.StatusInternalServerError, err)
		}

	})

	e.GET("/fimp/api/stats/metrics/meters", func(c echo.Context) error {

		result := make(map[string]interface{})
		result["restart_time"] = statsStore.GetResetTime()
		result["metrics"] = statsStore.GetMeterMetrics()

		if err == nil {
			return c.JSON(http.StatusOK, result)
		} else {
			log.Error("Faild to fetch errors ",err)
			return c.JSON(http.StatusInternalServerError, err)
		}

	})

	e.GET("/fimp/api/registry/things", func(c echo.Context) error {

		var things []registry.Thing
		var locationId int
		locationIdStr := c.QueryParam("locationId")
		locationId, _ = strconv.Atoi(locationIdStr)

		if locationId != 0 {
			things, err = thingRegistryStore.GetThingsByLocationId(registry.ID(locationId))
		} else {
			things, err = thingRegistryStore.GetAllThings()
		}
		thingsWithLocation := thingRegistryStore.ExtendThingsWithLocation(things)
		if err == nil {
			return c.JSON(http.StatusOK, thingsWithLocation)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}

	})

	e.GET("/fimp/api/registry/services", func(c echo.Context) error {
		serviceName := c.QueryParam("serviceName")
		locationIdStr := c.QueryParam("locationId")
		thingIdStr := c.QueryParam("thingId")
		thingId, _ := strconv.Atoi(thingIdStr)
		locationId , _ := strconv.Atoi(locationIdStr)
		filterWithoutAliasStr:= c.QueryParam("filterWithoutAlias")
		var filterWithoutAlias bool
		if filterWithoutAliasStr == "true" {
			filterWithoutAlias = true
		}
		services, err := thingRegistryStore.GetExtendedServices(serviceName,filterWithoutAlias,registry.ID(thingId),registry.ID(locationId))
		if err == nil {
			return c.JSON(http.StatusOK, services)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	//e.POST("/fimp/api/registry/service-fields", func(c echo.Context) error {
	//	// The service update only selected fields and not entire object
	//	service := registry.Service{}
	//	err := c.Bind(&service)
	//	if err == nil {
	//		log.Info("<REST> Saving service fields")
	//		thingRegistryStore.UpsertService(&service)
	//		return c.NoContent(http.StatusOK)
	//	} else {
	//		log.Info("<REST> Can't bind service")
	//		return c.JSON(http.StatusInternalServerError, err)
	//	}
	//})
	//
	//e.POST("/fimp/api/registry/thing-fields", func(c echo.Context) error {
	//	// The service update only selected fields and not entire object
	//	thing := registry.Thing{}
	//	err := c.Bind(&thing)
	//	if err == nil {
	//		log.Info("<REST> Saving thing fields")
	//		thingRegistryStore.UpsertThing(&thing)
	//		return c.NoContent(http.StatusOK)
	//	} else {
	//		log.Info("<REST> Can't bind thing")
	//		return c.JSON(http.StatusInternalServerError, err)
	//	}
	//})

	e.PUT("/fimp/api/registry/service", func(c echo.Context) error {
		service := registry.Service{}
		err := c.Bind(&service)
		if err == nil {
			log.Info("<REST> Saving service")
			thingRegistryStore.UpsertService(&service)
			return c.NoContent(http.StatusOK)
		} else {
			log.Info("<REST> Can't bind service")
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.PUT("/fimp/api/registry/location", func(c echo.Context) error {
		location := registry.Location{}
		err := c.Bind(&location)
		if err == nil {
			log.Info("<REST> Saving location")
			thingRegistryStore.UpsertLocation(&location)
			return c.NoContent(http.StatusOK)
		} else {
			log.Info("<REST> Can't bind location")
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.GET("/fimp/api/registry/interfaces", func(c echo.Context) error {
		//thingAddr := c.QueryParam("thingAddr")
		//thingTech := c.QueryParam("thingTech")
		//serviceName := c.QueryParam("serviceName")
		//intfMsgType := c.QueryParam("intfMsgType")
		//locationIdStr := c.QueryParam("locationId")
		//var locationId int
		//locationId, _ = strconv.Atoi(locationIdStr)
		//var thingId int
		//thingIdStr := c.QueryParam("thingId")
		//thingId, _ = strconv.Atoi(thingIdStr)
		//services, err := thingRegistryStore.GetFlatInterfaces(thingAddr, thingTech, serviceName, intfMsgType, registry.ID(locationId), registry.ID(thingId))
		services := []registry.ServiceExtendedView{}
		if err == nil {
			return c.JSON(http.StatusOK, services)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.GET("/fimp/api/registry/locations", func(c echo.Context) error {
		locations, err := thingRegistryStore.GetAllLocations()
		if err == nil {
			return c.JSON(http.StatusOK, locations)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	})

	e.GET("/fimp/api/registry/thing/:tech/:address", func(c echo.Context) error {
		things, err := thingRegistryStore.GetThingExtendedViewByAddress(c.Param("tech"), c.Param("address"))
		if err == nil {
			return c.JSON(http.StatusOK, things)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}

	})
	e.DELETE("/fimp/api/registry/clear_all", func(c echo.Context) error {
		thingRegistryStore.ClearAll()
		return c.NoContent(http.StatusOK)
	})

	e.POST("/fimp/api/registry/reindex", func(c echo.Context) error {
		thingRegistryStore.ReindexAll()
		return c.NoContent(http.StatusOK)
	})

	e.PUT("/fimp/api/registry/thing", func(c echo.Context) error {
		thing := registry.Thing{}
		err := c.Bind(&thing)
		fmt.Println(err)
		if err == nil {
			log.Info("<REST> Saving thing")
			thingRegistryStore.UpsertThing(&thing)
			return c.NoContent(http.StatusOK)
		} else {
			log.Info("<REST> Can't bind thing")
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	})

	e.DELETE("/fimp/api/registry/thing/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		thingId, _ := strconv.Atoi(idStr)
		err := thingRegistryStore.DeleteThing(registry.ID(thingId))
		if err == nil {
			return c.NoContent(http.StatusOK)
		}
		log.Error("<REST> Can't delete thing ")
		return c.JSON(http.StatusInternalServerError, err)
	})

	e.DELETE("/fimp/api/registry/location/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		thingId, _ := strconv.Atoi(idStr)
		err := thingRegistryStore.DeleteLocation(registry.ID(thingId))
		if err == nil {
			return c.NoContent(http.StatusOK)
		}
		log.Error("<REST> Failed to delete thing . Error : ",err)
		return c.JSON(http.StatusInternalServerError, err)
	})

	e.GET("/fimp/vinculum/devices", func(c echo.Context) error {
		resp, _ := vinculumClient.GetMessage([]string{"device"})
		return c.JSON(http.StatusOK, resp.Msg.Data.Param.Device)
	})

	e.GET("/fimp/vinculum/shortcuts", func(c echo.Context) error {
		resp, _ := vinculumClient.GetShortcuts()
		return c.JSON(http.StatusOK, resp)
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
			resp = &flow
		} else {
			resp = flowManager.GetFlowById(id).FlowMeta
		}

		return c.JSON(http.StatusOK, resp)
	})

	e.GET("/fimp/flow/context/:flowid", func(c echo.Context) error {
		id := c.Param("flowid")
		if id != "-"{
			ctx := flowManager.GetGlobalContext().GetRecords(id)
			return c.JSON(http.StatusOK, ctx)
		}
		var ctx []flowmodel.ContextRecord
		return c.JSON(http.StatusOK, ctx)


	})

	e.POST("/fimp/flow/definition/:id", func(c echo.Context) error {
		id := c.Param("id")
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		flowManager.UpdateFlowFromJsonAndSaveToStorage(id, body)
		return c.NoContent(http.StatusOK)
	})

	e.POST("/fimp/flow/ctrl/:id/:op", func(c echo.Context) error {
		id := c.Param("id")
		op := c.Param("op")

		switch op {
		case "send-inclusion-report" :
			flowManager.SendInclusionReport(id)
		case "send-exclusion-report" :
			flowManager.SendExclusionReport(id)
		}

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
	e.File("/fimp/systems-man", index)
	e.File("/fimp/flow/context", index)
	e.File("/fimp/flow/overview", index)
	e.File("/fimp/flow/flow-editor/*", index)
	e.File("/fimp/flight-recorder", index)
	e.File("/fimp/thing-view/*", index)
	e.File("/fimp/registry/things/*", index)
	e.File("/fimp/registry/services/*", index)
	e.File("/fimp/registry/locations", index)
	e.File("/fimp/registry/admin", index)
	e.Static("/fimp/static", "static/fimpui/dist/")
	e.Logger.Debug(e.Start(":8081"))
	log.Info("Exiting the app")


}
