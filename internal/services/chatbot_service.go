// internal/services/chatbot_service.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ChatbotService struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

type GeminiRequest struct {
	Contents          []GeminiContent          `json:"contents"`
	SystemInstruction *GeminiSystemInstruction `json:"systemInstruction,omitempty"`
	GenerationConfig  *GeminiGenerationConfig  `json:"generationConfig,omitempty"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiSystemInstruction struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiGenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	TopK            int     `json:"topK,omitempty"`
	TopP            float64 `json:"topP,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

const medicalPrompt = `Eres un asistente médico virtual especializado en atención primaria. Tu función es:

RESPONSABILIDADES:
- Proporcionar información médica general basada en evidencia científica
- Realizar triaje básico de síntomas (leve, moderado, urgente)
- Explicar condiciones médicas comunes en términos comprensibles
- Orientar sobre primeros auxilios básicos
- Informar sobre medicamentos de venta libre
- Identificar cuándo es necesaria atención médica profesional

LIMITACIONES IMPORTANTES:
- NO diagnosticas enfermedades
- NO prescribes medicamentos
- NO reemplazas la consulta médica profesional
- NO manejas emergencias médicas (deriva a servicios de emergencia)

ESTILO DE COMUNICACIÓN:
- Empático y comprensivo
- Lenguaje claro y accesible
- Preguntas de seguimiento para clarificar síntomas
- Siempre recomienda consulta profesional para casos complejos

CONTEXTO MÉDICO: Enfócate en medicina general, síntomas comunes, prevención y educación sanitaria.

IMPORTANTE: Si detectas síntomas de emergencia, siempre recomienda acudir inmediatamente a urgencias o llamar al número de emergencias local.`

func NewChatbotService() *ChatbotService {
	return &ChatbotService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:  os.Getenv("GEMINI_API_KEY"),
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-latest:generateContent",
	}
}

func (s *ChatbotService) ProcessMessage(message string) (string, error) {
	// Validar API key
	if s.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY no está configurada")
	}

	// Crear request para Gemini
	geminiReq := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{{Text: message}},
				Role:  "user",
			},
		},
		SystemInstruction: &GeminiSystemInstruction{
			Parts: []GeminiPart{{Text: medicalPrompt}},
		},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     0.7,
			TopK:            40,
			TopP:            0.95,
			MaxOutputTokens: 1024,
		},
	}

	// Llamar a la API de Gemini
	response, err := s.callGeminiAPI(geminiReq)
	if err != nil {
		return "", fmt.Errorf("error llamando a Gemini API: %w", err)
	}

	return response, nil
}

func (s *ChatbotService) callGeminiAPI(req GeminiRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s?key=%s", s.baseURL, s.apiKey)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response content from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

func (s *ChatbotService) HealthCheck() (map[string]interface{}, error) {
	// Verificar API key
	if s.apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY no configurada")
	}

	// Test básico de conectividad
	testReq := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{{Text: "Hello"}},
				Role:  "user",
			},
		},
		GenerationConfig: &GeminiGenerationConfig{
			MaxOutputTokens: 10,
		},
	}

	startTime := time.Now()
	_, err := s.callGeminiAPI(testReq)
	responseTime := time.Since(startTime)

	status := map[string]interface{}{
		"gemini_api":    "connected",
		"response_time": responseTime.Milliseconds(),
		"timestamp":     time.Now(),
	}

	if err != nil {
		status["gemini_api"] = "error"
		status["error"] = err.Error()
		return status, err
	}

	return status, nil
}