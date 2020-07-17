package main

import (
	"net/http"
	"github.com/golang/glog"
	"flag"
	"os"
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
)

var (
	// By default travel_control binary expects a /html folder in its directory
	webroot = "./html"
	statusDemoIndex = 0
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

func response(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		response, _ = json.Marshal(ResponseError{Error: err.Error()})
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func GetStatus(w http.ResponseWriter, _ *http.Request) {
	glog.Info("GetStatus \n")
	response(w, http.StatusOK, demoData())
}

func PutSettings(w http.ResponseWriter, r *http.Request) {
	glog.Info("PutSettings \n")

	params := mux.Vars(r)
	city := params["city"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, ResponseError{Error: "Error reading body", Detail: err.Error()})
		return
	}

	var settings Settings
	if err := json.Unmarshal(body, &settings); err != nil {
		response(w, http.StatusBadRequest, ResponseError{Error: "Error unmarshall settings", Detail: err.Error()})
		return
	}

	glog.Infof("Received PutSettings for [%s] [%v]", city, settings)
	response(w, http.StatusOK, settings)
}

func main() {
	setup()
	glog.Infof("Starting Travel Control")

	router := mux.NewRouter()

	// Dynamic routes

	router.HandleFunc("/status", GetStatus).Methods("GET")
	router.HandleFunc("/settings/{city}", PutSettings).Methods("PUT")

	// Static routes

	// Travel Control console is available under "/console" suffix
	router.Path("/console").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webroot+"/index.html")
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

func demoData() []CityStatus {

	status1 := Status{
		Requests{
			100,
			Devices{
				75,
				25,
			},
			Users{
				10,
				90,
			},
			TravelType{
				60,
				20,
				20,
			},
		},
	}

	status2 := Status{
		Requests{
			300,
			Devices{
				175,
				125,
			},
			Users{
				110,
				190,
			},
			TravelType{
				150,
				50,
				100,
			},
		},
	}

	status3 := Status{
		Requests{
			200,
			Devices{
				150,
				50,
			},
			Users{
				50,
				150,
			},
			TravelType{
				100,
				80,
				20,
			},
		},
	}

	status4 := Status{
		Requests{
			1000,
			Devices{
				750,
				250,
			},
			Users{
				100,
				900,
			},
			TravelType{
				600,
				200,
				200,
			},
		},
	}

	status := []Status{status1, status2, status3, status4}

	statusDemoIndex = statusDemoIndex + 1;

	return []CityStatus{
		{
			"Paris",
			[]float64{2.337418, 48.861310},
			"France",
			Settings{
				100,
				Devices{
					75,
					25,
				},
				Users{
					10,
					90,
				},
				TravelType{
					60,
					20,
					20,
				},
			},
			status[statusDemoIndex % len(status)],
		},
		{
			"Rome",
			[]float64{12.492194, 41.890668},
			"Italy",
			Settings{
				100,
				Devices{
					75,
					25,
				},
				Users{
					10,
					90,
				},
				TravelType{
					60,
					20,
					20,
				},
			},
			status[(statusDemoIndex+1)%len(status)],
		},
		{
			"London",
			[]float64{-0.128018, 51.508178},
			"United Kingdom",
			Settings{
				100,
				Devices{
					75,
					25,
				},
				Users{
					10,
					90,
				},
				TravelType{
					60,
					20,
					20,
				},
			},
			status[(statusDemoIndex+2)%len(status)],
		},
	}
}
