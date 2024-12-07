package ui

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
	p := widget.NewLabel("")
	p.Importance = widget.HighImportance
	w := &DayForecastWidget{
		day:            widget.NewLabel(""),
		precipitation:  p,
		symbol:         widget.NewIcon(resourceBlankSvg),
		temperatureMax: widget.NewLabel(""),
		temperatureMin: widget.NewLabel(""),
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *DayForecastWidget) Set(f forecastDay, icon fyne.Resource) {
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
	w.symbol.SetResource(icon)
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
