package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ErikKalkoken/weatherapp/internal/openmeteo"
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

func (w *HourForecastWidget) Set(f openmeteo.ForecastHour, icon fyne.Resource) {
	var text string
	if f.IsCurrent {
		text = "Now"
	} else {
		text = fmt.Sprintf("%02d", f.Time.Hour())
	}
	w.hour.SetText(text)
	w.temperature.SetText(fmt.Sprintf("%.0f°", f.Temperature2m))
	w.precipitation.SetText(fmt.Sprintf("%d%%", f.PrecipitationProbability))
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
