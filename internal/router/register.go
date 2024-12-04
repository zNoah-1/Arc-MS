package router

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

	if r.Body != nil {
		var bodyBytes []byte
		var err error
		bodyBytes, err = io.ReadAll(r.Body)

		if len(bodyBytes) > 0 {
			//fmt.Println(bodyBytes)
			bodyData := strings.Split(string(bodyBytes), "&")
			//fmt.Println(bodyData)
			serverInfo := make(map[string]any)

			i := 0
			for i < len(bodyData) {
				dataEntry := strings.Split(bodyData[i], "=")
				serverInfo[dataEntry[0]] = dataEntry[1]
				i++
			}

			serverInfo["ip"] = httputil.UserIpAddr(r)

			var id int = lastId
			serverInfo["id"] = id
			lastId++

			serverInfo["disabled"] = false
			serverInfo["lastUpdate"] = time.Now().Unix()

			logger.Debug(serverInfo)

			*serverList = append(*serverList, serverInfo)

			//fmt.Println(fmt.Sprintf("%x", bodyBytes))
			if err != nil {
				fmt.Printf("error: %v", err)
				http.Error(w, "Unknown error", http.StatusInternalServerError)
				return
			} else {
				fmt.Fprintf(w, "%d", id)
			}
		} else {
			http.Error(w, "Not enough data", http.StatusBadRequest)
			fmt.Printf("Body: No Body Supplied\n")
		}
	}
}
