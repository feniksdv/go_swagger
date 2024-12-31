package main

import (
	"swagger/internal/app/layout"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Swagger")

	// Создаем контейнер для отображения текущей страницы
	contentLabel := widget.NewLabel("Выберите пункт меню.")
	contentContainer := container.NewVBox(contentLabel)

	// Функция для обновления содержимого окна
	updateContent := func(newContent fyne.CanvasObject) {
		contentContainer.Objects = []fyne.CanvasObject{newContent}
		contentContainer.Refresh()
	}

	mainMenu := layout.GetMenu(a, w, updateContent)

	w.SetMainMenu(mainMenu)

	w.SetContent(contentContainer)

	w.ShowAndRun()
}
