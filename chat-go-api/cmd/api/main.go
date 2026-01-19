package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"chat-go-api/internal/db"
	httpapi "chat-go-api/internal/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// tenta conectar no Postgres com retry
	var sqlDB *sql.DB
	var err error

	// ✅ substitui isso: sqlDB, err := db.NewPool()
	// ✅ por retry:
	retries := 20
	for i := 1; i <= retries; i++ {
		sqlDB, err = db.NewDB()
		if err == nil {
			break
		}
		log.Printf("[startup] aguardando Postgres... tentativa %d/%d", i, retries)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("erro conectando no banco: %v", err)
	}
	defer sqlDB.Close()

	repo := db.NewRepo(sqlDB)

	// cria tabelas
	ctx := context.Background()
	if err := repo.EnsureTables(ctx); err != nil {
		log.Fatalf("erro garantindo tabelas: %v", err)
	}

	r := gin.Default()

	// CORS simples para React
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	h := httpapi.NewHandlers(repo)
	httpapi.RegisterRoutes(r, h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("API running on :%s", port)
	_ = r.Run("0.0.0.0:" + port)
}
