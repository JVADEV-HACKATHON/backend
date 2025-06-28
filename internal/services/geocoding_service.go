package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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

	// Agregar "La Paz, Bolivia" si no está incluido para mayor precisión
	if !strings.Contains(strings.ToLower(cleanAddress), "la paz") {
		cleanAddress = fmt.Sprintf("%s, La Paz, Bolivia", cleanAddress)
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

	// Agregar "La Paz, Bolivia" si no está incluido
	if !strings.Contains(strings.ToLower(cleanAddress), "la paz") {
		cleanAddress = fmt.Sprintf("%s, La Paz, Bolivia", cleanAddress)
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

// ValidateCoordinates valida que las coordenadas estén dentro de los límites de La Paz
func (g *GeocodingService) ValidateCoordinates(lat, lng float64) bool {
	// Límites aproximados de La Paz, Bolivia
	// Latitud: -16.7 a -16.4
	// Longitud: -68.2 a -67.9
	return lat >= -16.7 && lat <= -16.4 && lng >= -68.2 && lng <= -67.9
}
