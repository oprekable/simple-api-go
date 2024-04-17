package httphelper

import (
	"net"
	"net/http"
)

const (
	CHeaderCFConnectingIP = "CF-Connecting-IP"
	CHeaderForwardedFor   = "X-Forwarded-For"
	CHeaderXRealIP        = "X-Real-IP"
)

func GetVisitorIP(r *http.Request) (returnData string) {
	returnData = r.Header.Get(CHeaderCFConnectingIP)
	if returnData != "" && net.ParseIP(returnData) != nil {
		return
	}

	returnData = r.Header.Get(CHeaderForwardedFor)
	if returnData != "" && net.ParseIP(returnData) != nil {
		return
	}

	returnData = r.Header.Get(CHeaderXRealIP)
	if returnData != "" && net.ParseIP(returnData) != nil {
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return
	}

	returnData = net.ParseIP(ip).String()
	return
}
