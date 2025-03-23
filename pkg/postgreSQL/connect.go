package postgreSQL

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDB() *gorm.DB {
	dsn := "host=localhost user=admin password=admin dbname=mypoatgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Filed to connect db: %s", err)
	}
	return db
}
