package routes

import (
	"hospital-api/internal/handlers"
	"hospital-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes() *gin.Engine {
	// Crear router de Gin
	router := gin.New()

	// Middleware globales
	router.Use(middleware.JSONLoggerMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.SetupCORS())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Crear instancias de handlers
	authHandler := handlers.NewAuthHandler()
	pacienteHandler := handlers.NewPacienteHandler()
	historialHandler := handlers.NewHistorialHandler()

	// Rutas públicas
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Hospital API is running",
				"version": "1.0.0",
			})
		})

		// Autenticación (sin middleware de auth)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}
	}

	// Rutas protegidas (requieren autenticación)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Perfil del hospital autenticado
		auth := protected.Group("/auth")
		{
			auth.GET("/profile", authHandler.GetProfile)
		}

		// Gestión de pacientes
		pacientes := protected.Group("/pacientes")
		{
			pacientes.POST("/", pacienteHandler.CreatePaciente)
			pacientes.GET("/", pacienteHandler.GetAllPacientes)
			pacientes.GET("/search", pacienteHandler.SearchPacientes)
			pacientes.GET("/:id", pacienteHandler.GetPaciente)
			pacientes.PUT("/:id", pacienteHandler.UpdatePaciente)
			pacientes.DELETE("/:id", pacienteHandler.DeletePaciente)
		}

		// Gestión de historial clínico
		historial := protected.Group("/historial")
		{
			historial.POST("/", historialHandler.CreateHistorial)
			historial.GET("/:id", historialHandler.GetHistorial)
			historial.PUT("/:id", historialHandler.UpdateHistorial)
			historial.DELETE("/:id", historialHandler.DeleteHistorial)
			historial.GET("/paciente/:paciente_id", historialHandler.GetHistorialByPaciente)
			historial.GET("/hospital", historialHandler.GetHistorialByHospital)
		}

		// Epidemiología y mapas de calor
		epidemiologia := protected.Group("/epidemiologia")
		{
			epidemiologia.GET("/stats", historialHandler.GetEpidemiologicalStats)
			epidemiologia.GET("/contagious", historialHandler.GetContagiousHistorial)
		}
	}

	return router
}
