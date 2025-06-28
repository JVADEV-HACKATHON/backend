package handlers

import (
	"net/http"
	"strconv"

	"hospital-api/internal/services"
	"hospital-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HospitalHandler struct {
	hospitalService *services.HospitalService
	validator       *validator.Validate
}

// NewHospitalHandler crea una nueva instancia del handler de hospitales
func NewHospitalHandler() *HospitalHandler {
	return &HospitalHandler{
		hospitalService: services.NewHospitalService(),
		validator:       validator.New(),
	}
}

// GetAllHospitales obtiene todos los hospitales con sus coordenadas
// @Summary Obtener todos los hospitales
// @Description Obtiene una lista de todos los hospitales con sus coordenadas
// @Tags hospitales
// @Produce json
// @Security BearerAuth
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /hospitales [get]
func (h *HospitalHandler) GetAllHospitales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	hospitales, total, err := h.hospitalService.GetAllHospitales(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener hospitales", "FETCH_ERROR", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, hospitales, "Hospitales obtenidos exitosamente", page, limit, total)
}

// GetHospital obtiene un hospital específico por ID
// @Summary Obtener hospital por ID
// @Description Obtiene un hospital específico con sus coordenadas
// @Tags hospitales
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del hospital"
// @Success 200 {object} models.HospitalResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /hospitales/{id} [get]
func (h *HospitalHandler) GetHospital(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", "")
		return
	}

	hospital, err := h.hospitalService.GetHospitalByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "NOT_FOUND", "")
		return
	}

	utils.SuccessResponse(c, hospital.ToResponse(), "Hospital obtenido exitosamente")
}

// GetHospitalesNearby obtiene hospitales cercanos a unas coordenadas
// @Summary Obtener hospitales cercanos
// @Description Obtiene hospitales cercanos a unas coordenadas específicas
// @Tags hospitales
// @Produce json
// @Security BearerAuth
// @Param lat query number true "Latitud"
// @Param lng query number true "Longitud"
// @Param radius query number false "Radio en kilómetros" default(5)
// @Success 200 {array} models.HospitalResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /hospitales/nearby [get]
func (h *HospitalHandler) GetHospitalesNearby(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "5")

	if latStr == "" || lngStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Se requieren parámetros lat y lng", "MISSING_PARAMETERS", "")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Latitud inválida", "INVALID_LATITUDE", "")
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Longitud inválida", "INVALID_LONGITUDE", "")
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Radio inválido", "INVALID_RADIUS", "")
		return
	}

	hospitales, err := h.hospitalService.GetHospitalesNearby(lat, lng, radius)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al buscar hospitales cercanos", "SEARCH_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, hospitales, "Hospitales cercanos obtenidos exitosamente")
}

func (h *HospitalHandler) GetAllHospitalesPublic(c *gin.Context) {
	hospitales, err := h.hospitalService.GetAllHospitalesSinPaginacion()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener hospitales", "FETCH_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, hospitales, "Hospitales obtenidos exitosamente")
}
