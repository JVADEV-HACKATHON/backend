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

// Register maneja el registro de nuevos hospitales
// @Summary Registro de hospital
// @Description Registra un nuevo hospital en el sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param register body services.RegisterRequest true "Datos de registro"
// @Success 201 {object} services.RegisterResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest

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

	// Intentar registro
	response, err := h.authService.Register(req)
	if err != nil {
		// Determinar el código de estado basado en el error
		statusCode := http.StatusBadRequest
		errorCode := "REGISTRATION_FAILED"

		if err.Error() == "el email ya está registrado" || err.Error() == "el teléfono ya está registrado" {
			statusCode = http.StatusConflict
			errorCode = "ALREADY_EXISTS"
		}

		utils.ErrorResponse(c, statusCode, err.Error(), errorCode, "")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": response.Message,
		"data":    response.Hospital,
	})
}
