package http

import (
	"net/http"
	"os"
)

func initKernelRoutes() {
	http.HandleFunc("/proc/kernel/hostname", hostnameHandler)
}

func hostnameHandler(w http.ResponseWriter, r *http.Request) {
	name, err := os.Hostname()
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}

	RenderDataJson(w, name)
}
