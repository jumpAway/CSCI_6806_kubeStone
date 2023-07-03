package server

import (
	"fmt"
	"kubeStone/pkg/config"
	"net/http"
	"strconv"
)

var cfg config.Config

type HandlerFunc func(http.ResponseWriter, *http.Request)
type engine struct {
	router map[string]HandlerFunc
}

func (e *engine) addRoute(method string, path string, handler HandlerFunc) {
	key := method + "-" + path
	e.router[key] = handler
}

func (e *engine) InitHandler() {
	e.router = make(map[string]HandlerFunc)
	e.addRoute("GET", "/getServer", SearchSer)
	e.addRoute("POST", "/testServer", TestSer)
	e.addRoute("POST", "/addServer", AddSer)
	e.addRoute("POST", "/getClusterNS", getClusterNS)
	e.addRoute("POST", "/getClusterRes", getClusterRes)
	e.addRoute("GET", "/getCluster", searchCluster)
	e.addRoute("POST", "/createCluster", CreateCluster)
	e.addRoute("POST", "/byGPT", byGPT)
	e.addRoute("POST", "/gptHistory", GptHistory)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //CORS一定要设置在ServeHTTP中
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" { //一定要有这个判断
		w.WriteHeader(http.StatusOK)
		return
	}

	key := r.Method + "-" + r.URL.Path
	if handle, exist := e.router[key]; exist {
		handle(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("not found")); err != nil {
			return
		}
	}
}
func Run() error {
	var err error
	cfg, err = config.InitConfig()
	if err != nil {
		fmt.Println("Error initializing Config", err)
		return err
	}

	engine := new(engine)
	engine.InitHandler()

	err = http.ListenAndServe(":"+strconv.Itoa(cfg.Server.Port), engine)
	return err
}
