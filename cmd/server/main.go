package main

import (
	"log"
	"os"

	"hospital-api/internal/config"
	"hospital-api/internal/database"
	"hospital-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Configurar modo de Gin
	gin.SetMode(cfg.Server.GinMode)

	// Conectar a la base de datos
	database.ConnectDatabase()

	// Configurar rutas
	router := routes.SetupRoutes()

	// Obtener puerto del entorno o usar el de configuración
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("🚀 Servidor iniciando en puerto %s", port)
	log.Printf("🏥 Hospital API v1.0.0")
	log.Printf("📋 Endpoints disponibles:")
	log.Printf("   🔐 POST /api/v1/auth/login")
	log.Printf("   👤 GET  /api/v1/auth/profile")
	log.Printf("   🧑‍⚕️ CRUD /api/v1/pacientes")
	log.Printf("   📊 CRUD /api/v1/historial")
	log.Printf("   🦠 GET  /api/v1/epidemiologia/stats")
	log.Printf("   🗺️  GET  /api/v1/epidemiologia/contagious")
	log.Printf("   ✅ GET  /api/v1/health")

	// Iniciar servidor
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
