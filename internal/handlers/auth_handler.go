package handlers

import (
	"net/http"

	"hospital-api/internal/services"
	"hospital-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
}

// NewAuthHandler crea una nueva instancia del handler de autenticación
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
		validator:   validator.New(),
	}
}

// Login maneja el login de hospitales
// @Summary Login de hospital
// @Description Autentica un hospital y retorna un token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param login body services.LoginRequest true "Credenciales de login"
// @Success 200 {object} services.LoginResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos", "INVALID_INPUT", err.Error())
		return
	}

	// Validar datos
	if err := h.validator.Struct(req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Intentar login
	response, err := h.authService.Login(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "AUTH_FAILED", "")
		return
	}

	utils.SuccessResponse(c, response, "Login exitoso")
}

// GetProfile obtiene el perfil del hospital autenticado
// @Summary Perfil del hospital
// @Description Obtiene la información del hospital autenticado
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.HospitalResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Hospital no autenticado", "NOT_AUTHENTICATED", "")
		return
	}

	// Convertir a uint
	id, ok := hospitalID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error interno", "INTERNAL_ERROR", "Invalid hospital ID type")
		return
	}

	// Aquí podrías obtener más información del hospital desde la base de datos
	// Por ahora, retornamos la información básica disponible en el contexto
	email, _ := c.Get("hospital_email")

	utils.SuccessResponse(c, gin.H{
		"hospital_id": id,
		"email":       email,
	}, "Perfil obtenido exitosamente")
}
