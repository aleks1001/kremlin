package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"flag"
)

var (
	router    *mux.Router
)

var routes = Routes{
	Route{
		"AppStatus",
		"GET",
		"/appstatus",
		getAppStatus,
	},
	Route{
		"Root",
		"GET",
		"/",
		getAppStatus,
	},
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func main() {
	ipPtr := flag.String("ip", "0.0.0.0", "localhost")
	portPtr := flag.String("port", "3232", "port to listen on")
	server := newServer(fmt.Sprintf("%v:%v", *ipPtr, *portPtr))
	msg := fmt.Sprintf("Starting server on %s:%s", *ipPtr, *portPtr)
	log.Println(msg)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error: %v", err)
	}
}

func initRouting() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		r.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handlers.CORS()(handler))
	}
	return r
}

func newServer(addr string) *http.Server {
	router = initRouting()
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func getAppStatus(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
