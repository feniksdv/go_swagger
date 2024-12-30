package layout

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/joho/godotenv"
)

func GetMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	file := fyne.NewMenu("Файл",
		fyne.NewMenuItem("Импорт", func() {
			//TODO
		}),
		fyne.NewMenuItem("Экспорт", func() {
			//TODO
		}),
		fyne.NewMenuItem("Закрыть", func() {
			a.Quit()
		}),
	)

	settings := fyne.NewMenu("Настройки",
		fyne.NewMenuItem("Сгененировать ключ", func() {
			//TODO
		}),
		fyne.NewMenuItem("Настроить", func() {
			//TODO
		}),
	)

	mainMenu := fyne.NewMainMenu(file, settings)

	width, heigth := getSizeForWindow()

	w.Resize(fyne.NewSize(width, heigth))

	return mainMenu
}

func getSizeForWindow() (float32, float32) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	width, err := strconv.ParseFloat(os.Getenv("WIDTH"), 32)
	if err != nil {
		fmt.Printf("\nПреобразование из строки width в Float32 - ошибка: %s", err.Error())
	}

	heigth, err := strconv.ParseFloat(os.Getenv("HEIGTH"), 32)
	if err != nil {
		fmt.Printf("\nПреобразование из строки heigth в Float32 - ошибка: %s", err.Error())
	}

	return float32(width), float32(heigth)
}
