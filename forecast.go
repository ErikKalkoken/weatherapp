package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type forecastHour struct {
	time                     time.Time
	temperature2m            float64
	precipitationProbability int
	weatherCode              int
	isDay                    bool
	isCurrent                bool
}

type forecastDay struct {
	time                         time.Time
	temperature2mMin             float64
	temperature2mMax             float64
	precipitationProbabilityMean int
	weatherCode                  int
}

type forecastResponse struct {
	Elevation            float64 `json:"elevation"`
	GenerationTimeMS     float64 `json:"generationtime_ms"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Error                bool    `json:"error"`
	Reason               string  `json:"reason"`

	Current      map[string]any    `json:"current"`
	CurrentUnits map[string]string `json:"current_units"`
	Daily        map[string][]any  `json:"daily"`
	DailyUnits   map[string]string `json:"daily_units"`
	Hourly       map[string][]any  `json:"hourly"`
	HourlyUnits  map[string]string `json:"hourly_units"`
}

// getMyLocation returns the current location from the IP address of this machine.
func getForecast(lat float64, lon float64) (forecastHour, []forecastHour, []forecastDay, error) {
	response, err := fetchData(lat, lon)
	if err != nil {
		return forecastHour{}, nil, nil, err
	}
	current, err := parseCurrent(response)
	if err != nil {
		return forecastHour{}, nil, nil, err
	}
	vv, err := parseHourly(response)
	if err != nil {
		return forecastHour{}, nil, nil, err
	}
	hourly := make([]forecastHour, 0)
	for _, v := range vv {
		if v.time.After(time.Now().UTC().Truncate(time.Hour).Add(time.Hour)) {
			hourly = append(hourly, v)
		}
		if len(hourly) == 24 {
			break
		}
	}
	daily, err := parseDaily(response)
	if err != nil {
		return forecastHour{}, nil, nil, err
	}
	return current, hourly, daily, nil
}

func fetchData(lat float64, lon float64) (forecastResponse, error) {
	v := url.Values{}
	v.Add("latitude", fmt.Sprint(lat))
	v.Add("longitude", fmt.Sprint(lon))
	v.Add("timezone", "GMT")
	v.Add("forecast_days", "10")
	v.Add("current", "temperature_2m,precipitation_probability,weather_code,is_day")
	v.Add("daily", "temperature_2m_max,temperature_2m_min,precipitation_probability_mean,weather_code")
	v.Add("hourly", "temperature_2m,precipitation_probability,weather_code,is_day")
	u := "https://api.open-meteo.com/v1/forecast/?" + v.Encode()
	resp, err := http.Get(u)
	if err != nil {
		return forecastResponse{}, fmt.Errorf("making request to open meteo API: %w", err)
	}
	defer resp.Body.Close()

	var response forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("decoding response: %w", err)
	}
	if response.Error {
		return forecastResponse{}, fmt.Errorf("Error from open meteo: %s", response.Reason)
	}

	return response, nil
}

func parseCurrent(response forecastResponse) (forecastHour, error) {
	c := forecastHour{isCurrent: true}
	x, err := time.Parse("2006-01-02T15:04", response.Current["time"].(string))
	if err != nil {
		return c, err
	}
	c.time = x.UTC()
	v1, ok := response.Current["temperature_2m"]
	if !ok {
		return c, fmt.Errorf("missing data")
	}
	c.temperature2m = v1.(float64)
	v2, ok := response.Current["precipitation_probability"]
	if !ok {
		return c, fmt.Errorf("missing data")
	}
	c.precipitationProbability = int(v2.(float64))
	v3, ok := response.Current["weather_code"]
	if !ok {
		return c, fmt.Errorf("missing data")
	}
	c.weatherCode = int(v3.(float64))
	v4, ok := response.Current["is_day"]
	if !ok {
		return c, fmt.Errorf("missing data")
	}
	c.isDay = v4.(float64) == 1.0
	return c, nil
}

func parseHourly(response forecastResponse) ([]forecastHour, error) {
	hourly := make([]forecastHour, len(response.Hourly["time"]))
	for i, v := range response.Hourly["time"] {
		x, err := time.Parse("2006-01-02T15:04", v.(string))
		if err != nil {
			return nil, err
		}
		hourly[i].time = x.UTC()
	}
	vv1, ok := response.Hourly["temperature_2m"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv1 {
		hourly[i].temperature2m = v.(float64)
	}
	vv2, ok := response.Hourly["precipitation_probability"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv2 {
		hourly[i].precipitationProbability = int(v.(float64))
	}
	vv3, ok := response.Hourly["weather_code"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv3 {
		hourly[i].weatherCode = int(v.(float64))
	}
	vv4, ok := response.Hourly["is_day"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv4 {
		hourly[i].isDay = v.(float64) == 1.0
	}
	return hourly, nil
}

func parseDaily(response forecastResponse) ([]forecastDay, error) {
	daily := make([]forecastDay, len(response.Daily["time"]))
	for i, v := range response.Daily["time"] {
		x, err := time.Parse("2006-01-02", v.(string))
		if err != nil {
			return nil, err
		}
		daily[i].time = x.UTC()
	}
	vv1, ok := response.Daily["temperature_2m_min"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv1 {
		daily[i].temperature2mMin = v.(float64)
	}
	vv2, ok := response.Daily["temperature_2m_max"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv2 {
		daily[i].temperature2mMax = v.(float64)
	}
	vv3, ok := response.Daily["precipitation_probability_mean"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv3 {
		daily[i].precipitationProbabilityMean = int(v.(float64))
	}
	vv4, ok := response.Daily["weather_code"]
	if !ok {
		return nil, fmt.Errorf("missing data")
	}
	for i, v := range vv4 {
		daily[i].weatherCode = int(v.(float64))
	}
	return daily, nil
}
