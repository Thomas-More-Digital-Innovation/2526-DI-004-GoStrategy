// Package handlers provides handlers for the REST API.
package handlers

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// DebugStats returns runtime memory and goroutine statistics
func (h *Handler) DebugStats(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	h.Mutex.RLock()
	sessionCount := len(h.Sessions)
	h.Mutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"goroutines":     runtime.NumGoroutine(),
		"activeSessions": sessionCount,
		"memory": gin.H{
			"alloc":      m.Alloc / 1024 / 1024,
			"totalAlloc": m.TotalAlloc / 1024 / 1024,
			"sys":        m.Sys / 1024 / 1024,
			"numGC":      m.NumGC,
		},
	})
}
