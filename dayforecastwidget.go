package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DayForecastWidget struct {
	widget.BaseWidget
	day            *widget.Label
	precipitation  *widget.Label
	symbol         *widget.Icon
	temperatureMax *widget.Label
	temperatureMin *widget.Label
}

func NewDayForecastWidget() *DayForecastWidget {
	p := widget.NewLabel("99")
	p.Importance = widget.HighImportance
	w := &DayForecastWidget{
		day:            widget.NewLabel("99"),
		precipitation:  p,
		symbol:         widget.NewIcon(weatherSymbols[undefined]),
		temperatureMax: widget.NewLabel("99"),
		temperatureMin: widget.NewLabel("99"),
	}
	w.ExtendBaseWidget(w)
	return w
}

func NewDayForecastWidget2(f forecastDay) *DayForecastWidget {
	w := NewDayForecastWidget()
	w.Set(f)
	return w
}

func (w *DayForecastWidget) Set(f forecastDay) {
	var text string
	if f.time.Day() == time.Now().UTC().Day() {
		text = "Today"
	} else {
		text = f.time.Weekday().String()
	}
	w.day.SetText(text)
	w.temperatureMin.SetText(fmt.Sprintf("%.0f°", f.temperature2mMin))
	w.temperatureMax.SetText(fmt.Sprintf("%.0f°", f.temperature2mMax))
	w.precipitation.SetText(fmt.Sprintf("%d%%", f.precipitationProbabilityMean))
	m := weatherCodeMappings[f.weatherCode]
	r, ok := weatherSymbols[m.short]
	if !ok {
		r = weatherSymbols[undefined]
	}
	w.symbol.SetResource(r)
}

func (w *DayForecastWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewGridWithColumns(
		5,
		w.day,
		container.NewCenter(w.symbol),
		container.NewCenter(w.precipitation),
		container.NewCenter(w.temperatureMin),
		container.NewCenter(w.temperatureMax),
	)
	return widget.NewSimpleRenderer(c)
}
