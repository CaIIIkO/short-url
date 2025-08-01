include .env
export $(shell sed 's/=.*//' .env)


#============Docker============
SERVICE=app
# Собрать образы
build:
	docker-compose build $(SERVICE)

# Запустить контейнеры (в фоне) с пересборкой образа
up:
	docker-compose up -d --build --force-recreate $(SERVICE)

# Остановить контейнеры
down:
	docker-compose down

# Перезапустить контейнер (с пересборкой и пересозданием)
restart: down up

# Просмотр логов сервиса
logs:
	docker-compose logs -f $(SERVICE)


#============МИГРАЦИИ============
goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "$(DATABASE_DSN_MIGRATIONS)" create rename_me sql

goose-up:
	goose -dir ./migrations postgres "$(DATABASE_DSN_MIGRATIONS)" up

goose-down:
	goose -dir ./migrations postgres "$(DATABASE_DSN_MIGRATIONS)" down

goose-status:
	goose -dir ./migrations postgres "$(DATABASE_DSN_MIGRATIONS)" status
