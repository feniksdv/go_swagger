package main

import (
	"swagger/internal/app/layout"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Swagger")

	mainMenu := layout.GetMenu(a, w)
	w.SetMainMenu(mainMenu)

	w.ShowAndRun()
}
