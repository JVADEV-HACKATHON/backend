package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIErrorResponse estructura estándar para respuestas de error
type APIErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// APISuccessResponse estructura estándar para respuestas exitosas
type APISuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// PaginatedResponse estructura para respuestas paginadas
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination información de paginación
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ErrorResponse envía una respuesta de error estandarizada
func ErrorResponse(c *gin.Context, statusCode int, message, code, details string) {
	c.JSON(statusCode, APIErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
		Details: details,
	})
}

// SuccessResponse envía una respuesta exitosa estandarizada
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APISuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// PaginatedSuccessResponse envía una respuesta exitosa paginada
func PaginatedSuccessResponse(c *gin.Context, data interface{}, message string, page, limit int, total int64) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Message: message,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// ValidationErrorResponse envía una respuesta de error de validación
func ValidationErrorResponse(c *gin.Context, err error) {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, getFieldErrorMessage(fieldError))
		}
	} else {
		errors = append(errors, err.Error())
	}

	ErrorResponse(c, http.StatusBadRequest, "Errores de validación", "VALIDATION_ERROR", strings.Join(errors, "; "))
}

// getFieldErrorMessage retorna un mensaje de error personalizado para cada tipo de validación
func getFieldErrorMessage(fieldError validator.FieldError) string {
	field := fieldError.Field()
	tag := fieldError.Tag()

	switch tag {
	case "required":
		return field + " es requerido"
	case "email":
		return field + " debe ser un email válido"
	case "min":
		return field + " debe tener al menos " + fieldError.Param() + " caracteres"
	case "max":
		return field + " debe tener como máximo " + fieldError.Param() + " caracteres"
	case "gt":
		return field + " debe ser mayor que " + fieldError.Param()
	case "gte":
		return field + " debe ser mayor o igual que " + fieldError.Param()
	case "lt":
		return field + " debe ser menor que " + fieldError.Param()
	case "lte":
		return field + " debe ser menor o igual que " + fieldError.Param()
	case "oneof":
		return field + " debe ser uno de: " + fieldError.Param()
	case "latitude":
		return field + " debe ser una latitud válida"
	case "longitude":
		return field + " debe ser una longitud válida"
	default:
		return field + " es inválido"
	}
}
