package router

import (
	"html"
	"net/http"

	"github.com/zNoah-1/Arc-MS/internal/logger"
)

func DefineEndpoints(serverList *[]map[string]any) {
	//GET
	http.HandleFunc("/ms/api/rules", rules)
	http.HandleFunc("/ms/api/games/RingRacers/version", version)
	http.HandleFunc("/ms/api/games/RingRacers/4/servers", func(w http.ResponseWriter, r *http.Request) { serverListing(w, r, serverList) })

	//POST
	http.HandleFunc("/ms/api/servers/", func(w http.ResponseWriter, r *http.Request) { serverModify(w, r, serverList) })
	http.HandleFunc("/ms/api/games/RingRacers/4/servers/register", func(w http.ResponseWriter, r *http.Request) { register(w, r, serverList) })

	//ALL
	http.HandleFunc("/", defaultRoute)
}

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("request received: ", r)
	printReceivedRequest(r, true)
	http.Error(w, "Not implemented", http.StatusNotFound)
}

func printReceivedRequest(req *http.Request, verbose bool) {
	if verbose {
		logger.Debug("Request received: ", req)
		return
	}

	logger.Info("Received request: ", html.EscapeString(req.Method), " ", html.EscapeString(req.URL.Path))
}
