package main

import (
	"Rvvijays/logapi/api"
	"Rvvijays/logapi/automated"
	"Rvvijays/logapi/util"
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/valyala/fasthttp"
)

func init() {
	// initialize configs
	util.InitConfig()
	util.LogsMap = cmap.New()

	// goroutine to write logs from memory
	if util.Config.WriteLogInterval > 0 {
		go automated.LogMapToLocalFile()
	}

	// goroutine to uploads file to storage server
	// check config.json, you can make it off by setting uploadMaster 0
	if util.Config.UploadLogInterval > 0 && util.Config.ServerStorage {
		go automated.LocalFileToServer()
	}

}

func main() {

	router := fasthttprouter.New()
	router.POST("/api/logs/search", api.SearchLog)

	router.POST("/api/logs/add", api.AddLog)

	router.GET("/api/logs/current", api.GetCurrentLogs)

	fmt.Println("running on ", util.Config.Port)

	log.Fatal(fasthttp.ListenAndServe(":"+util.Config.Port, router.Handler))
}
