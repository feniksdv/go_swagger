package layout

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func buildFileUploadPage(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabel("Выберите папку с проектом:")
	folderPathLabel := widget.NewLabel("Папка не выбрана")
	folderPath := ""

	saveButton := widget.NewButton("Сохранить", func() {
		fmt.Println("Сохранение выполнено.")
	})
	saveButton.Hide()

	selectFolderButton := widget.NewButton("Обзор", func() {
		dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				folderPathLabel.SetText("Ошибка: " + err.Error())
				saveButton.Hide()
				return
			}
			if folder == nil {
				folderPathLabel.SetText("Папка не выбрана")
				saveButton.Hide()
				return
			}
			folderPathLabel.SetText("Выбрана папка: " + folder.Path())
			folderPath = folder.Path()
			saveButton.Show()
		}, w)
	})

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	runProcess := widget.NewButton("Собрать данные", func() {
		if folderPath == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите папку", w)
			return
		}
		progressBar.SetValue(0)
		progressBar.Show()
		go collectData(folderPath, progressBar, w)
	})

	// Основное содержимое страницы
	content := container.NewVBox(
		label,
		folderPathLabel,
		selectFolderButton,
		saveButton,
		runProcess,
		progressBar,
	)

	return content
}
func collectData(folderPath string, progressBar *widget.ProgressBar, w fyne.Window) {
	foundFiles := []string{}

	// Канал для отправки прогресса
	progressChan := make(chan float64)

	// Горутина для обновления прогресс-бара
	go func() {
		for progress := range progressChan {
			progressBar.SetValue(progress)
		}
	}()

	// Выполняем работу в отдельной горутине
	go func() {
		defer close(progressChan) // Закрываем канал после завершения работы

		// Сканируем директории
		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == "Integration" {
				// Добавляем все файлы из папки "Integration"
				files, err := os.ReadDir(path)
				if err != nil {
					return err
				}

				for _, file := range files {
					if !file.IsDir() {
						foundFiles = append(foundFiles, filepath.Join(path, file.Name()))
					}
				}
			}
			return nil
		})

		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Симулируем обновление прогресса
		for i := 1; i <= len(foundFiles); i++ {
			time.Sleep(20 * time.Millisecond) // Задержка для демонстрации
			progressChan <- float64(i) / 100.0
		}

		// Выводим найденные файлы
		result := fmt.Sprintf("Найденные файлы (%d):\n%v", len(foundFiles), foundFiles)
		dialog.ShowInformation("Результат", result, w)
	}()
}
