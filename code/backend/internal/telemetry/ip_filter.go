// Package telemetry provides metrics collection and telemetry filtering capabilities.
package telemetry

import (
	"net"
	"net/http"
	"strings"

	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/utils"

	"github.com/gin-gonic/gin"
)

var (
	ginUnauthorizedIPaddress  = gin.H{"error": "unauthorized IP address"}
	ginInvalidClientIPaddress = gin.H{"error": "invalid client ip"}
)

// isIPAllowed reports whether ip is covered by any entry in allowedIPs (plain IP or CIDR).
func isIPAllowed(ip net.IP, rawIP string, allowedIPs []string) bool {
	for _, entry := range allowedIPs {
		if strings.Contains(entry, "/") {
			_, ipNet, err := net.ParseCIDR(entry)
			if err == nil && ipNet.Contains(ip) {
				return true
			}
		} else if entry == rawIP {
			return true
		}
	}
	return false
}

// IPFilterMiddleware restricts endpoint access to whitelisted IPs/CIDRs from env.
func IPFilterMiddleware() gin.HandlerFunc {
	allowedIPsStr := utils.GetEnv("TELEMETRY_ALLOWED_IPS", "")
	if allowedIPsStr == "" {
		logging.Error("No allowed IPs configured. No one can access telemetry", nil)
		return func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusForbidden, ginUnauthorizedIPaddress) }
	}

	allowedIPs := strings.Split(allowedIPsStr, ",")
	for i, ip := range allowedIPs {
		allowedIPs[i] = strings.TrimSpace(ip)
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		parsedClientIP := net.ParseIP(clientIP)
		if parsedClientIP == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, ginInvalidClientIPaddress)
			return
		}
		if !isIPAllowed(parsedClientIP, clientIP, allowedIPs) {
			c.AbortWithStatusJSON(http.StatusForbidden, ginUnauthorizedIPaddress)
			return
		}
		c.Next()
	}
}
