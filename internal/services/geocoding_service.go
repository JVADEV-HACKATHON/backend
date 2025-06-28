package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"hospital-api/internal/utils"

	"googlemaps.github.io/maps"
)

type GeocodingService struct {
	client *maps.Client
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type AddressComponents struct {
	FormattedAddress string      `json:"formatted_address"`
	District         string      `json:"district"`
	Neighborhood     string      `json:"neighborhood"`
	City             string      `json:"city"`
	Country          string      `json:"country"`
	Coordinates      Coordinates `json:"coordinates"`
}

// NewGeocodingService crea una nueva instancia del servicio de geocodificación
func NewGeocodingService() (*GeocodingService, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GOOGLE_MAPS_API_KEY no está configurada")
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creando cliente de Google Maps: %v", err)
	}

	return &GeocodingService{
		client: client,
	}, nil
}

// GetCoordinatesFromAddress obtiene las coordenadas de una dirección
func (g *GeocodingService) GetCoordinatesFromAddress(address string) (*Coordinates, error) {
	// Limpiar y formatear la dirección
	cleanAddress := strings.TrimSpace(address)
	if cleanAddress == "" {
		return nil, errors.New("la dirección no puede estar vacía")
	}

	// Agregar "Santa Cruz, Bolivia" si no está incluido para mayor precisión
	if !strings.Contains(strings.ToLower(cleanAddress), "santa cruz") {
		cleanAddress = fmt.Sprintf("%s, Santa Cruz de la Sierra, Bolivia", cleanAddress)
	}

	// Realizar la geocodificación
	req := &maps.GeocodingRequest{
		Address: cleanAddress,
	}

	resp, err := g.client.Geocode(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error en geocodificación: %v", err)
	}

	if len(resp) == 0 {
		return nil, errors.New("no se encontraron resultados para la dirección proporcionada")
	}

	// Tomar el primer resultado (más relevante)
	result := resp[0]

	return &Coordinates{
		Latitude:  result.Geometry.Location.Lat,
		Longitude: result.Geometry.Location.Lng,
	}, nil
}

// GetAddressComponents obtiene información completa de una dirección
func (g *GeocodingService) GetAddressComponents(address string) (*AddressComponents, error) {
	// Limpiar y formatear la dirección
	cleanAddress := strings.TrimSpace(address)
	if cleanAddress == "" {
		return nil, errors.New("la dirección no puede estar vacía")
	}

	// Agregar "Santa Cruz, Bolivia" si no está incluido
	if !strings.Contains(strings.ToLower(cleanAddress), "santa cruz") {
		cleanAddress = fmt.Sprintf("%s, Santa Cruz de la Sierra, Bolivia", cleanAddress)
	}

	// Realizar la geocodificación
	req := &maps.GeocodingRequest{
		Address: cleanAddress,
	}

	resp, err := g.client.Geocode(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error en geocodificación: %v", err)
	}

	if len(resp) == 0 {
		return nil, errors.New("no se encontraron resultados para la dirección proporcionada")
	}

	// Tomar el primer resultado
	result := resp[0]

	components := &AddressComponents{
		FormattedAddress: result.FormattedAddress,
		Coordinates: Coordinates{
			Latitude:  result.Geometry.Location.Lat,
			Longitude: result.Geometry.Location.Lng,
		},
	}

	// Extraer componentes específicos
	for _, component := range result.AddressComponents {
		for _, componentType := range component.Types {
			switch componentType {
			case "sublocality", "sublocality_level_1":
				if components.District == "" {
					components.District = component.LongName
				}
			case "sublocality_level_2", "neighborhood":
				if components.Neighborhood == "" {
					components.Neighborhood = component.LongName
				}
			case "locality", "administrative_area_level_2":
				if components.City == "" {
					components.City = component.LongName
				}
			case "country":
				components.Country = component.LongName
			}
		}
	}

	// Si no se encontró distrito, usar la ciudad
	if components.District == "" {
		components.District = components.City
	}

	return components, nil
}

// ValidateCoordinates valida que las coordenadas estén dentro de los límites de Santa Cruz
func (g *GeocodingService) ValidateCoordinates(lat, lng float64) bool {
	// Límites aproximados de Santa Cruz de la Sierra, Bolivia
	// Latitud: -17.9 a -17.7
	// Longitud: -63.3 a -63.0
	return lat >= -17.9 && lat <= -17.7 && lng >= -63.3 && lng <= -63.0
}

// EvaluarPrecisionGeocoding evalúa qué tan precisa es la ubicación geocodificada
func (g *GeocodingService) EvaluarPrecisionGeocoding(address *AddressComponents) map[string]interface{} {
	result := make(map[string]interface{})

	// Iniciar con una confianza base
	confidence := 0.5

	// 1. Verificar si tiene coordenadas precisas (más de 5 decimales)
	latStr := fmt.Sprintf("%.7f", address.Coordinates.Latitude)
	lngStr := fmt.Sprintf("%.7f", address.Coordinates.Longitude)
	latDecimals := len(latStr) - strings.IndexByte(latStr, '.') - 1
	lngDecimals := len(lngStr) - strings.IndexByte(lngStr, '.') - 1

	// Calcular precisión basada en decimales
	if latDecimals >= 6 && lngDecimals >= 6 {
		confidence += 0.2
		result["precision_nivel"] = "muy alta"
	} else if latDecimals >= 5 && lngDecimals >= 5 {
		confidence += 0.15
		result["precision_nivel"] = "alta"
	} else if latDecimals >= 4 && lngDecimals >= 4 {
		confidence += 0.1
		result["precision_nivel"] = "media"
	} else {
		result["precision_nivel"] = "baja"
	}

	// 2. Verificar componentes de dirección
	if address.District != "" {
		confidence += 0.1
	}
	if address.Neighborhood != "" {
		confidence += 0.15
	}

	// 3. Verificar si hay número en la dirección
	if strings.Count(address.FormattedAddress, " ") > 1 &&
		regexp.MustCompile(`\d+`).MatchString(address.FormattedAddress) {
		confidence += 0.1
	}

	// 4. Verificar que esté dentro de La Paz (esto ya se hace en el servicio)
	if g.ValidateCoordinates(address.Coordinates.Latitude, address.Coordinates.Longitude) {
		confidence += 0.1
	} else {
		confidence -= 0.5 // Penalizar fuertemente si está fuera de La Paz
	}

	// Limitar a rango 0-1
	if confidence > 1.0 {
		confidence = 1.0
	} else if confidence < 0.0 {
		confidence = 0.0
	}

	result["confidence"] = confidence

	// Distancia al centro de Santa Cruz (Plaza 24 de Septiembre: -17.7834, -63.1821)
	centroSantaCruzLat, centroSantaCruzLng := -17.7834, -63.1821
	distancia := utils.CalcularDistanciaHaversine(
		address.Coordinates.Latitude,
		address.Coordinates.Longitude,
		centroSantaCruzLat,
		centroSantaCruzLng)

	result["distancia_centro_ciudad_km"] = distancia

	// Sugerencia de precisión
	if confidence > 0.8 {
		result["sugerencia"] = "Ubicación muy precisa"
	} else if confidence > 0.6 {
		result["sugerencia"] = "Ubicación aceptablemente precisa"
	} else if confidence > 0.4 {
		result["sugerencia"] = "Ubicación con precisión moderada"
	} else {
		result["sugerencia"] = "Ubicación poco precisa, considere verificar manualmente"
	}

	return result
}

// Eliminado y movido a utils.geospatial
