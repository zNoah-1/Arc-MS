package httputil

import (
	"net"
	"net/http"
)

func UserIpAddr(req *http.Request) string {
	//Consider "X-FORWARDED-FOR" for reverse proxy setup
	ipAddr, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ipAddr
}
