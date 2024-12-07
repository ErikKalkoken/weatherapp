package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ErikKalkoken/weatherapp/internal/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CurrentWeatherWidget struct {
	widget.BaseWidget
	city        *widget.Label
	temperature *widget.RichText
	description *widget.Label
}

func NewCurrentWeatherWidget() *CurrentWeatherWidget {
	w := &CurrentWeatherWidget{
		city:        widget.NewLabel(""),
		temperature: widget.NewRichTextFromMarkdown(""),
		description: widget.NewLabel(""),
	}
	w.ExtendBaseWidget(w)
	return w
}

func NewCurrentWeatherWidget2(location api.Location, current ForecastHour) *CurrentWeatherWidget {
	w := NewCurrentWeatherWidget()
	w.Set(location, current)
	return w
}

func (w *CurrentWeatherWidget) Set(l api.Location, c ForecastHour) {
	city := fmt.Sprintf("%s / %s", l.City, l.Country)
	w.city.SetText(city)
	t := fmt.Sprintf("# %.0f°", c.Temperature2m)
	w.temperature.ParseMarkdown(t)
	m := weatherCodeMappings[c.WeatherCode]
	x := cases.Title(language.English)
	description := x.String(m.description)
	w.description.SetText(description)
}

func (w *CurrentWeatherWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewVBox(
		container.NewCenter(w.city),
		container.NewCenter(w.temperature),
		container.NewCenter(w.description),
	)
	return widget.NewSimpleRenderer(c)
}
