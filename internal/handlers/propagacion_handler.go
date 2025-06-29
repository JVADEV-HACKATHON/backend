package handlers

import (
	"net/http"
	"strconv"
		"fmt"
	"math"
	"strings"
	"time"

	"hospital-api/internal/services"
	"hospital-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


type PropagacionHandler struct {
	propagacionService *services.PropagacionService
	validator          *validator.Validate
}

// NewPropagacionHandler crea una nueva instancia del handler de propagación
func NewPropagacionHandler() *PropagacionHandler {
	return &PropagacionHandler{
		propagacionService: services.NewPropagacionService(),
		validator:          validator.New(),
	}
}

// AnalyzeSpreadVelocity analiza la velocidad de propagación de enfermedades
// @Summary Analizar velocidad de propagación
// @Description Analiza la velocidad de propagación de una enfermedad específica basada en densidad poblacional y patrones históricos
// @Tags propagacion
// @Produce json
// @Security BearerAuth
// @Param enfermedad query string true "Nombre de la enfermedad"
// @Param dias query int false "Días de análisis histórico" default(30)
// @Success 200 {object} services.VelocidadPropagacion
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 500 {object} utils.APIErrorResponse
// @Router /propagacion/analizar [get]
func (h *PropagacionHandler) AnalyzeSpreadVelocity(c *gin.Context) {
	enfermedad := c.Query("enfermedad")
	if enfermedad == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'enfermedad' es requerido", "MISSING_PARAMETER", "")
		return
	}

	diasStr := c.DefaultQuery("dias", "30")
	dias, err := strconv.Atoi(diasStr)
	if err != nil || dias < 7 || dias > 365 {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'dias' debe ser un número entre 7 y 365", "INVALID_PARAMETER", "")
		return
	}

	analisis, err := h.propagacionService.AnalyzeSpreadVelocity(enfermedad, dias)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al analizar velocidad de propagación", "ANALYSIS_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, analisis, "Análisis de velocidad de propagación completado exitosamente")
}

// GetDistrictPrediction obtiene predicciones específicas para un distrito
// @Summary Predicción por distrito
// @Description Obtiene predicciones de propagación específicas para un distrito de Santa Cruz
// @Tags propagacion
// @Produce json
// @Security BearerAuth
// @Param distrito path string true "Nombre del distrito"
// @Param enfermedad query string true "Nombre de la enfermedad"
// @Param dias query int false "Días de análisis histórico" default(30)
// @Success 200 {object} services.PrediccionPropagacion
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /propagacion/distrito/{distrito} [get]
func (h *PropagacionHandler) GetDistrictPrediction(c *gin.Context) {
	distrito := c.Param("distrito")
	if distrito == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'distrito' es requerido", "MISSING_PARAMETER", "")
		return
	}

	enfermedad := c.Query("enfermedad")
	if enfermedad == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'enfermedad' es requerido", "MISSING_PARAMETER", "")
		return
	}

	diasStr := c.DefaultQuery("dias", "30")
	dias, err := strconv.Atoi(diasStr)
	if err != nil || dias < 7 || dias > 365 {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'dias' debe ser un número entre 7 y 365", "INVALID_PARAMETER", "")
		return
	}

	prediccion, err := h.propagacionService.GetSpreadPredictionsByDistrict(distrito, enfermedad, dias)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "No se encontraron predicciones para el distrito especificado", "NOT_FOUND", err.Error())
		return
	}

	utils.SuccessResponse(c, prediccion, "Predicción de propagación obtenida exitosamente")
}

// GetSpreadComparison compara velocidades de propagación entre diferentes enfermedades
// @Summary Comparar propagación entre enfermedades
// @Description Compara las velocidades de propagación entre diferentes enfermedades en Santa Cruz
// @Tags propagacion
// @Produce json
// @Security BearerAuth
// @Param enfermedades query string true "Enfermedades separadas por coma (ej: dengue,zika,influenza)"
// @Param dias query int false "Días de análisis histórico" default(30)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.APIErrorResponse
// @Router /propagacion/comparar [get]
func (h *PropagacionHandler) GetSpreadComparison(c *gin.Context) {
	enfermedadesStr := c.Query("enfermedades")
	if enfermedadesStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'enfermedades' es requerido (separadas por coma)", "MISSING_PARAMETER", "")
		return
	}

	diasStr := c.DefaultQuery("dias", "30")
	dias, err := strconv.Atoi(diasStr)
	if err != nil || dias < 7 || dias > 365 {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'dias' debe ser un número entre 7 y 365", "INVALID_PARAMETER", "")
		return
	}

	// Separar enfermedades por coma
	enfermedades := strings.Split(enfermedadesStr, ",")
	if len(enfermedades) < 2 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Se requieren al menos 2 enfermedades para comparar", "INSUFFICIENT_DATA", "")
		return
	}

	comparacion := make(map[string]interface{})
	var analisisCompletos []services.VelocidadPropagacion

	// Analizar cada enfermedad
	for _, enfermedad := range enfermedades {
		enfermedad = strings.TrimSpace(enfermedad)
		if enfermedad != "" {
			analisis, err := h.propagacionService.AnalyzeSpreadVelocity(enfermedad, dias)
			if err != nil {
				comparacion[enfermedad] = map[string]string{
					"error": err.Error(),
				}
			} else {
				comparacion[enfermedad] = analisis
				analisisCompletos = append(analisisCompletos, *analisis)
			}
		}
	}

	// Generar resumen comparativo
	resumen := h.generarResumenComparativo(analisisCompletos)
	comparacion["resumen_comparativo"] = resumen

	utils.SuccessResponse(c, comparacion, "Comparación de velocidades de propagación completada")
}

