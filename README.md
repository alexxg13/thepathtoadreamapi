Перед запуском программы вы должны вставить свой токен в переменную KEY_AI по адресу: token/token.go
Запустить докер с PostgreSQL 
docker run --name my_postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=mypostgres -p 5432:5432 -d postgres

После всех этих действий можем запускать проект
go run main.go
 
