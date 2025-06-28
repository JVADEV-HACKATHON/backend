package handlers

import (
	"net/http"
	"strconv"
	"time"

	"hospital-api/internal/models"
	"hospital-api/internal/services"
	"hospital-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HistorialHandler struct {
	historialService *services.HistorialService
	validator        *validator.Validate
}

// NewHistorialHandler crea una nueva instancia del handler de historial clínico
func NewHistorialHandler() *HistorialHandler {
	return &HistorialHandler{
		historialService: services.NewHistorialService(),
		validator:        validator.New(),
	}
}

// CreateHistorial crea un nuevo registro de historial clínico
// @Summary Crear historial clínico
// @Description Crea un nuevo registro en el historial clínico de un paciente
// @Tags historial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param historial body models.HistorialClinico true "Datos del historial clínico"
// @Success 201 {object} models.HistorialClinico
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Router /historial [post]
// CreateHistorial crea un nuevo registro de historial clínico con geocodificación
func (h *HistorialHandler) CreateHistorial(c *gin.Context) {
	var request models.HistorialClinicoRequest

	// Obtener ID del hospital del contexto (desde JWT)
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Hospital no autenticado", "NOT_AUTHENTICATED", "")
		return
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	// Validar datos básicos
	if err := h.validator.Struct(request); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Crear servicio de geocodificación
	geocodingService, err := services.NewGeocodingService()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error configurando servicio de mapas", "GEOCODING_CONFIG_ERROR", err.Error())
		return
	}

	// Obtener información completa de la dirección
	addressComponents, err := geocodingService.GetAddressComponents(request.PatientAddress)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No se pudo geocodificar la dirección proporcionada", "GEOCODING_ERROR", err.Error())
		return
	}

	// Validar que las coordenadas estén en La Paz
	if !geocodingService.ValidateCoordinates(addressComponents.Coordinates.Latitude, addressComponents.Coordinates.Longitude) {
		utils.ErrorResponse(c, http.StatusBadRequest, "La dirección debe estar ubicada en La Paz, Bolivia", "INVALID_LOCATION", "")
		return
	}

	// Convertir a modelo de base de datos
	historial := request.ToHistorialClinico()

	// Asignar hospital desde el JWT
	historial.IDHospital = hospitalID.(uint)

	// Asignar coordenadas y información obtenida
	historial.PatientLatitude = addressComponents.Coordinates.Latitude
	historial.PatientLongitude = addressComponents.Coordinates.Longitude
	historial.PatientAddress = addressComponents.FormattedAddress // Usar dirección formateada

	// Si no se proporcionó distrito, usar el obtenido de Google Maps
	if historial.PatientDistrict == "" {
		historial.PatientDistrict = addressComponents.District
	}

	// Si no se proporcionó barrio, usar el obtenido de Google Maps
	if historial.PatientNeighborhood == "" {
		historial.PatientNeighborhood = addressComponents.Neighborhood
	}

	// Crear historial
	if err := h.historialService.CreateHistorial(historial); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al crear historial", "CREATE_ERROR", err.Error())
		return
	}

	// Respuesta con información completa
	response := map[string]interface{}{
		"historial": historial,
		"geocoding_info": map[string]interface{}{
			"original_address":  request.PatientAddress,
			"formatted_address": addressComponents.FormattedAddress,
			"coordinates":       addressComponents.Coordinates,
			"district":          addressComponents.District,
			"neighborhood":      addressComponents.Neighborhood,
		},
	}

	utils.SuccessResponse(c, response, "Historial clínico creado exitosamente con geocodificación")
}

