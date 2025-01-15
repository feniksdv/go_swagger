package layout

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"swagger/internal/database/apis"
	"swagger/internal/database/entities"
	"swagger/internal/database/entity_fields"
	"swagger/internal/database/settings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ProtectedField struct {
	Name string
	Desc string
}

type ParamsApi struct {
	Info    string
	Method  string
	Path    string
	Private string
	Entity  string
}

func buildFileUploadPage(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabel("Выберите папку с проектом:")
	folderPathLabel := widget.NewLabel("Папка не выбрана")
	folderPath := ""
	var runProcess *widget.Button

	if path, found := settings.GetSettings(); found {
		folderPath = path
		log.Printf("Папка определенна = %s", path)

	}

	selectFolderButton := widget.NewButton("Обзор", func() {
		dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				folderPathLabel.SetText("Ошибка: " + err.Error())
				return
			}
			if folder == nil {
				folderPathLabel.SetText("Папка не выбрана")
				return
			}
			folderPathLabel.SetText("Выбрана папка: " + folder.Path())
			folderPath = folder.Path()
		}, w)
	})

	// Отображение начального пути, если он есть
	if folderPath != "" {
		folderPathLabel.SetText("Начальная папка: " + folderPath)
	} else {
		folderPathLabel.SetText("Папка не выбрана")
	}

	saveButton := widget.NewButton("Сохранить", func() {
		if folderPath == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите папку", w)
			return
		}
		settings.CreateOrUpdate(folderPath)
		fmt.Println("Сохранение выполнено.")
	})

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	runProcess = widget.NewButton("Собрать данные", func() {
		if folderPath == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите папку", w)
			return
		}

		progressBar.SetValue(0)
		progressBar.Show()
		runProcess.Hide()
		saveButton.Hide()
		selectFolderButton.Hide()

		func() {
			collectData(folderPath, progressBar, w)

			// После завершения задачи обновляем UI в главном потоке
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Задача завершена",
				Content: "Сканирование данных завершено.",
			})

			w.Content().Refresh() // Обновляем интерфейс
			progressBar.Hide()
			runProcess.Show()
			saveButton.Show()
			selectFolderButton.Show()
		}()
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
	foundFilesIntegration := []string{}
	foundFilesEntity := []string{}

	// Сканируем директории
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Пропускаем обработку каталогов, но продолжаем обход
			return nil
		}

		// Проверяем, находится ли файл в каталоге Integration
		if strings.Contains(path, "/Integration/") {
			foundFilesIntegration = append(foundFilesIntegration, path)
		}

		// Проверяем, находится ли файл в каталоге Entity
		if strings.Contains(path, "/Entity/") {
			foundFilesEntity = append(foundFilesEntity, path)
		}

		return nil
	})

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	// Симулируем обновление прогресса
	go func() {
		for i := 1; i <= len(foundFilesEntity); i++ {
			time.Sleep(500 * time.Millisecond) // Задержка для демонстрации
			progressBar.SetValue(float64(i) / 100.0)
		}
	}()

	// Выводим найденные файлы
	result := fmt.Sprintf("Найденные файлы (%d):(%d):\n%v", len(foundFilesIntegration), len(foundFilesEntity), foundFilesIntegration)
	dialog.ShowInformation("Результат", result, w)

	// сохраняем entity + fields
	addEntity(foundFilesEntity)
	// сохраняем все методы API
	addIntegration(foundFilesIntegration)
}

func addEntity(foundFilesEntity []string) {
	for _, file := range foundFilesEntity {
		info, protectedFields, err := parseEntityFile(file)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		id := entities.CreateOrUpdate(file, info)
		fmt.Printf("Значение info: %s\n", info)
		fmt.Println("Поля protected:")
		for _, field := range protectedFields {
			entity_fields.CreateOrUpdate(id, field.Name, field.Desc)
			fmt.Printf("- %s: %s\n", field.Name, field.Desc)
		}
	}
}

func addIntegration(foundFilesIntegration []string) {
	for _, file := range foundFilesIntegration {
		paramsApi, err := parseIntegrationFile(file)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}

		for _, param := range paramsApi {
			params := []string{
				param.Info,
				param.Entity,
				param.Path,
				param.Method,
				param.Private,
			}
			apis.CreateOrUpdate(params)
		}
	}
}

func parseEntityFile(filePath string) (string, []ProtectedField, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	var infoValue string
	protectedFields := []ProtectedField{}

	// Регулярные выражения для поиска нужных строк
	infoRegex := regexp.MustCompile(`public const info = '(.*?)';`)
	protectedRegex := regexp.MustCompile(`protected \$(\w+);.*?desc:(.*?)(\||@|$)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ищем значение переменной info
		if matches := infoRegex.FindStringSubmatch(line); matches != nil {
			infoValue = matches[1]
		}

		// Ищем protected поля с описанием
		if matches := protectedRegex.FindStringSubmatch(line); matches != nil {
			field := ProtectedField{
				Name: matches[1],
				Desc: matches[2],
			}
			protectedFields = append(protectedFields, field)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, fmt.Errorf("ошибка при чтении файла: %v", err)
	}

	return infoValue, protectedFields, nil
}

func parseIntegrationFile(filePath string) ([]ParamsApi, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	var currentParams *ParamsApi
	var result []ParamsApi

	// Регулярные выражения для поиска нужных строк
	infoRegex := regexp.MustCompile(`public const info = ['"](.*?)['"];`)
	methodRegex := regexp.MustCompile(`(?:^|\s)(?:final\s+)?class\s+(\w+)`)
	privateRegex := regexp.MustCompile(`public const auth = (.*?);`)
	entityRegex := regexp.MustCompile(`\b(?:public|protected|private)\s+(?:\w+\s+)?\$entity\s*=\s*([^;]+);`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.ReplaceAll(line, "\uFEFF", "") // Удаление BOM

		// Если находим начало нового метода, инициализируем новый объект ParamsApi
		if methodMatch := methodRegex.FindStringSubmatch(line); methodMatch != nil {
			// Если текущий объект уже существует, добавляем его в результат
			if currentParams != nil {
				result = append(result, *currentParams)
			}

			// Инициализируем новый объект для нового метода
			currentParams = &ParamsApi{
				Method:  methodMatch[1], // Имя метода
				Path:    filePath,       // Добавляем путь
				Private: "false",
			}
		}

		// Если currentParams инициализирован, добавляем данные
		if currentParams != nil {
			// Ищем информацию о методе
			if infoMatch := infoRegex.FindStringSubmatch(line); infoMatch != nil {
				currentParams.Info = infoMatch[1]
			}

			// Ищем публичность метода
			if privateMatch := privateRegex.FindStringSubmatch(line); privateMatch != nil {
				currentParams.Private = privateMatch[1]
			}

			// Ищем сущность
			if entityMatch := entityRegex.FindStringSubmatch(line); entityMatch != nil {
				currentParams.Entity = entityMatch[1]
			}
		}
	}

	// Если остались данные в currentParams, добавляем их в результат
	if currentParams != nil {
		result = append(result, *currentParams)
	}

	return result, nil
}
