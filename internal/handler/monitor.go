package handler

import (
	"net/http"

	"github.com/PatrikMaltacm/life-uptime/internal/model"
	"github.com/PatrikMaltacm/life-uptime/internal/repository"
	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	repo *repository.MonitorRepository
}

func NewMonitorHandler(repo *repository.MonitorRepository) *MonitorHandler {
	return &MonitorHandler{repo: repo}
}

func (h *MonitorHandler) GetMonitor(ctx *gin.Context) {
	id := ctx.Param("id")

	monitor, err := h.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Monitor não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, monitor)
}

func (h *MonitorHandler) GetAllMonitors(ctx *gin.Context) {
	monitors, err := h.repo.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Dados não encontrados"})
		return
	}
	ctx.JSON(http.StatusOK, monitors)
}

func (h *MonitorHandler) CreateMonitor(ctx *gin.Context) {
	var payload model.MonitorRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validação falhou",
			"details": err.Error(),
		})
		return
	}

	res, err := h.repo.Create(ctx.Request.Context(), payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar monitor"})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
