package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handlers) {
	r.GET("/", h.Health)

	r.POST("/chat", h.Chat)

	admin := r.Group("/admin")
	{
		admin.GET("/plans", h.AdminListPlans)
		admin.POST("/plans", h.AdminCreatePlan)

		admin.GET("/schedules", h.AdminListSchedules)
		admin.POST("/schedules", h.AdminCreateSchedule)
	}
}
