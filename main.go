package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var lastId int = 0
var serverList []map[string]any

func main() {
	printWithTimestamp("Starting up...")

	//enterTestData()
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
			"id":         0,
			"lastUpdate": int64(1732238430),
			"disabled":   false,
		},
	}
}

func defineEndpoints() {
	//GET
	http.HandleFunc("/ms/api/rules", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			printReceivedRequest(r, true)
			http.Error(w, "Not implemented", http.StatusNotFound)
			return
		}
		//fmt.Println("Received", html.EscapeString(r.Method), html.EscapeString(r.URL.Path))
		printReceivedRequest(r, false)

		fmt.Fprint(w, "Rules: Soon (tm)\n\n")
	})

	//GET
	http.HandleFunc("/ms/api/games/RingRacers/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			printReceivedRequest(r, true)
			http.Error(w, "Not implemented", http.StatusNotFound)
			return
		}
		//fmt.Println("Received", html.EscapeString(r.Method), html.EscapeString(r.URL.Path))
		printReceivedRequest(r, false)

		fmt.Fprint(w, "4 v2.3\n")
	})

	//GET
	http.HandleFunc("/ms/api/games/RingRacers/4/servers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			printReceivedRequest(r, true)
			http.Error(w, "Not implemented", http.StatusNotFound)
			return
		}
		//fmt.Println("Received", html.EscapeString(r.Method), html.EscapeString(r.URL.Path))
		printReceivedRequest(r, false)

		var response string
		i := 0
		for i < len(serverList) {
			serverData := serverList[i]

			if serverData["disabled"] == false {
				response = response + fmt.Sprint(serverData["ip"], " ", serverData["port"], " ", serverData["contact"], "\n")
			}

			i++
		}

		fmt.Fprint(w, html.EscapeString(response))
	})

	//POST
	http.HandleFunc("/ms/api/servers/", func(w http.ResponseWriter, r *http.Request) {
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

			ipAddr := getUserIpAddr(r)

			i := 0
			for i < len(serverList) {
				serverData := serverList[i]

				if serverData["id"] == id {
					if serverData["ip"] == ipAddr {
						if parts[1] == "unlist" {
							serverList = append(serverList[:i], serverList[i+1:]...)
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
	})

	//POST
	http.HandleFunc("/ms/api/games/RingRacers/4/servers/register", func(w http.ResponseWriter, r *http.Request) {
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

				serverInfo["ip"] = getUserIpAddr(r)

				var id int = lastId
				serverInfo["id"] = id
				lastId++

				serverInfo["disabled"] = false
				serverInfo["lastUpdate"] = time.Now().Unix()

				printWithTimestamp(serverInfo)

				serverList = append(serverList, serverInfo)

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
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("request received: ", r)
		printReceivedRequest(r, true)
		http.Error(w, "Not implemented", http.StatusNotFound)
	})
}

func setupTimedTasks() {
	go checkForInactiveServers()
}

func checkForInactiveServers() {
	for true {
		printWithTimestamp("Checking for inactive servers...")

		unixTime := time.Now().Unix()
		i := 0
		for i < len(serverList) {
			serverData := serverList[i]
			lastUpdate := serverData["lastUpdate"].(int64)

			if unixTime-lastUpdate > 3600 { //3600 = 1h
				printWithTimestamp(fmt.Sprint("Removing server ", serverData["ip"], " for inactivity"))
				serverList = append(serverList[:i], serverList[i+1:]...)
				i--
			} else if unixTime-lastUpdate > 1800 { //1800 = 30 min
				if serverData["disabled"].(bool) {
					break
				}

				printWithTimestamp(fmt.Sprint("Disabling server ", serverData["ip"], " for inactivity"))
				serverData["disabled"] = true
			}

			i++
		}

		time.Sleep(60 * time.Second)
	}
}

func getUserIpAddr(req *http.Request) string {
	//Consider "X-FORWARDED-FOR" for reverse proxy setup
	ipAddr, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ipAddr
}

func printReceivedRequest(req *http.Request, verbose bool) {
	if verbose {
		printWithTimestamp("Request received: ", req)

		return
	}

	printWithTimestamp("Received request: ", html.EscapeString(req.Method), " ", html.EscapeString(req.URL.Path))
}

func printWithTimestamp(output ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05") + ":"
	fmt.Println(timestamp, fmt.Sprint(output...))

}

func startServer() {
	port := flag.Int("port", 8080, "Port to publish the server")
	flag.Parse()

	printWithTimestamp("Starting server on port ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", *port), nil))
}
