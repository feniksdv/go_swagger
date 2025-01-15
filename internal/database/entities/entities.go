package entities

import (
	"database/sql"
	"log"
	"swagger/internal/database"
	"time"
)

type entities struct {
	Id         int       `json:"id"`
	Entity     string    `json:"entity"`
	Info       string    `json:"info"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func CreateOrUpdate(entity, info string) int64 {
	db := database.Connect()
	defer db.Close()

	var id int64
	err := db.QueryRow("SELECT id FROM entities WHERE entity = ?", entity).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Запись не найдена, создаем новую.")
			id = create(entity, info)
			return id
		}
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Таблица entities не пуста. Обновляем запись.")
	update(entity, info, id)
	return id
}

func create(entity, info string) int64 {
	obj := entities{
		Entity:     entity,
		Info:       info,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	db := database.Connect()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO entities (entity, info, created_at, updated_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Ошибка при подготовке запроса: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		obj.Entity,
		obj.Info,
		obj.Created_at.Format("2006-01-02 15:04:05"),
		obj.Updated_at.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	// Получение ID созданной записи
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Ошибка получения ID последней вставки: %v", err)
	}

	log.Printf("Данные успешно добавлены в таблицу entities, ID: %d\n", lastInsertID)
	return lastInsertID
}

func update(entity, info string, id int64) {
	db := database.Connect()
	defer db.Close()

	query := "UPDATE entities SET entity = ?, info = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"

	_, err := db.Exec(query, entity, info, id)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	log.Println("Данные успешно обновлены в таблицу entities")
}
