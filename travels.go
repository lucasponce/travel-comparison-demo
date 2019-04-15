package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

type Discount struct {
	User string `json:"user"`
	Discount float32 `json:"discount"`
}

var (
	cityNames = Cities{
		Cities: []string{"paris", "rome", "london"},
	}
	travelInfo = []TravelInfo{
		{
			City: "paris",
			Flights: []Flight{
				{
					Airline: "Red Airlines",
					Price: 300.00,
				},
				{
					Airline: "Green Airlines",
					Price: 350.00,
				},
			},
			Hotels: []Hotel{
				{
					Hotel: "Grand Hotel Paris",
					Price: 1000,
				},
				{
					Hotel: "Little Paris Hotel",
					Price: 800,
				},
			},
			Cars: []Car{
				{
					CarModel: "Sports Car",
					Price: 1000,
				},
				{
					CarModel: "Economy Car",
					Price: 300,
				},
			},
			Insurances: []Insurance{
				{
					Company: "Blue Insurances",
					Price: 250,
				},
				{
					Company: "Yellow Insurances",
					Price: 180,
				},
			},
		},
		{
			City: "rome",
			Flights: []Flight{
				{
					Airline: "Red Airlines",
					Price: 400.00,
				},
				{
					Airline: "Green Airlines",
					Price: 450.00,
				},
			},
			Hotels: []Hotel{
				{
					Hotel: "Grand Hotel Rome",
					Price: 1500,
				},
				{
					Hotel: "Little Rome Hotel",
					Price: 1200,
				},
			},
			Cars: []Car{
				{
					CarModel: "Sports Car",
					Price: 1200,
				},
				{
					CarModel: "Economy Car",
					Price: 500,
				},
			},
			Insurances: []Insurance{
				{
					Company: "Blue Insurances",
					Price: 300,
				},
				{
					Company: "Yellow Insurances",
					Price: 220,
				},
			},
		},
		{
			City: "london",
			Flights: []Flight{
				{
					Airline: "Red Airlines",
					Price: 500.00,
				},
				{
					Airline: "Green Airlines",
					Price: 550.00,
				},
			},
			Hotels: []Hotel{
				{
					Hotel: "Grand Hotel London",
					Price: 2000,
				},
				{
					Hotel: "Little London Hotel",
					Price: 2200,
				},
			},
			Cars: []Car{
				{
					CarModel: "Sports Car",
					Price: 1800,
				},
				{
					CarModel: "Economy Car",
					Price: 1000,
				},
			},
			Insurances: []Insurance{
				{
					Company: "Blue Insurances",
					Price: 380,
				},
				{
					Company: "Yellow Insurances",
					Price: 320,
				},
			},
		},
	}
	localService = "http://localhost:8000"
	flightsService = localService
	hotelsService = localService
	carsService = localService
	insurancesService = localService
	discountsService = localService
	currentService = "travels"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion

	chaosMonkey = false
	chaosMonkeySleep = 500 * time.Millisecond // Milliseconds to wait if chaosMonkey is enabled
	chaosMonkeyCity = ""
	chaosMonkeyUser = ""
)

func main() {
	setupServices()
	glog.Infof("Starting %s \n", instance)
	router := mux.NewRouter()
	router.HandleFunc("/travels", GetCities).Methods("GET")
	router.HandleFunc("/travels/{city}", GetTravelQuote).Methods("GET")
	router.HandleFunc("/flights/{city}", GetFlights).Methods("GET")
	router.HandleFunc("/hotels/{city}", GetHotels).Methods("GET")
	router.HandleFunc("/cars/{city}", GetCars).Methods("GET")
	router.HandleFunc("/insurances/{city}", GetInsurances).Methods("GET")
	router.HandleFunc("/discounts/{user}", GetDiscounts).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func NotFound(w http.ResponseWriter, msg string) {
	glog.Infof("[%s] NotFound: %s \n", instance, msg)

	response, _ := json.Marshal(map[string]string{"error": msg})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(response)
}

func GetCities(w http.ResponseWriter, _ *http.Request) {
	glog.Infof("[%s] GetCities \n", instance)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cityNames)
}

