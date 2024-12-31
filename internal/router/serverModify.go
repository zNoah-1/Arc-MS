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

	pathList := pathList(r.URL.Path)
	id, err := strconv.Atoi(pathList[0])
	ipAddr := httputil.UserIpAddr(r)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	//Get server
	server, lastIndex := getServer(id, serverList)

	if server == nil {
		http.Error(w, "Can't find server with this ID", http.StatusNotFound)
		return
	}

	if server["ip"] != ipAddr {
		http.Error(w, "I'm sorry Dave, I'm afraid I can't do that", http.StatusForbidden)
		return
	}

	switch requestType(pathList) {
	case "unlist":
		*serverList = append((*serverList)[:lastIndex], (*serverList)[lastIndex+1:]...)
	case "update":
		server["lastUpdate"] = time.Now().Unix()
		server["disabled"] = false
	default:
		http.Error(w, "Not implemented", http.StatusNotFound)
	}
}

func pathList(uri string) []string {
	path := strings.TrimPrefix(uri, "/ms/api/servers/")
	return strings.Split(path, "/")
}

func requestType(pathList []string) string {
	if len(pathList) != 2 {
		return "unknown"
	}

	switch subpath := pathList[1]; subpath {
	case "unlist":
		return "unlist"
	case "update":
		return "update"
	default:
		return "unknown"
	}
}

func getServer(id int, serverList *[]map[string]any) (map[string]any, int) {
	i := 0
	for i < len(*serverList) {
		server := (*serverList)[i]

		if server["id"] == id {
			return server, i
		}
		i++
	}
	return nil, -1
}
