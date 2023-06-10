package server

import (
	"fmt"
	"kubeStone/m/v2/pkg/config"
	"net/http"
	"strconv"
)

var cfg config.Config

type HandlerFunc func(http.ResponseWriter, *http.Request)
type engine struct {
	router map[string]HandlerFunc
}

func (e *engine) InitHandler() {
	e.router = make(map[string]HandlerFunc)

}
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
