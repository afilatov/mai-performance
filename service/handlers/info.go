package handlers

import (
	"net"
	"net/http"
)

type infoHandler struct {
	ip string
}

func newInfoHandler() *infoHandler {
	ip, _ := getLocalIp()
	return &infoHandler{ip: ip}
}

type infoResponse struct {
	Status string `json:"status"`
	IP string `json:"ip"`
}

func (h *infoHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := sendResponse(w, http.StatusOK, infoResponse{Status: "success", IP: h.ip}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}
