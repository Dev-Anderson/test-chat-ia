package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"chat-go-api/internal/db"
	httpapi "chat-go-api/internal/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	//cria admin inicial se nao existir
	adminEmail := strings.TrimSpace(strings.ToLower(os.Getenv("ADMIN_EMAIL")))
	adminPass := strings.TrimSpace(os.Getenv("ADMIN_PASS"))

	if adminEmail != "" && adminPass != "" {
		u, err := repo.GetUserByEmail(ctx, adminEmail)
		if err != nil {
			log.Fatalf("erro buscando admin inicial: %v", err)
		}

		if u == nil {
			hash, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
			_, err := repo.CreateUser(ctx, adminEmail, string(hash), "admin")
			if err != nil {
				log.Fatalf("erro criando admin: %v", err)
			}
			log.Printf("[startup] admin criado: %s", adminEmail)
		}
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