func GetTravelQuote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cityName := params["city"]
	user := r.Header.Get("user")

	glog.Infof("[%s] GetTravelQuote for [city: %s] and [user: %s] \n", instance, cityName, user)

	travelQuote := TravelQuote{
		City: cityName,
		CreatedAt: time.Now().Format(time.RFC3339),
		Status: "Not valid",
	}

	wg := sync.WaitGroup{}
	wg.Add(4)
	errChan := make(chan error, 4)

	go func() {
		defer wg.Done()
		request, _ := http.NewRequest("GET", flightsService + "/flights/" + cityName, nil)
		request.Header.Set("user", user)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			errChan <- err
			return
		}
		if response.StatusCode >= 400 {
			errChan <- errors.New(response.Status)
			return
		}
		json.NewDecoder(response.Body).Decode(&travelQuote.Flights)
	}()

	go func() {
		defer wg.Done()
		request, _ := http.NewRequest("GET", hotelsService + "/hotels/" + cityName, nil)
		request.Header.Set("user", user)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			errChan <- err
		}
		if response.StatusCode >= 400 {
			errChan <- errors.New(response.Status)
			return
		}
		json.NewDecoder(response.Body).Decode(&travelQuote.Hotels)
	}()

	go func() {
		defer wg.Done()
		request, _ := http.NewRequest("GET", carsService + "/cars/" + cityName, nil)
		request.Header.Set("user", user)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			errChan <- err
		}
		if response.StatusCode >= 400 {
			errChan <- errors.New(response.Status)
			return
		}
		json.NewDecoder(response.Body).Decode(&travelQuote.Cars)
	}()

	go func() {
		defer wg.Done()
		request, _ := http.NewRequest("GET", insurancesService + "/insurances/" + cityName, nil)
		request.Header.Set("user", user)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			errChan <- err
		}
		if response.StatusCode >= 400 {
			errChan <- errors.New(response.Status)
			return
		}
		json.NewDecoder(response.Body).Decode(&travelQuote.Insurances)
	}()

	wg.Wait()
	if len(errChan) != 0 {
		NotFound(w, "Travel Quote for " + cityName + " not found")
		return
	}

	travelQuote.Status = "Valid"

	releaseTheMonkey(cityName, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelQuote)
}

func GetFlights(w http.ResponseWriter, r *http.Request) {
	travelInfo, err := getTravelInfo(r)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	travelInfo = applyDiscounts(r, &travelInfo, "flights")

	glog.Infof("[%s] GetFlights for city %s \n", instance, travelInfo.City)

	releaseTheMonkey(travelInfo.City, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelInfo.Flights)
}

func GetHotels(w http.ResponseWriter, r *http.Request) {
	travelInfo, err := getTravelInfo(r)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	travelInfo = applyDiscounts(r, &travelInfo, "hotels")

	glog.Infof("[%s] GetHotels for city %s \n", instance, travelInfo.City)

	releaseTheMonkey(travelInfo.City, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelInfo.Hotels)
}

func GetCars(w http.ResponseWriter, r *http.Request) {
	travelInfo, err := getTravelInfo(r)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	travelInfo = applyDiscounts(r, &travelInfo, "cars")

	glog.Infof("[%s] GetCars for city %s \n", instance, travelInfo.City)

	releaseTheMonkey(travelInfo.City, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelInfo.Cars)
}

func GetInsurances(w http.ResponseWriter, r *http.Request) {
	travelInfo, err := getTravelInfo(r)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	travelInfo = applyDiscounts(r, &travelInfo, "insurances")

	glog.Infof("[%s] GetInsurances for city %s \n", instance, travelInfo.City)

	releaseTheMonkey(travelInfo.City, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelInfo.Insurances)
}

func GetDiscounts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	discount := Discount{
		User: user,
		Discount: 0,
	}
	discountFrom := r.Header.Get("discountFrom")
	if user != "" {
		switch user {
		case "first":
			discount.Discount = 0.10
		case "vip":
			discount.Discount = 0.25
		default:
			discount.Discount = 0.05
		}
	}

	glog.Infof("[%s] GetDiscounts for user %s from %s \n", instance, user, discountFrom)

	releaseTheMonkey("", user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discount)
}

