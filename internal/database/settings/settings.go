package settings

import (
	"database/sql"
	"log"
	"swagger/internal/database"
	"time"
)

type setting struct {
	Id         int       `json:"id"`
	Path       string    `json:"path"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func CreateOrUpdate(folderPath string) {
	db := database.Connect()
	defer db.Close()

	var id int
	err := db.QueryRow("SELECT id FROM settings").Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Запись не найдена, создаем новую.")
			create(folderPath)
			return
		}
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Таблица settings не пуста. Обновляем запись.")
	update(folderPath, id)
}

func create(folderPath string) {
	setting := setting{
		Path:       folderPath,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	db := database.Connect()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO settings(path, created_at, updated_at) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("Ошибка при подготовке запроса: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		setting.Path,
		setting.Created_at.Format("2006-01-02 15:04:05"),
		setting.Updated_at.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Данные успешно добавлены в таблицу settings")
}

func update(folderPath string, id int) {
	db := database.Connect()
	defer db.Close()

	query := "UPDATE settings SET path = ?, created_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?"

	_, err := db.Exec(query, folderPath, id)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	log.Println("Данные успешно обновлены в таблицу settings")
}

func GetSettings() (string, bool) {
	db := database.Connect()
	defer db.Close()

	var path string
	err := db.QueryRow("SELECT path FROM settings").Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Запись не найдена.")
			return "", false
		}
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Таблица settings не пуста. Возвращаем ID.")
	return path, true
}
