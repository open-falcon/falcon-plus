package http

import (
	"github.com/open-falcon/agent/g"
	"log"
	"net/http"
)

func init() {
	initHealthRoutes()
}

func Start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	if g.Config().Debug {
		log.Println("listening", addr)
	}

	log.Fatalln(s.ListenAndServe())
}
