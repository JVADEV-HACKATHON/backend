package middleware

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configura CORS para la aplicación
func SetupCORS() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // En producción, especificar dominios específicos
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}

// JSONLoggerMiddleware middleware personalizado para logging
func JSONLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf(`{"time":"%s","method":"%s","path":"%s","protocol":"%s","status_code":%d,"latency":"%s","client_ip":"%s","user_agent":"%s","error_message":"%s"}%s`,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Request.UserAgent(),
			param.ErrorMessage,
			"\n",
		)
	})
}

// ErrorHandlerMiddleware middleware para manejo de errores
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Manejar errores si los hay
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(400, gin.H{
					"error":   "Error de validación en los datos",
					"details": err.Error(),
					"success": false,
					"code":    "VALIDATION_ERROR",
				})
			case gin.ErrorTypePublic:
				c.JSON(500, gin.H{
					"error":   "Error interno del servidor",
					"details": err.Error(),
					"success": false,
					"code":    "INTERNAL_ERROR",
				})
			default:
				c.JSON(500, gin.H{
					"error":   "Error interno del servidor",
					"success": false,
					"code":    "UNKNOWN_ERROR",
				})
			}
		}
	}
}
