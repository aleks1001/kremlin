package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"flag"
	r "gopkg.in/gorethink/gorethink.v4"
	"io/ioutil"

	h "github.com/aleks1001/kremlin/src"
	"encoding/json"
)

var (
	router    *mux.Router
	session   *r.Session
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
	Route{
		"Hotel",
		"POST",
		"/hotel",
		postHotel,
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
	var err error
	session, err = r.Connect(r.ConnectOpts{Address: "localhost:28015"})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	ipPtr := flag.String("ip", "0.0.0.0", "localhost")
	portPtr := flag.String("port", "3232", "port to listen on")
	server := newServer(fmt.Sprintf("%v:%v", *ipPtr, *portPtr))
	msg := fmt.Sprintf("Starting server on %s:%s", *ipPtr, *portPtr)
	log.Println(msg)
	err = server.ListenAndServe()
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

func postHotel(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	var hotels h.Hotels
	json.Unmarshal(body, &hotels)

	fmt.Println(len(hotels.Hotels))

	for _, hotel := range hotels.Hotels {
		_, err = r.DB("test").Table("hotel").Insert(hotel).RunWrite(session)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
