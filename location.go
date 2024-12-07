package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ipResponse struct {
	City        string
	Country     string
	CountryCode string
	Lat         float64
	Lon         float64
	Message     string
	Region      string
	RegionName  string
	Status      string
	Timezone    string
	Zip         string
}

type location struct {
	City      string
	Country   string
	Latitude  float64
	Longitude float64
}

// getMyLocation returns the current location from the IP address of this machine.
func getMyLocation() (loc location, err error) {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return location{}, fmt.Errorf("making request to IP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return location{}, fmt.Errorf("making request to IP API %s: %w", resp.Status, err)
	}

	var response ipResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return location{}, fmt.Errorf("decoding response: %w", err)
	}
	if response.Status == "fail" {
		return location{}, fmt.Errorf(response.Message)
	}
	l := location{
		Latitude:  response.Lat,
		Longitude: response.Lon,
		City:      response.City,
		Country:   response.Country,
	}
	return l, nil
}
