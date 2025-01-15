package apis

import (
	"database/sql"
	"log"
	"swagger/internal/database"
	"time"
)

type api struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	TokenId     int       `json:"token_id"`
	Body        string    `json:"body"`
	Private     string    `json:"private"`
	Entity      string    `json:"entity"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

func CreateOrUpdate(params []string) {
	db := database.Connect()
	defer db.Close()

	obj := api{
		Title:      params[0],
		Entity:     params[1],
		Path:       params[2],
		Method:     params[3],
		Private:    params[4],
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	var id int
	err := db.QueryRow("SELECT id FROM apis WHERE path = ?", params[2]).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Запись не найдена, создаем новую.")
			create(obj)
			return
		}
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Таблица apis не пуста. Обновляем запись.")
	update(id, obj)
}

func create(params api) {
	db := database.Connect()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO apis (title, method, path, private, entity, created_at, updated_at) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatalf("Ошибка при подготовке запроса: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		params.Title,
		params.Method,
		params.Path,
		params.Private,
		params.Entity,
		params.Created_at.Format("2006-01-02 15:04:05"),
		params.Updated_at.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Println("Данные успешно добавлены в таблицу apis")
}

func update(id int, params api) {
	db := database.Connect()
	defer db.Close()

	query := "UPDATE apis SET title=?, method=?, path=?, private=?, entity=?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"

	_, err := db.Exec(query, params.Title, params.Method, params.Path, params.Private, params.Entity, id)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	log.Println("Данные успешно обновлены в таблицу apis")
}
