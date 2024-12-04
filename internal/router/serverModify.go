package router

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zNoah-1/Arc-MS/internal/util/httputil"
)

func serverModify(w http.ResponseWriter, r *http.Request, serverList *[]map[string]any) {
	if r.Method != http.MethodPost {
		printReceivedRequest(r, true)
		http.Error(w, "Not implemented", http.StatusNotFound)
		return
	}
	//fmt.Println("Received", html.EscapeString(r.Method), html.EscapeString(r.URL.Path))
	printReceivedRequest(r, false)

	path := strings.TrimPrefix(r.URL.Path, "/ms/api/servers/")
	parts := strings.Split(path, "/")

	if len(parts) == 2 {
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Process with the ID
		//fmt.Fprintf(w, "Server ID: %d\n", id)

		ipAddr := httputil.UserIpAddr(r)

		i := 0
		for i < len(*serverList) {
			serverData := (*serverList)[i]

			if serverData["id"] == id {
				if serverData["ip"] == ipAddr {
					if parts[1] == "unlist" {
						*serverList = append((*serverList)[:i], (*serverList)[i+1:]...)
					} else if parts[1] == "update" {
						serverData["lastUpdate"] = time.Now().Unix()
						serverData["disabled"] = false
					} else {
						http.Error(w, "Not implemented", http.StatusNotFound)
					}
					return
				}

				http.Error(w, "You are not allowed to do that", http.StatusForbidden)
				return
			}
			i++
		}
		http.Error(w, "Server Not Found", http.StatusInternalServerError)
	} else {
		http.Error(w, "Not implemented", http.StatusNotFound)
	}
}
