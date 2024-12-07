package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ErikKalkoken/weatherapp/internal/api"
)

const (
	forecastedHours = 24
	forecastedDays  = 10
)

type ui struct {
	Content fyne.CanvasObject
	window  fyne.Window

	current *CurrentWeatherWidget
	hours   []*HourForecastWidget
	days    []*DayForecastWidget
}

func New(w fyne.Window) *ui {
	u := &ui{
		window:  w,
		current: NewCurrentWeatherWidget(),
		hours:   make([]*HourForecastWidget, forecastedHours+1),
		days:    make([]*DayForecastWidget, forecastedDays),
	}

	hoursGrid := container.NewGridWithRows(1)
	for i := range forecastedHours + 1 {
		h := NewHourForecastWidget()
		hoursGrid.Add(h)
		u.hours[i] = h
	}
	hoursBox := container.NewBorder(
		makeTitle("Hourly Forecast"),
		nil,
		nil,
		nil,
		container.NewHScroll(hoursGrid),
	)

	dayGrid := container.NewGridWithColumns(1)
	for i := range forecastedDays {
		d := NewDayForecastWidget()
		dayGrid.Add(d)
		u.days[i] = d
	}
	daysBox := container.NewBorder(
		makeTitle("10-Day Forecast"),
		nil,
		nil,
		nil,
		container.NewVScroll(dayGrid),
	)
	c := container.NewBorder(
		container.NewVBox(u.current, hoursBox),
		nil,
		nil,
		nil,
		daysBox,
	)
	u.Content = c
	return u
}

func (u *ui) Refresh() {
	var err error
	loc, err := api.GetMyLocation()
	if err != nil {
		log.Fatal(err)
	}
	current, hours, days, err := getForecast(loc.Latitude, loc.Longitude)
	if err != nil {
		log.Fatal(err)
	}
	u.current.Set(loc, current)
	u.hours[0].Set(current, resourceBlankSvg)
	for i, f := range hours {
		if i+1 >= len(u.hours) {
			break
		}
		u.hours[i+1].Set(f, iconFromCode(f.WeatherCode, f.IsDay))
	}
	for i, f := range days {
		if i >= len(u.days) {
			break
		}
		u.days[i].Set(f, iconFromCode(f.WeatherCode, true))
	}
}

func makeTitle(s string) *fyne.Container {
	return container.NewStack(canvas.NewRectangle(theme.Color(theme.ColorNameButton)), widget.NewLabel(s))
}