func (h *HistorialHandler) GeocodeAddress(c *gin.Context) {
	var request struct {
		Address string `json:"address" validate:"required,min=5"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	if err := h.validator.Struct(request); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Crear servicio de geocodificación
	geocodingService, err := services.NewGeocodingService()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error configurando servicio de mapas", "GEOCODING_CONFIG_ERROR", err.Error())
		return
	}

	// Geocodificar dirección
	addressComponents, err := geocodingService.GetAddressComponents(request.Address)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No se pudo geocodificar la dirección", "GEOCODING_ERROR", err.Error())
		return
	}

	// Validar ubicación
	isValid := geocodingService.ValidateCoordinates(addressComponents.Coordinates.Latitude, addressComponents.Coordinates.Longitude)

	response := map[string]interface{}{
		"address_components": addressComponents,
		"is_valid_location":  isValid,
		"location_note":      "La dirección debe estar en La Paz, Bolivia",
	}

	utils.SuccessResponse(c, response, "Dirección geocodificada exitosamente")
}

// GetHistorial obtiene un registro de historial clínico por ID
// @Summary Obtener historial clínico
// @Description Obtiene un registro específico del historial clínico con información relacionada
// @Tags historial
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del historial clínico"
// @Success 200 {object} models.HistorialClinico
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /historial/{id} [get]
func (h *HistorialHandler) GetHistorial(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", "")
		return
	}

	historial, err := h.historialService.GetHistorialByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "NOT_FOUND", "")
		return
	}

	utils.SuccessResponse(c, historial, "Historial clínico obtenido exitosamente")
}

// GetHistorialByPaciente obtiene el historial clínico de un paciente específico
// @Summary Obtener historial por paciente
// @Description Obtiene todos los registros del historial clínico de un paciente específico
// @Tags historial
// @Produce json
// @Security BearerAuth
// @Param paciente_id path int true "ID del paciente"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /historial/paciente/{paciente_id} [get]
func (h *HistorialHandler) GetHistorialByPaciente(c *gin.Context) {
	pacienteIDParam := c.Param("paciente_id")
	pacienteID, err := strconv.ParseUint(pacienteIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de paciente inválido", "INVALID_ID", "")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	historiales, total, err := h.historialService.GetHistorialByPaciente(uint(pacienteID), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener historial", "FETCH_ERROR", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, historiales, "Historial del paciente obtenido exitosamente", page, limit, total)
}

// GetHistorialByHospital obtiene el historial clínico del hospital autenticado
// @Summary Obtener historial por hospital
// @Description Obtiene todos los registros del historial clínico del hospital autenticado
// @Tags historial
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /historial/hospital [get]
func (h *HistorialHandler) GetHistorialByHospital(c *gin.Context) {
	// Obtener ID del hospital del contexto (desde JWT)
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Hospital no autenticado", "NOT_AUTHENTICATED", "")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	historiales, total, err := h.historialService.GetHistorialByHospital(hospitalID.(uint), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener historial", "FETCH_ERROR", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, historiales, "Historial del hospital obtenido exitosamente", page, limit, total)
}

// UpdateHistorial actualiza un registro de historial clínico
// @Summary Actualizar historial clínico
// @Description Actualiza un registro existente del historial clínico
// @Tags historial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del historial clínico"
// @Param historial body models.HistorialClinico true "Datos actualizados del historial"
// @Success 200 {object} utils.APISuccessResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /historial/{id} [put]
func (h *HistorialHandler) UpdateHistorial(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", "")
		return
	}

	var updates models.HistorialClinico
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	if err := h.historialService.UpdateHistorial(uint(id), &updates); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar historial", "UPDATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Historial clínico actualizado exitosamente")
}

// DeleteHistorial elimina un registro de historial clínico
// @Summary Eliminar historial clínico
// @Description Elimina un registro del historial clínico (soft delete)
// @Tags historial
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del historial clínico"
// @Success 200 {object} utils.APISuccessResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /historial/{id} [delete]
func (h *HistorialHandler) DeleteHistorial(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", "")
		return
	}

	if err := h.historialService.DeleteHistorial(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar historial", "DELETE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Historial clínico eliminado exitosamente")
}

// GetEpidemiologicalStats obtiene estadísticas epidemiológicas para mapas de calor
// @Summary Estadísticas epidemiológicas
// @Description Obtiene estadísticas epidemiológicas incluyendo datos para mapas de calor
// @Tags epidemiologia
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)" format(date)
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)" format(date)
// @Success 200 {object} services.EpidemiologicalStats
// @Failure 400 {object} utils.APIErrorResponse
// @Router /epidemiologia/stats [get]
func (h *HistorialHandler) GetEpidemiologicalStats(c *gin.Context) {
	// Fechas por defecto: últimos 30 días
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// Parsear fechas si se proporcionan
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	stats, err := h.historialService.GetEpidemiologicalStats(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener estadísticas", "STATS_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, stats, "Estadísticas epidemiológicas obtenidas exitosamente")
}

// GetContagiousHistorial obtiene historiales de casos contagiosos
// @Summary Obtener casos contagiosos
// @Description Obtiene todos los registros marcados como contagiosos
// @Tags epidemiologia
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /epidemiologia/contagious [get]
func (h *HistorialHandler) GetContagiousHistorial(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	historiales, total, err := h.historialService.GetContagiousHistorial(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener casos contagiosos", "FETCH_ERROR", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, historiales, "Casos contagiosos obtenidos exitosamente", page, limit, total)
}

// EvaluateGeocodePrecision evalúa la precisión de una geocodificación
func (h *HistorialHandler) EvaluateGeocodePrecision(c *gin.Context) {
	var request struct {
		Address string `json:"address" validate:"required,min=5"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	if err := h.validator.Struct(request); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Crear servicio de geocodificación
	geocodingService, err := services.NewGeocodingService()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error configurando servicio de mapas", "GEOCODING_CONFIG_ERROR", err.Error())
		return
	}

	// Geocodificar dirección
	addressComponents, err := geocodingService.GetAddressComponents(request.Address)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No se pudo geocodificar la dirección", "GEOCODING_ERROR", err.Error())
		return
	}

	// Validar ubicación
	isValid := geocodingService.ValidateCoordinates(addressComponents.Coordinates.Latitude, addressComponents.Coordinates.Longitude)

	// Evaluar precisión
	evaluacion := geocodingService.EvaluarPrecisionGeocoding(addressComponents)

	response := map[string]interface{}{
		"address_components": addressComponents,
		"is_valid_location":  isValid,
		"location_note":      "La dirección debe estar en La Paz, Bolivia",
		"evaluacion":         evaluacion,
	}

	utils.SuccessResponse(c, response, "Precisión de geocodificación evaluada exitosamente")
}

