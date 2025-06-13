package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	REQUEST_IN_PROGRESS = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_apps_request_In_Progress",
		Help: "no of application request in progress",
	})
)

func main() {
	// Start the application
	startMyApp()
}

func startMyApp() {
	router := mux.NewRouter()
	router.HandleFunc("/user/{name}", func(rw http.ResponseWriter, r *http.Request) {
		REQUEST_IN_PROGRESS.Inc()
		vars := mux.Vars(r)
		name := vars["name"]
		greetings := fmt.Sprintf("Hi Buddy, your name is %s :)", name)
		rw.Write([]byte(greetings))
		time.Sleep(5 * time.Second)
		REQUEST_IN_PROGRESS.Dec()
	}).Methods("GET")

	log.Println("Starting the application server, Guage metrics...")
	router.Path("/metrics").Handler(promhttp.Handler())
	http.ListenAndServe("0.0.0.0:8000", router)
}
