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
	REQUEST_RESPONSE_TIME = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "go_apps_response_latency_sec",
		Help: "response_latency_sec",
	}, []string{"paths"})
)

func main() {
	// Start the application
	startMyApp()
}

func routeMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start_time := time.Now()
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		next.ServeHTTP(w, r)
		time_taken := time.Since(start_time)
		REQUEST_RESPONSE_TIME.WithLabelValues(path).Observe(time_taken.Seconds())
	})

}

func startMyApp() {
	router := mux.NewRouter()
	router.HandleFunc("/user/{name}", func(rw http.ResponseWriter, r *http.Request) {
		// REQUEST_IN_PROGRESS.Inc()
		vars := mux.Vars(r)
		name := vars["name"]
		greetings := fmt.Sprintf("Hi Buddy, your name is %s :)", name)
		rw.Write([]byte(greetings))
		time.Sleep(5 * time.Second)
		// REQUEST_IN_PROGRESS.Dec()
	}).Methods("GET")

	router.Use(routeMiddleware)
	log.Println("Starting the application server, summary metrics...")
	router.Path("/metrics").Handler(promhttp.Handler())
	http.ListenAndServe("0.0.0.0:8000", router)
}