// GetDensityAnalysis obtiene análisis detallado de densidad poblacional
// @Summary Análisis de densidad poblacional
// @Description Obtiene análisis detallado de la densidad poblacional por distrito en Santa Cruz
// @Tags propagacion
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /propagacion/densidad [get]
func (h *PropagacionHandler) GetDensityAnalysis(c *gin.Context) {
	// Obtener datos de densidad poblacional
	densityData := map[string]interface{}{
		"ciudad": "Santa Cruz de la Sierra",
		"fecha_analisis": time.Now().Format("2006-01-02"),
		"distritos": h.obtenerDatosDensidad(),
		"estadisticas_generales": h.calcularEstadisticasGenerales(),
		"recomendaciones_vigilancia": h.generarRecomendacionesVigilancia(),
	}

	utils.SuccessResponse(c, densityData, "Análisis de densidad poblacional obtenido exitosamente")
}

// GetSpreadRoutes obtiene las rutas de propagación más probables
// @Summary Rutas de propagación
// @Description Obtiene las rutas de propagación más probables entre distritos para una enfermedad
// @Tags propagacion
// @Produce json
// @Security BearerAuth
// @Param enfermedad query string true "Nombre de la enfermedad"
// @Param origen query string false "Distrito de origen (opcional)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.APIErrorResponse
// @Router /propagacion/rutas [get]
func (h *PropagacionHandler) GetSpreadRoutes(c *gin.Context) {
	enfermedad := c.Query("enfermedad")
	if enfermedad == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'enfermedad' es requerido", "MISSING_PARAMETER", "")
		return
	}

	origen := c.Query("origen")
	dias := 30

	analisis, err := h.propagacionService.AnalyzeSpreadVelocity(enfermedad, dias)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al analizar rutas de propagación", "ANALYSIS_ERROR", err.Error())
		return
	}

	// Filtrar rutas por origen si se especifica
	rutas := analisis.RutasPropagacion
	if origen != "" {
		var rutasFiltradas []services.RutaPropagacion
		for _, ruta := range rutas {
			if strings.EqualFold(ruta.DistritoOrigen, origen) {
				rutasFiltradas = append(rutasFiltradas, ruta)
			}
		}
		rutas = rutasFiltradas
	}

	// Organizar respuesta
	response := map[string]interface{}{
		"enfermedad": enfermedad,
		"origen_filtro": origen,
		"total_rutas": len(rutas),
		"rutas_propagacion": rutas,
		"matriz_conectividad": h.generarMatrizConectividad(),
		"recomendaciones": h.generarRecomendacionesRutas(rutas),
	}

	utils.SuccessResponse(c, response, "Rutas de propagación obtenidas exitosamente")
}

// Métodos auxiliares

func (h *PropagacionHandler) generarResumenComparativo(analisis []services.VelocidadPropagacion) map[string]interface{} {
	if len(analisis) == 0 {
		return map[string]interface{}{"error": "No hay datos suficientes para comparar"}
	}

	var velocidadPromedio, velocidadMaxima float64
	var distritosAfectados []string
	enfermedadMasRapida := ""
	velocidadMasAlta := 0.0

	for _, a := range analisis {
		velocidadPromedio += a.VelocidadPromedio
		if a.VelocidadMaxima > velocidadMaxima {
			velocidadMaxima = a.VelocidadMaxima
		}

		if a.VelocidadPromedio > velocidadMasAlta {
			velocidadMasAlta = a.VelocidadPromedio
			enfermedadMasRapida = a.Enfermedad
		}

		for _, distrito := range a.DistritosAfectados {
			// Agregar distrito si no está ya en la lista
			encontrado := false
			for _, d := range distritosAfectados {
				if d == distrito.Distrito {
					encontrado = true
					break
				}
			}
			if !encontrado {
				distritosAfectados = append(distritosAfectados, distrito.Distrito)
			}
		}
	}

	velocidadPromedio /= float64(len(analisis))

	return map[string]interface{}{
		"total_enfermedades_analizadas": len(analisis),
		"velocidad_promedio_general":    math.Round(velocidadPromedio*100) / 100,
		"velocidad_maxima_registrada":   velocidadMaxima,
		"enfermedad_propagacion_rapida": enfermedadMasRapida,
		"total_distritos_afectados":     len(distritosAfectados),
		"distritos_con_casos":           distritosAfectados,
		"nivel_alerta_general":          h.determinarNivelAlerta(velocidadPromedio),
	}
}

