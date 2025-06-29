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
	propagacionHandler := handlers.NewPropagacionHandler() // AGREGADO: handler faltante

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

		// CORREGIDO: Sintaxis correcta para el grupo de propagación
		propagacionGroup := api.Group("/propagacion")
		{
			// Análisis principal de velocidad de propagación
			propagacionGroup.GET("/analizar", propagacionHandler.AnalyzeSpreadVelocity)
			
			// Predicciones específicas por distrito
			propagacionGroup.GET("/distrito/:distrito", propagacionHandler.GetDistrictPrediction)
			
			// Comparación entre enfermedades
			propagacionGroup.GET("/comparar", propagacionHandler.GetSpreadComparison)
			
			// Análisis de densidad poblacional
			propagacionGroup.GET("/densidad", propagacionHandler.GetDensityAnalysis)
			
			// Rutas de propagación
			propagacionGroup.GET("/rutas", propagacionHandler.GetSpreadRoutes)
		}

		// Grupo público para datos de referencia (sin autenticación)
		publicGroup := api.Group("/public/propagacion") // CORREGIDO: ruta simplificada
		{
			// Información básica de distritos de Santa Cruz
			publicGroup.GET("/distritos", func(c *gin.Context) {
				distritos := map[string]interface{}{
					"ciudad": "Santa Cruz de la Sierra",
					"distritos": []map[string]interface{}{
						{
							"nombre":       "Equipetrol",
							"habitantes":   85000,
							"area_km2":     12.5,
							"densidad":     6800,
							"tipo":         "Residencial-Comercial",
							"coordenadas":  map[string]float64{"lat": -17.7690416, "lng": -63.1956686},
						},
						{
							"nombre":       "Norte",
							"habitantes":   320000,
							"area_km2":     45.8,
							"densidad":     6986,
							"tipo":         "Residencial-Popular",
							"coordenadas":  map[string]float64{"lat": -17.7987909, "lng": -63.210345},
						},
						{
							"nombre":       "Plan Tres Mil",
							"habitantes":   180000,
							"area_km2":     22.3,
							"densidad":     8072,
							"tipo":         "Popular-Alta Densidad",
							"coordenadas":  map[string]float64{"lat": -17.798792, "lng": -63.210345},
						},
						{
							"nombre":       "Villa 1ro de Mayo",
							"habitantes":   95000,
							"area_km2":     18.7,
							"densidad":     5080,
							"tipo":         "Residencial",
							"coordenadas":  map[string]float64{"lat": -17.7379806, "lng": -63.2484834},
						},
						{
							"nombre":       "Sur",
							"habitantes":   125000,
							"area_km2":     28.4,
							"densidad":     4401,
							"tipo":         "Residencial-Comercial",
							"coordenadas":  map[string]float64{"lat": -17.7441931, "lng": -63.1801563},
						},
						{
							"nombre":       "Oeste",
							"habitantes":   75000,
							"area_km2":     35.2,
							"densidad":     2131,
							"tipo":         "Residencial-Periférico",
							"coordenadas":  map[string]float64{"lat": -17.7439533, "lng": -63.1756103},
						},
						{
							"nombre":       "Este",
							"habitantes":   60000,
							"area_km2":     42.1,
							"densidad":     1425,
							"tipo":         "Periférico-Rural",
							"coordenadas":  map[string]float64{"lat": -17.7728417, "lng": -63.2374135},
						},
						{
							"nombre":       "Centro",
							"habitantes":   45000,
							"area_km2":     8.2,
							"densidad":     5488,
							"tipo":         "Comercial-Histórico",
							"coordenadas":  map[string]float64{"lat": -17.7807346, "lng": -63.1890985},
						},
					},
					"estadisticas": map[string]interface{}{
						"poblacion_total": 1970000,
						"area_total_km2":  187.2,
						"densidad_promedio": 5245,
					},
				}
				
				c.JSON(200, map[string]interface{}{
					"success": true,
					"message": "Información de distritos de Santa Cruz obtenida exitosamente",
					"data":    distritos,
				})
			})
			
			// Matriz de conectividad entre distritos
			publicGroup.GET("/conectividad", func(c *gin.Context) {
				conectividad := map[string]interface{}{
					"matriz_conectividad": map[string][]string{
						"Equipetrol":       {"Norte", "Centro", "Sur"},
						"Norte":            {"Equipetrol", "Plan Tres Mil", "Este"},
						"Plan Tres Mil":    {"Norte", "Sur", "Este"},
						"Villa 1ro de Mayo": {"Oeste", "Centro"},
						"Sur":              {"Equipetrol", "Plan Tres Mil", "Centro"},
						"Oeste":            {"Villa 1ro de Mayo", "Centro"},
						"Este":             {"Norte", "Plan Tres Mil"},
						"Centro":           {"Equipetrol", "Sur", "Oeste", "Villa 1ro de Mayo"},
					},
					"descripcion": "Matriz de conectividad entre distritos de Santa Cruz de la Sierra",
					"criterios": []string{
						"Proximidad geográfica",
						"Conexiones de transporte público",
						"Flujo poblacional diario",
						"Corredores comerciales",
					},
				}
				
				c.JSON(200, map[string]interface{}{
					"success": true,
					"message": "Matriz de conectividad obtenida exitosamente",
					"data":    conectividad,
				})
			})
		} // CORREGIDO: Cierre correcto del publicGroup
	} // CORREGIDO: Cierre correcto del api group

	return router
}