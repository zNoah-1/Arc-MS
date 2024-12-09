package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zNoah-1/Arc-MS/internal/logger"
	"github.com/zNoah-1/Arc-MS/internal/util/httputil"
)

var lastId = 0

func register(w http.ResponseWriter, r *http.Request, serverList *[]map[string]any) {
	if r.Method != http.MethodPost {
		printReceivedRequest(r, true)
		http.Error(w, "Not implemented", http.StatusNotFound)
		return
	}
	//fmt.Println("Received", html.EscapeString(r.Method), html.EscapeString(r.URL.Path))
	printReceivedRequest(r, false)

	bodyBytes, err := httputil.BodyBytes(r)

	if len(bodyBytes) == 0 {
		http.Error(w, "Not enough data", http.StatusBadRequest)
		logger.Info("Body: No Body Supplied")
		return
	}

	if err != nil {
		http.Error(w, "An error has ocurred", http.StatusInternalServerError)
		return
	}

	id := lastId
	ip := httputil.UserIpAddr(r)
	port := httputil.ResponseValue(bodyBytes, "port")
	contact := httputil.ResponseValue(bodyBytes, "contact")

	serverInfo := make(map[string]any)
	serverInfo["id"] = id
	serverInfo["ip"] = ip
	serverInfo["port"] = port
	serverInfo["contact"] = contact
	serverInfo["disabled"] = false
	serverInfo["lastUpdate"] = time.Now().Unix()

	*serverList = append(*serverList, serverInfo)
	lastId++
	fmt.Fprintf(w, "%d", id)
}
