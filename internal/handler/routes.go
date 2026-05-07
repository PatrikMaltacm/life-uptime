package handler

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.RouterGroup, h *MonitorHandler) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "running"})
	})

	monitor := r.Group("/monitors")
	{
		monitor.GET("/:id", h.GetMonitor)
		monitor.GET("/", h.GetAllMonitors)
		monitor.POST("/", h.CreateMonitor)
	}
}
