package main

import (
	"log"
	"meu-chat/internal/config"
	"meu-chat/internal/database"
	"meu-chat/internal/http/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	c := config.Load()

	db := database.Connect(c)
	_ = db

	r := gin.Default()
	routes.Register(r)

	log.Printf("%s rodando na porta %s", c.AppName, c.AppPort)
	r.Run(":" + c.AppPort)
}
