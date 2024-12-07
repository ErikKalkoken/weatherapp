package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ErikKalkoken/weatherapp/internal/ui"
)

const (
	updateTicker = 60 * time.Second
)

func main() {
	a := app.New()
	w := a.NewWindow("Weather")
	u := ui.New(w)
	w.SetContent(u.Content)
	w.Resize(fyne.NewSize(300, 600))
	ticker := time.NewTicker(updateTicker)
	go func() {
		for {
			if err := u.Refresh(); err != nil {
				log.Println("ERROR: ", err)
			}
			<-ticker.C
		}
	}()
	w.ShowAndRun()
}
