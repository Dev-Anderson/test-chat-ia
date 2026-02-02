package database

import (
	"fmt"
	"log"
	"meu-chat/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
	)

	var db *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("database connected")
			return db
		}

		log.Printf("database not ready (attempt %d/10). retrying...", i)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("failed to connect database after multiple attempts:", err)
	return nil
}