func applyDiscounts(r *http.Request, travelInfo *TravelInfo, discountFrom string) TravelInfo {
	user := r.Header.Get("user")
	if user == "" {
		return *travelInfo
	}

	discount := float32(1)
	request, _ := http.NewRequest("GET", discountsService + "/discounts/" + user, nil)
	request.Header.Set("discountFrom", discountFrom)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		glog.Errorf("No discount. Discount service is not available")
		return *travelInfo
	}
	discountQuote := Discount{}
	json.NewDecoder(response.Body).Decode(&discountQuote)
	discount = discount - discountQuote.Discount

	glog.Infof("[%s] Applying discount %f for %s \n", instance, discount, user)

	for i, flight := range travelInfo.Flights {
		travelInfo.Flights[i].Price = flight.Price * discount
	}
	for i, hotel := range travelInfo.Hotels {
		travelInfo.Hotels[i].Price = hotel.Price * discount
	}
	for i, car := range travelInfo.Cars {
		travelInfo.Cars[i].Price = car.Price * discount
	}
	for i, insurance := range travelInfo.Insurances {
		travelInfo.Insurances[i].Price = insurance.Price * discount
	}
	return *travelInfo
}

func getTravelInfo(r *http.Request) (TravelInfo, error) {
	params := mux.Vars(r)
	cityName := params["city"]
	for _, travel := range travelInfo {
		if travel.City == cityName {
			return deepCopy(travel), nil
		}
	}
	return TravelInfo{}, errors.New("City " + cityName + " not found")
}

func deepCopy(t TravelInfo) TravelInfo {
	out := TravelInfo{}
	out.City = t.City
	out.Flights = make([]Flight, len(t.Flights))
	for i, f := range t.Flights {
		out.Flights[i].Airline = f.Airline
		out.Flights[i].Price = f.Price
	}
	out.Hotels = make([]Hotel, len(t.Hotels))
	for i, h := range t.Hotels {
		out.Hotels[i].Hotel = h.Hotel
		out.Hotels[i].Price = h.Price
	}
	out.Cars = make([]Car, len(t.Cars))
	for i, c := range t.Cars {
		out.Cars[i].CarModel = c.CarModel
		out.Cars[i].Price = c.Price
	}
	out.Insurances = make([]Insurance, len(t.Insurances))
	for i, in := range t.Insurances {
		out.Insurances[i].Company = in.Company
		out.Insurances[i].Price = in.Price
	}
	return out
}

func releaseTheMonkey(city, user string) {
	if chaosMonkey {
		glog.Infof("[%s] ChaosMonkey introduced %s \n", instance, chaosMonkeySleep.String())
		if (city != "" && city == chaosMonkeyCity) || (user != "" && user == chaosMonkeyUser) || (chaosMonkeyCity == "" && chaosMonkeyUser == "") {
			time.Sleep(chaosMonkeySleep)
		}
	}
}

func setupServices() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	ss := os.Getenv("CURRENT_SERVICE")
	if ss != "" {
		currentService = ss
	}
	sv := os.Getenv("CURRENT_VERSION")
	if sv != "" {
		currentVersion = sv
	}

	instance = currentService + "/" + currentVersion

	fs := os.Getenv("FLIGHTS_SERVICE")
	if fs != "" {
		flightsService = fs
	}
	hs := os.Getenv("HOTELS_SERVICE")
	if hs != "" {
		hotelsService = hs
	}
	cs := os.Getenv("CARS_SERVICE")
	if cs != "" {
		carsService = cs
	}
	is := os.Getenv("INSURANCES_SERVICE")
	if is != "" {
		insurancesService = is
	}
	ds := os.Getenv("DISCOUNTS_SERVICE")
	if ds != "" {
		discountsService = ds
	}
	if os.Getenv("CHAOS_MONKEY") == "true" {
		chaosMonkey = true
		sleep := os.Getenv("CHAOS_MONKEY_SLEEP")
		if value, err := strconv.Atoi(sleep); err == nil {
			chaosMonkeySleep = time.Duration(value) * time.Millisecond
		}
		chaosMonkeyCity = os.Getenv("CHAOS_MONKEY_CITY")
		chaosMonkeyUser = os.Getenv("CHAOS_MONKEY_USER")
	}
}