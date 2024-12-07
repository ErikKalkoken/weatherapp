package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

type iconName uint

const (
	undefined iconName = iota
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
	icon        iconName
	iconNight   iconName // alternate to be used at night (when defined)
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

var weatherIcons map[iconName]fyne.Resource

func main() {
	a := app.New()
	loadWeatherIcons()
	w := a.NewWindow("Weather")
	u := newUI(w)
	w.SetContent(u.content)
	w.Resize(fyne.NewSize(300, 600))
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			u.refresh()
			<-ticker.C
		}
	}()
	w.ShowAndRun()
}

func loadWeatherIcons() {
	weatherIcons = map[iconName]fyne.Resource{
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
