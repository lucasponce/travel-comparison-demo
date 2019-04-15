package main

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Cities struct {
	Cities []string `json:"cities"`
}

type Flight struct {
	Airline string `json:"airline"`
	Price float32 `json:"price"`
}

type Hotel struct {
	Hotel string `json:"hotel"`
	Price float32 `json:"price"`
}

type Car struct {
	CarModel string `json:"carModel"`
	Price float32 `json:"price"`
}

type Insurance struct {
	Company string `json:"company"`
	Price float32 `json:"price"`
}

type TravelInfo struct {
	City string `json:"city"`
	Flights []Flight `json:"flights"`
	Hotels []Hotel `json:"hotels"`
	Cars []Car `json:"cars"`
	Insurances []Insurance `json:"insurances"`
}

type TravelQuote struct {
	City string `json:"city"`
	CreatedAt string `json:"createdAt"`
	Status string `json:"status"`
	Flights []Flight `json:"flights"`
	Hotels []Hotel `json:"hotels"`
	Cars []Car `json:"cars"`
	Insurances []Insurance `json:"insurances"`
}

var (
	currentPortal = ""
	travelsService = "http://localhost:8000/travels"
	cityPortal = ""
	userType = ""
	requestSleep = 500 * time.Millisecond // Milliseconds to wait for a new request
)

func main() {
	setupServices()
	glog.Infof("Starting Portal %s \n", currentPortal)

	requestId := 0
	for {
		request, _ := http.NewRequest("GET", travelsService + "/travels/" + cityPortal, nil)
		request.Header.Set("user", userType)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			glog.Errorf("[%s] Request %d for city %s returned error: %s", currentPortal, requestId, cityPortal, err.Error())
			continue
		}
		if response.StatusCode >= 400 {
			glog.Errorf("[%s] Request %d for city %s returned error: %s",currentPortal, requestId, cityPortal, response.Status)
			continue
		}
		travelQuote := TravelQuote{}
		json.NewDecoder(response.Body).Decode(&travelQuote)
		glog.Infof("[%s] Quote received %d for %s. %v",currentPortal, requestId, cityPortal, travelQuote)
		time.Sleep(requestSleep)
	}
}

func setupServices() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	currentPortal = os.Getenv("CURRENT_PORTAL")
	travelsService = os.Getenv("TRAVELS_SERVICE")
	cityPortal = os.Getenv("CITY_PORTAL")
	userType = os.Getenv("USER_TYPE")
	sleep := os.Getenv("REQUEST_SLEEP")
	if value, err := strconv.Atoi(sleep); err == nil {
		requestSleep = time.Duration(value) * time.Millisecond
	}
}