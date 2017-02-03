package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/task/g"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// start http server
func Start() {
	go startHttpServer()
}
func startHttpServer() {
	if !g.Config().Http.Enable {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	// init url mapping
	configCommonRoutes()
	configProcHttpRoutes()
	configIndexHttpRoutes()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http:startHttpServer, ok, listening ", addr)
	log.Fatalln(s.ListenAndServe())
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}
