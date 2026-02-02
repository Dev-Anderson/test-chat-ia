package routes

import (
	"meu-chat/internal/http/handlers"
	"meu-chat/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/", handlers.Health)

	auth := r.Group("/admin")
	auth.Use(middleware.Auth(), middleware.OnlyAdmin())
	{
		auth.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin aqui"})
		})
	}
}
