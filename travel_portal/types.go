package main

type Settings struct {
	RequestRatio 	int `json:"request_ratio"`
	Devices 		Devices `json:"devices"`
	Users 			Users `json:"users"`
	TravelType 		TravelType `json:"travel_type"`
}

type Status struct {
	Requests Requests `json:"requests"`
	Error    bool `json:"error"`
}

type Requests struct {
	Total      int `json:"total"`
	Devices    Devices `json:"devices"`
	Users      Users `json:"users"`
	TravelType TravelType `json:"travel_type"`
}

type Devices struct {
	Web    int `json:"web"`
	Mobile int `json:"mobile"`
}

type Users struct {
	Registered int `json:"registered"`
	New        int `json:"new"`
}

type TravelType struct {
	T1 int `json:"t1"`
	T2 int `json:"t2"`
	T3 int `json:"t3"`
}

type PortalStatus struct {
	Name  		string `json:"name"`
	Coordinates []float64 `json:"coordinates"`
	Country 	string `json:"country"`
	Settings 	Settings `json:"settings"`
	Status  	Status `json:"status"`
}

type ResponseError struct {
	Error  string `json:"error,omitempty"`
	Detail string `json:"detail,omitempty"`
}

