package main

import (
	"net/http"
	"github.com/golang/glog"
	"flag"
	"os"
	"github.com/gorilla/mux"
	"encoding/json"
)

var (
	// By default travel_control binary expects a /html folder in its directory
	webroot = "./html"
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	wr := os.Getenv("WEBROOT")
	if wr != "" {
		webroot = wr
		glog.Infof("WEBROOT=%s", webroot)
	}
}

func GetStatus(w http.ResponseWriter, _ *http.Request) {
	glog.Info("GetStatus \n")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Everything is going ok")
}

func main() {
	setup()
	glog.Infof("Starting Travel Control")

	router := mux.NewRouter()

	// Dynamic routes

	router.HandleFunc("/status", GetStatus).Methods("GET")

	// Static routes

	// Travel Control console is available under "/console" suffix
	router.Path("/console").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webroot + "/index.html")
	})

	// Static resources are served from "/" path
	// Then index.html doesn't need to adapt the paths
	// i.e.
	// 		<link rel="stylesheet" href="style.css" />
	//		<script src="script.js"></script>
	//		d3.json("countries-110m.json")
	//
	// Don't forget to adjust this if it changes
	// Note that this setup allows to serve static from content and also to run it from built-in server in Goland
	staticFileServer := http.FileServer(http.Dir(webroot))
	router.PathPrefix("/").Handler(staticFileServer)

	glog.Fatal(http.ListenAndServe(":8080", router))
}

