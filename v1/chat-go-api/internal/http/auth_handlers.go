package http

import (
	"chat-go-api/internal/auth"
	"chat-go-api/internal/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email e senha são obrigatórios"})
		return
	}
	u, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro buscando usuario"})
		return
	}
	if u == nil || !u.Active {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PassHas), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}

	token, err := auth.GenerateToken(u.ID, u.Email, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro gerando token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"user": gin.H{
			"id":    u.ID,
			"email": u.Email,
			"role":  u.Role,
		},
	})
}

var _ = db.User{}
