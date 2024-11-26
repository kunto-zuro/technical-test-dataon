package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"technical-test-dataon/models"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=paramadaksa dbname=tree_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop tabel lama (opsional)
	//err = DB.Migrator().DropTable(&models.Node{})
	//if err != nil {
	//	log.Fatal("Failed to drop table:", err)
	//}

	err = DB.AutoMigrate(&models.Node{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	return DB
}
