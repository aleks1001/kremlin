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
	"strconv"
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
	Route{
		"Hotels",
		"POST",
		"/hotels",
		postHotels,
	},
	Route{
		"GetByCityId",
		"GET",
		"/city/{id}",
		getByCityId,
	},
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func init() {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "test",
		MaxOpen:  40,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}

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

func getByCityId(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.ParseFloat(vars["id"], 64)

	var hotels h.Hotels
	res, err := r.Table("hotel").Filter(r.Row.Field("CityId").Eq(id)).OrderBy(r.Desc("OverallGuestRating")).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()
	err = res.All(&hotels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(hotels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getAppStatus(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func postHotels(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	var hotels h.HotelsResponse
	json.Unmarshal(body, &hotels)
	for _, hotel := range hotels.Hotels {
		var h = &h.Hotel{
			CheckIn:hotels.CheckIn,
			CheckOut:hotels.CheckOut,
			HotelId:hotel.HotelId,
			BrandId:hotel.BrandId,
			HotelName:hotel.HotelName,
			NeighborhoodId:hotel.Location["neighborhoodId"].(string),
			NeighborhoodName:hotel.Location["neighborhoodName"].(string),
			CityId:hotel.Location["cityId"].(float64),
			OverallGuestRating:hotel.OverallGuestRating,
			MinPrice:hotel.RatesSummary["minPrice"].(string),
			StarRating:hotel.StarRating,
			TotalReviewCount:hotel.TotalReviewCount,
		}
		_, err = r.DB("test").Table("hotel").Insert(h).RunWrite(session)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func postHotel(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	var hotel h.HotelResponse
	json.Unmarshal(body, &hotel)

	_, err = r.Table("hotel").Insert(hotel).RunWrite(session)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
