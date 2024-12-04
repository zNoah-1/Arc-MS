package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/zNoah-1/Arc-MS/internal/logger"
	"github.com/zNoah-1/Arc-MS/internal/router"
	"github.com/zNoah-1/Arc-MS/internal/task"
)

var serverList []map[string]any

func main() {
	logger.Info("Starting up...")

	enterTestData()
	defineEndpoints()
	setupTimedTasks()
	startServer()
}

func enterTestData() {
	serverList = []map[string]any{
		{
			"ip":         "1.1.1.1",
			"port":       777,
			"contact":    "%40.nua",
			"id":         -1,
			"lastUpdate": int64(1733095960),
			"disabled":   false,
		},
	}
}

func defineEndpoints() {
	router.DefineEndpoints(&serverList)
}

func setupTimedTasks() {
	task.InitInactiveCheck(&serverList)
}

func startServer() {
	port := flag.Int("port", 8080, "Port to publish the server")
	flag.Parse()

	logger.Info("Starting server on port ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", *port), nil))
}
