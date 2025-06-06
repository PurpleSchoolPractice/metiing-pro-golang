Metiing-pro-golang
===============
Meeting Pro Golang — это REST API для управления встречами, написанный на Go.
-------------------------------------------------------------------------------
Требования
--------------------------------------------------------------------------------------
* Go 1.23 или выше
* PostgreSQL 15 или выше
* Библиотеки:
* github.com/jinzhu/gorm (или gorm.io/gorm) для работы с базой данных
* github.com/go-chi/chi для маршрутизации
* golang.org/x/crypto/bcrypt для хеширования паролей
* github.com/DATA-DOG/go-sqlmock для тестирования
* github.com/spf13/cobra  для CLI
* github.com/spf13/viper  для чтения config
Установка
--------------------------------------------------------------------------------------
1. Клонируйте репозиторий:
##
```go
git clone https://github.com/PurpleSchoolPractice/metiing-pro-golang.git
```
2. Установите зависимости:
##

```go
go mod tidy
```
3. Запустите Postgres через Docker. В корне есть файл `docker-compose.yml` для запуска.
##
  * Запустите Docker
  ```go
  docker-compose up -d
  ```
4. Создайте базу данных через новый запрос:
##
```go
CREATE DATABASE ваша название Базы данных(postgres);
```

5. Настройте подключение к базе данных. Создайте `.env` с переменными окружения:
##
```go
DATABASE_DSN="host=localhost user=ваши данные password=ваши данные dbname=postgres port=5432 sslmode=disable"
```
6. Установите приложение
```go
go install
```
Использование
---------------
1. Запустите приложение:
##
```go
meeting 
    или
go run main.go
```
2. API доступно по ссылке <http://localhost:8080>



