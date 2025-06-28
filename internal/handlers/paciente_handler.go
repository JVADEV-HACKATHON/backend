package handlers

import (
	"net/http"
	"strconv"

	"hospital-api/internal/models"
	"hospital-api/internal/services"
	"hospital-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PacienteHandler struct {
	pacienteService *services.PacienteService
	validator       *validator.Validate
}

// NewPacienteHandler crea una nueva instancia del handler de pacientes
func NewPacienteHandler() *PacienteHandler {
	return &PacienteHandler{
		pacienteService: services.NewPacienteService(),
		validator:       validator.New(),
	}
}

// CreatePaciente crea un nuevo paciente
// @Summary Crear paciente
// @Description Crea un nuevo paciente en el sistema
// @Tags pacientes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param paciente body models.Paciente true "Datos del paciente"
// @Success 201 {object} models.Paciente
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Router /pacientes [post]
func (h *PacienteHandler) CreatePaciente(c *gin.Context) {
	var paciente models.Paciente

	// Bind JSON
	if err := c.ShouldBindJSON(&paciente); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	// Validar datos
	if err := h.validator.Struct(paciente); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Crear paciente
	if err := h.pacienteService.CreatePaciente(&paciente); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al crear paciente", "CREATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, paciente, "Paciente creado exitosamente")
}

// GetPaciente obtiene un paciente por ID
// @Summary Obtener paciente
// @Description Obtiene los datos de un paciente por su ID
// @Tags pacientes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del paciente"
// @Success 200 {object} models.Paciente
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /pacientes/{id} [get]
func (h *PacienteHandler) GetPaciente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", "")
		return
	}

	paciente, err := h.pacienteService.GetPacienteByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "NOT_FOUND", "")
		return
	}

	utils.SuccessResponse(c, paciente, "Paciente obtenido exitosamente")
}

// GetAllPacientes obtiene todos los pacientes con paginación
// @Summary Listar pacientes
// @Description Obtiene una lista paginada de todos los pacientes
// @Tags pacientes
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /pacientes [get]
func (h *PacienteHandler) GetAllPacientes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	pacientes, total, err := h.pacienteService.GetAllPacientes(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener pacientes", "FETCH_ERROR", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, pacientes, "Pacientes obtenidos exitosamente", page, limit, total)
}

// UpdatePaciente actualiza un paciente
// @Summary Actualizar paciente
// @Description Actualiza los datos de un paciente existente
// @Tags pacientes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del paciente"
// @Param paciente body models.Paciente true "Datos actualizados del paciente"
// @Success 200 {object} utils.APISuccessResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /pacientes/{id} [put]
func (h *PacienteHandler) UpdatePaciente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", "")
		return
	}

	var updates models.Paciente
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	// Validar datos (omitir validaciones required para updates parciales)
	// Aquí podrías implementar validaciones específicas para updates

	if err := h.pacienteService.UpdatePaciente(uint(id), &updates); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar paciente", "UPDATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Paciente actualizado exitosamente")
}

// DeletePaciente elimina un paciente
// @Summary Eliminar paciente
// @Description Elimina un paciente del sistema (soft delete)
// @Tags pacientes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del paciente"
// @Success 200 {object} utils.APISuccessResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /pacientes/{id} [delete]
func (h *PacienteHandler) DeletePaciente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", "")
		return
	}

	if err := h.pacienteService.DeletePaciente(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar paciente", "DELETE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Paciente eliminado exitosamente")
}

// SearchPacientes busca pacientes por nombre
// @Summary Buscar pacientes
// @Description Busca pacientes por nombre con paginación
// @Tags pacientes
// @Produce json
// @Security BearerAuth
// @Param q query string true "Término de búsqueda"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /pacientes/search [get]
func (h *PacienteHandler) SearchPacientes(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Término de búsqueda requerido", "MISSING_QUERY", "")
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

	pacientes, total, err := h.pacienteService.SearchPacientes(query, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error en la búsqueda", "SEARCH_ERROR", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, pacientes, "Búsqueda completada exitosamente", page, limit, total)
}
