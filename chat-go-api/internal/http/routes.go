package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handlers) {
	r.GET("/", h.Health)

	// auth
	r.POST("/auth/login", h.Login)

	// chat aberto
	r.POST("/chat", h.Chat)

	// admin protegido
	admin := r.Group("/admin")
	admin.Use(AuthRequired())
	{
		admin.GET("/plans", h.AdminListPlans)
		admin.POST("/plans", h.AdminCreatePlan)

		admin.GET("/schedules", h.AdminListSchedules)
		admin.POST("/schedules", h.AdminCreateSchedule)
	}
}
