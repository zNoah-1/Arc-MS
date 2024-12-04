package router

import (
	"fmt"
	"html"
	"net/http"
)

func serverListing(w http.ResponseWriter, r *http.Request, serverList *[]map[string]any) {
	if r.Method != http.MethodGet {
		printReceivedRequest(r, true)
		http.Error(w, "Not implemented", http.StatusNotFound)
		return
	}
	//fmt.Println("Received", html.EscapeString(r.Method), html.EscapeString(r.URL.Path))
	printReceivedRequest(r, false)

	var response string
	i := 0
	for i < len(*serverList) {
		serverData := (*serverList)[i]

		if serverData["disabled"] == false {
			response = response + fmt.Sprint(serverData["ip"], " ", serverData["port"], " ", serverData["contact"], "\n")
		}

		i++
	}

	fmt.Fprint(w, html.EscapeString(response))
}
