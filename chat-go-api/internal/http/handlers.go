package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"chat-go-api/internal/ai"
	"chat-go-api/internal/db"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	repo *db.Repo
}

func NewHandlers(repo *db.Repo) *Handlers {
	return &Handlers{repo: repo}
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type ChatRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func (h *Handlers) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	intent, err := ai.ClassifyIntent(ctx, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reply, err := h.buildReply(ctx, intent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"intent": intent,
		"reply":  reply,
	})
}

func (h *Handlers) buildReply(ctx context.Context, intent ai.Intent) (string, error) {
	switch intent {
	case ai.Greeting:
		return "👋 Olá! Bem-vindo ao Box! 😄\n\nComo posso ajudar?\n1) Sobre o Box\n2) Horários\n3) Planos\n4) Reservar uma aula", nil

	case ai.About:
		return "🏋️ Somos um Box focado em performance, saúde e comunidade. Quer ver horários ou planos?", nil

	case ai.Plans:
		plans, err := h.repo.ListPlans(ctx)
		if err != nil {
			return "", err
		}
		if len(plans) == 0 {
			return "Ainda não temos planos cadastrados.", nil
		}
		out := "💳 *Planos disponíveis:*\n"
		for _, p := range plans {
			desc := ""
			if p.Description != nil && *p.Description != "" {
				desc = fmt.Sprintf(" — %s", *p.Description)
			}
			out += fmt.Sprintf("- %s: R$ %s/mês%s\n", p.Name, p.Price, desc)
		}
		return out, nil

	case ai.Schedule:
		schedules, err := h.repo.ListSchedules(ctx)
		if err != nil {
			return "", err
		}
		if len(schedules) == 0 {
			return "Ainda não temos horários cadastrados.", nil
		}
		out := "⏰ *Horários disponíveis:*\n"
		for _, s := range schedules {
			end := ""
			if s.EndTime != nil {
				end = fmt.Sprintf(" às %s", s.EndTime.Format("15:04"))
			}
			mod := ""
			if s.Modality != nil && *s.Modality != "" {
				mod = fmt.Sprintf(" (%s)", *s.Modality)
			}
			coach := ""
			if s.Coach != nil && *s.Coach != "" {
				coach = fmt.Sprintf(" - Coach %s", *s.Coach)
			}
			out += fmt.Sprintf("- %s %s%s%s%s\n", s.Weekday, s.StartTime.Format("15:04"), end, mod, coach)
		}
		return out, nil

	case ai.BookClass:
		return "✅ Perfeito! Em breve vamos ativar reservas. Por enquanto posso te mostrar horários e planos 🙂", nil

	default:
		return "Não entendi totalmente 😅\nVocê quer:\n1) Sobre o Box\n2) Horários\n3) Planos\n4) Reservar uma aula", nil
	}
}

type PlanCreateReq struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
	Active      bool    `json:"active"`
}

func (h *Handlers) AdminCreatePlan(c *gin.Context) {
	var req PlanCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	id, err := h.repo.CreatePlan(c.Request.Context(), req.Name, req.Description, req.Price, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handlers) AdminListPlans(c *gin.Context) {
	plans, err := h.repo.ListPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

type ScheduleCreateReq struct {
	Weekday   string  `json:"weekday"`
	StartTime string  `json:"start_time"` // "18:00"
	EndTime   *string `json:"end_time"`
	Modality  *string `json:"modality"`
	Coach     *string `json:"coach"`
	Active    bool    `json:"active"`
}

func (h *Handlers) AdminCreateSchedule(c *gin.Context) {
	var req ScheduleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	id, err := h.repo.CreateSchedule(
		c.Request.Context(),
		req.Weekday,
		req.StartTime,
		req.EndTime,
		req.Modality,
		req.Coach,
		req.Active,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handlers) AdminListSchedules(c *gin.Context) {
	s, err := h.repo.ListSchedules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}