// GetHistorialByEnfermedad obtiene historiales clínicos por nombre de enfermedad
// @Summary Obtener historiales por enfermedad
// @Description Obtiene todos los registros del historial clínico que coincidan con el nombre de enfermedad especificado
// @Tags historial
// @Produce json
// @Security BearerAuth
// @Param enfermedad query string true "Nombre de la enfermedad"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} models.EnfermedadSearchResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /historial/enfermedad [get]
func (h *HistorialHandler) GetHistorialByEnfermedad(c *gin.Context) {
	enfermedad := c.Query("enfermedad")
	if enfermedad == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro 'enfermedad' es requerido", "MISSING_PARAMETER", "")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	historiales, total, err := h.historialService.GetHistorialByEnfermedad(enfermedad, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener historiales por enfermedad", "FETCH_ERROR", err.Error())
		return
	}

	// Convertir a formato de respuesta específico
	responseData := make([]models.HistorialEnfermedadResponse, len(historiales))
	for i, historial := range historiales {
		responseData[i] = historial.ToEnfermedadResponse()
	}

	// Crear respuesta en el formato solicitado
	response := models.EnfermedadSearchResponse{
		Message: "Datos obtenidos exitosamente",
		Total:   total,
		Data:    responseData,
	}

	c.JSON(http.StatusOK, response)
}
