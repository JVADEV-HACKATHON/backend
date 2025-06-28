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
func (h *HistorialHandler) CreateHistorial(c *gin.Context) {
	var historial models.HistorialClinico

	// Obtener ID del hospital del contexto (desde JWT)
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Hospital no autenticado", "NOT_AUTHENTICATED", "")
		return
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&historial); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	// Asignar hospital desde el JWT
	historial.IDHospital = hospitalID.(uint)

	// Validar datos
	if err := h.validator.Struct(historial); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Crear historial
	if err := h.historialService.CreateHistorial(&historial); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al crear historial", "CREATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, historial, "Historial clínico creado exitosamente")
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
