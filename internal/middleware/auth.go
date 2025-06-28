package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims define los claims del JWT
type JWTClaims struct {
	HospitalID uint   `json:"hospital_id"`
	Email      string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware middleware para verificar JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token de autorización requerido",
				"code":    "AUTH_TOKEN_REQUIRED",
				"success": false,
			})
			c.Abort()
			return
		}

		// Verificar formato "Bearer token"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Formato de token inválido",
				"code":    "INVALID_TOKEN_FORMAT",
				"success": false,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Verificar y parsear el token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token inválido",
				"code":    "INVALID_TOKEN",
				"success": false,
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// Agregar información del hospital al contexto
			c.Set("hospital_id", claims.HospitalID)
			c.Set("hospital_email", claims.Email)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token inválido",
				"code":    "INVALID_TOKEN_CLAIMS",
				"success": false,
			})
			c.Abort()
			return
		}
	}
}
