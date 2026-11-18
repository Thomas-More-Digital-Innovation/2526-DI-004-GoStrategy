package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the health status of the server
// @Summary Health check
// @Description Check if the API server is up and running
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *Handler) HealthHandler(c *gin.Context) {
	// Check database connection
	sqlDB, err := db.DB.DB()
	if err != nil {
		core.SendError(c, "Database connection error", http.StatusInternalServerError)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		core.SendError(c, "Database unreachable", http.StatusInternalServerError)
		return
	}

	core.SendJSON(c, models.HealthResponse{Status: "ok", DB: "connected"}, http.StatusOK)
}

// GetCSRFToken returns a message for CSRF token retrieval
// @Summary Get CSRF token
// @Description Retrieve a CSRF token for subsequent state-changing requests
// @Tags auth
// @Produce json
// @Success 200 {object} models.MessageResponse
// @Router /csrf [get]
func (h *Handler) GetCSRFToken(c *gin.Context) {
	// The CSRF token is already set in a cookie by the middleware
	core.SendJSON(c, models.MessageResponse{Message: "CSRF cookie set"}, http.StatusOK)
}