func (h *PropagacionHandler) obtenerDatosDensidad() []map[string]interface{} {
	// Esta información vendría del servicio, aquí la simulamos
	distritos := []map[string]interface{}{
		{
			"distrito":      "Plan Tres Mil",
			"habitantes":    180000,
			"area_km2":      22.3,
			"densidad":      8072,
			"tipo_zona":     "Popular-Alta Densidad",
			"riesgo_base":   "ALTO",
			"conectividad":  []string{"Norte", "Sur", "Este"},
		},
		{
			"distrito":      "Norte",
			"habitantes":    320000,
			"area_km2":      45.8,
			"densidad":      6986,
			"tipo_zona":     "Residencial-Popular",
			"riesgo_base":   "ALTO",
			"conectividad":  []string{"Equipetrol", "Plan Tres Mil", "Este"},
		},
		{
			"distrito":      "Equipetrol",
			"habitantes":    85000,
			"area_km2":      12.5,
			"densidad":      6800,
			"tipo_zona":     "Residencial-Comercial",
			"riesgo_base":   "MEDIO",
			"conectividad":  []string{"Norte", "Centro", "Sur"},
		},
	}

	return distritos
}

func (h *PropagacionHandler) calcularEstadisticasGenerales() map[string]interface{} {
	return map[string]interface{}{
		"poblacion_total_santa_cruz": 1970000,
		"area_total_km2":            187.2,
		"densidad_promedio":         5245,
		"distrito_mayor_densidad":   "Plan Tres Mil",
		"distrito_menor_densidad":   "Este",
		"total_distritos":           8,
	}
}

func (h *PropagacionHandler) generarRecomendacionesVigilancia() []string {
	return []string{
		"🏙️ Priorizar vigilancia epidemiológica en Plan Tres Mil y Norte por alta densidad poblacional",
		"🚌 Monitorear estaciones de transporte público como puntos de dispersión",
		"🏥 Distribuir recursos médicos proporcionalmente a la densidad poblacional",
		"📊 Implementar sistema de alerta temprana en distritos de alta conectividad",
		"🎯 Establecer centros de testeo móviles en zonas de alta densidad",
	}
}

func (h *PropagacionHandler) generarMatrizConectividad() map[string][]string {
	return map[string][]string{
		"Equipetrol":       {"Norte", "Centro", "Sur"},
		"Norte":            {"Equipetrol", "Plan Tres Mil", "Este"},
		"Plan Tres Mil":    {"Norte", "Sur", "Este"},
		"Villa 1ro de Mayo": {"Oeste", "Centro"},
		"Sur":              {"Equipetrol", "Plan Tres Mil", "Centro"},
		"Oeste":            {"Villa 1ro de Mayo", "Centro"},
		"Este":             {"Norte", "Plan Tres Mil"},
		"Centro":           {"Equipetrol", "Sur", "Oeste", "Villa 1ro de Mayo"},
	}
}

func (h *PropagacionHandler) generarRecomendacionesRutas(rutas []services.RutaPropagacion) []string {
	if len(rutas) == 0 {
		return []string{"No se detectaron rutas de propagación activas"}
	}

	recomendaciones := []string{
		"🛣️ Monitorear corredores de alta movilidad entre distritos conectados",
		"📍 Establecer puntos de control epidemiológico en rutas identificadas",
	}

	// Rutas rápidas (menos de 3 días)
	for _, ruta := range rutas {
		if ruta.DiasTransicion <= 3 {
			recomendaciones = append(recomendaciones, 
				fmt.Sprintf("⚡ Alerta: Propagación rápida detectada %s → %s (%d días)", 
					ruta.DistritoOrigen, ruta.DistritoDestino, ruta.DiasTransicion))
		}
	}

	return recomendaciones
}

func (h *PropagacionHandler) determinarNivelAlerta(velocidadPromedio float64) string {
	switch {
	case velocidadPromedio >= 10:
		return "CRÍTICO"
	case velocidadPromedio >= 5:
		return "ALTO"
	case velocidadPromedio >= 2:
		return "MEDIO"
	default:
		return "BAJO"
	}
}



