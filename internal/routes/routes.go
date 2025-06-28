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
	hospitalHandler := handlers.NewHospitalHandler()

	// Todas las rutas son públicas ahora
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

		// Autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/profile", authHandler.GetProfile)
		}

		// Gestión de pacientes
		pacientes := api.Group("/pacientes")
		{
			pacientes.POST("/", pacienteHandler.CreatePaciente)
			pacientes.GET("/search", pacienteHandler.SearchPacientes)
			pacientes.GET("/", pacienteHandler.GetAllPacientes)
			pacientes.GET("/:id", pacienteHandler.GetPaciente)
			pacientes.PUT("/:id", pacienteHandler.UpdatePaciente)
			pacientes.DELETE("/:id", pacienteHandler.DeletePaciente)
		}

		// Rutas adicionales de pacientes

		// Gestión de hospitales
		hospitales := api.Group("/hospitales")
		{
			hospitales.GET("/", hospitalHandler.GetAllHospitales)
			hospitales.GET("/nearby", hospitalHandler.GetHospitalesNearby)
			hospitales.GET("/:id", hospitalHandler.GetHospital)
		}

		// Ruta adicional de hospitales

		// Gestión de historial clínico
		historial := api.Group("/historial")
		{
			historial.POST("/", historialHandler.CreateHistorial)
			historial.GET("/", hospitalHandler.GetAllHospitales)
			historial.GET("/:id", historialHandler.GetHistorial)
			historial.PUT("/:id", historialHandler.UpdateHistorial)
			historial.DELETE("/:id", historialHandler.DeleteHistorial)
			historial.GET("/paciente/:paciente_id", historialHandler.GetHistorialByPaciente)
			historial.GET("/enfermedad", historialHandler.GetHistorialByEnfermedad)
		}

		// Endpoints para geocodificación
		api.POST("/geocode", historialHandler.GeocodeAddress)
		api.POST("/geocode/evaluate", historialHandler.EvaluateGeocodePrecision)

		// Epidemiología y mapas de calor
		epidemiologia := api.Group("/epidemiologia")
		{
			epidemiologia.GET("/stats", historialHandler.GetEpidemiologicalStats)
			epidemiologia.GET("/contagious", historialHandler.GetContagiousHistorial)
		}
	}

	return router
}
