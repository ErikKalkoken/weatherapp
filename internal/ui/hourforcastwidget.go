package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type HourForecastWidget struct {
	widget.BaseWidget
	hour          *widget.Label
	symbol        *widget.Icon
	temperature   *widget.Label
	precipitation *widget.Label
}

func NewHourForecastWidget() *HourForecastWidget {
	p := widget.NewLabel("")
	p.Importance = widget.HighImportance
	w := &HourForecastWidget{
		hour:          widget.NewLabel(""),
		symbol:        widget.NewIcon(resourceBlankSvg),
		temperature:   widget.NewLabel(""),
		precipitation: p,
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *HourForecastWidget) Set(f forecastHour, icon fyne.Resource) {
	var text string
	if f.isCurrent {
		text = "Now"
	} else {
		text = fmt.Sprintf("%02d", f.time.Hour())
	}
	w.hour.SetText(text)
	w.temperature.SetText(fmt.Sprintf("%.0f°", f.temperature2m))
	w.precipitation.SetText(fmt.Sprintf("%d%%", f.precipitationProbability))
	w.symbol.SetResource(icon)
}

func (w *HourForecastWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewVBox(
		container.NewCenter(w.hour),
		container.NewCenter(w.symbol),
		container.NewCenter(w.temperature),
		container.NewCenter(w.precipitation),
	)
	return widget.NewSimpleRenderer(c)
}
