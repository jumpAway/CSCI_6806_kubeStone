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

// addRoute is a method on the `engine` struct type that adds a new route to the engine's routing table.
func (e *engine) addRoute(method string, path string, handler HandlerFunc) {
	key := method + "-" + path
	e.router[key] = handler
}

// InitHandler is a method on the `engine` struct type that initializes the router and adds route handlers.
func (e *engine) InitHandler() {
	e.router = make(map[string]HandlerFunc)
	e.addRoute("GET", "/getServer", SearchSer)
	e.addRoute("POST", "/testServer", TestSer)
	e.addRoute("POST", "/addServer", AddSer)
}

/*
ServeHTTP is a method on the `engine` struct type that makes it satisfy the http.Handler interface.

This means an `engine` can be passed to functions like http.ListenAndServe as the main HTTP handler.
*/
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
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
