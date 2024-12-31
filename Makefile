.PHONY: run
run:
	go run cmd/swagger/main.go
// Применение миграций (up):
.PHONY: migrateup
migrateup:
	migrate -database sqlite3://internal/database/swagger.db -path internal/database/migrations up
//Откат миграций (down):
.PHONY: migratedown
migratedown:
	migrate -database sqlite3://internal/database/swagger.db -path internal/database/migrations down
//Откат на одну миграцию:
.PHONY: migratedown_1
migratedown_1:
	migrate -database sqlite3://internal/database/swagger.db -path internal/database/migrations down 1
//Просмотр текущей версии миграции:
.PHONY: migrateversion
migrateversion:
	migrate -database sqlite3://internal/database/swagger.db -path internal/database/migrations version
//Пропуск неисправных миграций:
.PHONY: migrateforce
migrateforce:
	migrate -database sqlite3://internal/database/swagger.db -path internal/database/migrations force <version>