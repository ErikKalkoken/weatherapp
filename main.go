package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ErikKalkoken/weatherapp/internal/ui"
)

func main() {
	a := app.New()
	ui.LoadWeatherIcons()
	w := a.NewWindow("Weather")
	u := ui.New(w)
	w.SetContent(u.Content)
	w.Resize(fyne.NewSize(300, 600))
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			u.Refresh()
			<-ticker.C
		}
	}()
	w.ShowAndRun()
}
