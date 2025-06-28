package utils

import "math"

// CalcularDistanciaHaversine calcula la distancia entre dos puntos geográficos usando la fórmula de Haversine
func CalcularDistanciaHaversine(lat1, lng1, lat2, lng2 float64) float64 {
	// Radio de la Tierra en km
	const R = 6371.0

	// Convertir a radianes
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	// Diferencias
	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	// Fórmula de Haversine
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
