package main

import (
	"fmt"
	"log"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type weatherShort uint

const (
	undefined weatherShort = iota
	cloudy
	dust
	fog
	hail
	hurricane
	lightning
	lightningRainy
	night
	nightPartlyCloudy
	partlyCloudy
	partlyRainy
	partlySnowy
	partlySnowyRainy
	partyLightning
	pouring
	rainy
	snowy
	snowyHeavy
	snowyRainy
	sunny
	windy
)

type weatherCodeMapping struct {
	description string
	short       weatherShort
	shortNight  weatherShort // alternate to be used at night (when defined)
}

var weatherCodeMappings = map[int]weatherCodeMapping{
	0:  {"clear sky", sunny, night},
	1:  {"mainly clear", sunny, night},
	2:  {"partly cloudy", partlyCloudy, nightPartlyCloudy},
	3:  {"overcast", cloudy, undefined},
	45: {"fog", fog, undefined},
	48: {"depositing rime fog", fog, undefined},
	51: {"light drizzle", rainy, rainy},
	52: {"moderate drizzle", rainy, undefined},
	53: {"dense drizzle", rainy, undefined},
	56: {"light freezing drizzle", snowyRainy, undefined},
	57: {"dense freezing drizzle", snowyRainy, undefined},
	61: {"slight rain", rainy, undefined},
	63: {"moderate rain", pouring, undefined},
	65: {"heavy rain", pouring, undefined},
	66: {"light freezing rain", snowyRainy, rainy},
	67: {"heavy freezing rain", pouring, undefined},
	71: {"slight snow fall", snowy, snowy},
	73: {"moderate snow fall", snowy, undefined},
	75: {"heavy snow fall", snowyHeavy, undefined},
	77: {"snow grains", snowy, undefined},
	80: {"slight rain showers", rainy, undefined},
	81: {"moderate rain showers", rainy, undefined},
	83: {"violent rain showers", pouring, undefined},
	85: {"slight snow showers", snowy, undefined},
	86: {"heavy snow showers", snowyHeavy, undefined},
	95: {"thunderstorms", lightning, undefined},
	96: {"thunderstorms with slight hail", lightningRainy, undefined},
	99: {"thunderstorms with heavy hail", hail, undefined},
}

var weatherSymbols map[weatherShort]fyne.Resource

func main() {
	a := app.New()
	w := a.NewWindow("Weather")
	loadWeatherSymbols()
	location, err := getMyLocation()
	if err != nil {
		log.Fatal(err)
	}
	current, hourly, daily, err := getForecast(location.Latitude, location.Longitude)
	if err != nil {
		log.Fatal(err)
	}
	w.SetContent(makeContent(location, current, hourly, daily))
	w.Resize(fyne.NewSize(300, 600))
	w.ShowAndRun()
}

func makeContent(location location, current forecastHour, hourly []forecastHour, daily []forecastDay) *fyne.Container {
	top := makeTopBox(location, current)
	hours := container.NewBorder(
		container.NewStack(canvas.NewRectangle(theme.Color(theme.ColorNameButton)), widget.NewLabel("Hourly Forecast")),
		nil,
		nil,
		nil,
		container.NewHScroll(makeHourGrid(current, hourly)),
	)
	days := container.NewBorder(
		container.NewStack(canvas.NewRectangle(theme.Color(theme.ColorNameButton)), widget.NewLabel("10-Day Forecast")),
		nil,
		nil,
		nil,
		container.NewVScroll(makeDayGrid(daily)),
	)
	c := container.NewBorder(
		container.NewVBox(top, hours),
		nil,
		nil,
		nil,
		days,
	)
	return c
}

func makeHourGrid(current forecastHour, hourly []forecastHour) *fyne.Container {
	grid := container.NewGridWithRows(1)
	for _, v := range slices.Concat([]forecastHour{current}, hourly) {
		grid.Add(NewHourForecastWidget2(v))
	}
	return grid
}

func makeDayGrid(daily []forecastDay) *fyne.Container {
	grid := container.NewGridWithColumns(1)
	for _, v := range daily {
		grid.Add(NewDayForecastWidget2(v))
	}
	return grid
}

func makeTopBox(location location, current forecastHour) *fyne.Container {
	city := fmt.Sprintf("%s / %s", location.City, location.Country)
	t := fmt.Sprintf("# %.0f°", current.temperature2m)
	m := weatherCodeMappings[current.weatherCode]
	x := cases.Title(language.English)
	description := x.String(m.description)
	c := container.NewVBox(
		container.NewCenter(widget.NewLabel(city)),
		container.NewCenter(widget.NewRichTextFromMarkdown(t)),
		container.NewCenter(widget.NewLabel(description)),
	)
	return c
}

func loadWeatherSymbols() {
	weatherSymbols = map[weatherShort]fyne.Resource{
		cloudy:            theme.NewThemedResource(resourceWeatherCloudySvg),
		dust:              theme.NewThemedResource(resourceWeatherDustSvg),
		fog:               theme.NewThemedResource(resourceWeatherFogSvg),
		hail:              theme.NewThemedResource(resourceWeatherHailSvg),
		lightning:         theme.NewThemedResource(resourceWeatherLightningSvg),
		lightningRainy:    theme.NewThemedResource(resourceWeatherLightningRainySvg),
		partyLightning:    theme.NewThemedResource(resourceWeatherPartlyLightningSvg),
		hurricane:         theme.NewThemedResource(resourceWeatherHurricaneOutlineSvg),
		night:             theme.NewThemedResource(resourceWeatherNightSvg),
		nightPartlyCloudy: theme.NewThemedResource(resourceWeatherNightPartlyCloudySvg),
		partlyCloudy:      theme.NewThemedResource(resourceWeatherPartlyCloudySvg),
		partlyRainy:       theme.NewThemedResource(resourceWeatherPartlyRainySvg),
		partlySnowy:       theme.NewThemedResource(resourceWeatherPartlySnowySvg),
		partlySnowyRainy:  theme.NewThemedResource(resourceWeatherPartlySnowyRainySvg),
		pouring:           theme.NewThemedResource(resourceWeatherPouringSvg),
		rainy:             theme.NewThemedResource(resourceWeatherRainySvg),
		snowy:             theme.NewThemedResource(resourceWeatherSnowySvg),
		snowyHeavy:        theme.NewThemedResource(resourceWeatherSnowyHeavySvg),
		snowyRainy:        theme.NewThemedResource(resourceWeatherSnowyRainySvg),
		sunny:             theme.NewThemedResource(resourceWeatherSunnySvg),
		windy:             theme.NewThemedResource(resourceWeatherWindySvg),
		undefined:         theme.QuestionIcon(),
	}
}
