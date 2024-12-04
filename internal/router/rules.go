package router

import (
	"fmt"
	"net/http"
)

func rules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		printReceivedRequest(r, true)
		http.Error(w, "Not implemented", http.StatusNotFound)
		return
	}
	//fmt.Println("Received", html.EscapeString(r.Method), html.EscapeString(r.URL.Path))
	printReceivedRequest(r, false)
	fmt.Fprint(w, "Rules: Soon (tm)\n\n")
}
