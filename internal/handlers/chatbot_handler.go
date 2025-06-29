// internal/handlers/chatbot_handler.go
package handlers

import (
	"hospital-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatbotHandler struct {
	chatbotService *services.ChatbotService
}

func NewChatbotHandler() *ChatbotHandler {
	return &ChatbotHandler{
		chatbotService: services.NewChatbotService(),
	}
}

// ChatRequest representa la estructura de la petición del chat
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

// ChatResponse representa la respuesta del chatbot
type ChatResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Chat maneja las conversaciones con el chatbot médico
func (h *ChatbotHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ChatResponse{
			Success: false,
			Error:   "Formato de mensaje inválido: " + err.Error(),
		})
		return
	}

	// Validar que el mensaje no esté vacío
	if len(req.Message) == 0 {
		c.JSON(http.StatusBadRequest, ChatResponse{
			Success: false,
			Error:   "El mensaje no puede estar vacío",
		})
		return
	}

	// Validar longitud del mensaje
	if len(req.Message) > 1000 {
		c.JSON(http.StatusBadRequest, ChatResponse{
			Success: false,
			Error:   "El mensaje es demasiado largo (máximo 1000 caracteres)",
		})
		return
	}

	// Procesar el mensaje a través del service
	response, err := h.chatbotService.ProcessMessage(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ChatResponse{
			Success: false,
			Error:   "Error procesando mensaje: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{
		Success:  true,
		Message:  "Respuesta generada exitosamente",
		Response: response,
	})
}

// HealthCheck verifica el estado del servicio de chatbot
func (h *ChatbotHandler) HealthCheck(c *gin.Context) {
	status, err := h.chatbotService.HealthCheck()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"status":  "unhealthy",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "healthy",
		"data":    status,
	})
}