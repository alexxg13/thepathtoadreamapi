package main

import (
	_ "10042025/docs"
	"10042025/handlers"
	"10042025/middleware"
	"10042025/pkg/postgreSQL"
	"10042025/token"
	"github.com/coalaura/mistral"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
)

// docker run --name my_postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=mypoatgres -p 5432:5432 -d postgres

// @title           THE PATH TO A DREAM API
// @version         1.0
// @description     API для мобильного приложения.
// @host            185.104.114.234:8080
// @BasePath
func main() {

	router := gin.Default()                          // Создаем роутер
	client := mistral.NewMistralClient(token.KEY_AI) // Клиент AI
	db := postgreSQL.ConnectDB()                     //Подключаем DB

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := db.AutoMigrate( // Подключаем Автомиграции БД
		&postgreSQL.User{},
		&postgreSQL.Dialog{},
	); err != nil {
		log.Fatalf("Error AutoMigrate: %s", err)
	}

	router.Use(middleware.LoggerMiddleware()) // Подключаем Middleware

	handler := handlers.NewHandler(db, client) // Подключаем Handlers
	handler.Register(router)                   // Регистрация хэндлера

	// Запускаем роутер
	start(router)

}

func start(router *gin.Engine) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(server.ListenAndServe())
}
