package entity_fields

import (
	"database/sql"
	"log"
	"swagger/internal/database"
	"time"
)

type entityFields struct {
	Id         int       `json:"id"`
	EntityId   int64     `json:"entity_id"`
	FieldName  string    `json:"field_name"`
	FieldDesc  string    `json:"field_desc"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func CreateOrUpdate(entityId int64, fieldName, fileldDesc string) {
	db := database.Connect()
	defer db.Close()

	var id int
	err := db.QueryRow("SELECT id FROM entity_fields WHERE entity_id = ? AND field_name = ?", entityId, fieldName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Запись не найдена, создаем новую.")
			create(entityId, fieldName, fileldDesc)
			return
		}
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Таблица entities не пуста. Обновляем запись.")
	update(id, fileldDesc)
}

func create(entityId int64, fieldName, fileldDesc string) {
	obj := entityFields{
		EntityId:   entityId,
		FieldName:  fieldName,
		FieldDesc:  fileldDesc,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	db := database.Connect()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO entity_fields (entity_id, field_name, field_desc, created_at, updated_at) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Fatalf("Ошибка при подготовке запроса: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		obj.EntityId,
		obj.FieldName,
		obj.FieldDesc,
		obj.Created_at.Format("2006-01-02 15:04:05"),
		obj.Updated_at.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Данные успешно добавлены в таблицу entities")
}

func update(id int, fileldDesc string) {
	db := database.Connect()
	defer db.Close()

	query := "UPDATE entity_fields SET field_desc = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"

	_, err := db.Exec(query, fileldDesc, id)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	log.Println("Данные успешно обновлены в таблицу entities")
}
