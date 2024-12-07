// Package location allows to determine the current location of a machine.
package location

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

type Location struct {
	City      string
	Country   string
	Latitude  float64
	Longitude float64
}

// Get returns the location associated with the IP address of this machine.
func Get() (loc Location, err error) {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return Location{}, fmt.Errorf("making request to IP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return Location{}, fmt.Errorf("making request to IP API %s: %w", resp.Status, err)
	}

	var response ipResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Location{}, fmt.Errorf("decoding response: %w", err)
	}
	if response.Status == "fail" {
		return Location{}, fmt.Errorf(response.Message)
	}
	l := Location{
		Latitude:  response.Lat,
		Longitude: response.Lon,
		City:      response.City,
		Country:   response.Country,
	}
	return l, nil
}
