package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"hospital-api/internal/database"
	"hospital-api/internal/models"

	"gorm.io/gorm"
)

type PropagacionService struct {
	db *gorm.DB
}

// Datos demográficos de Santa Cruz por distrito (habitantes por km²)
var densidadPoblacionalSantaCruz = map[string]DensidadDistrito{
	"Equipetrol": {
		Habitantes:  85000,
		AreaKm2:     12.5,
		Densidad:    6800, // hab/km²
		Conectividad: []string{"Norte", "Centro", "Sur"},
		TipoZona:    "Residencial-Comercial",
	},
	"Norte": {
		Habitantes:  320000,
		AreaKm2:     45.8,
		Densidad:    6986,
		Conectividad: []string{"Equipetrol", "Plan Tres Mil", "Este"},
		TipoZona:    "Residencial-Popular",
	},
	"Plan Tres Mil": {
		Habitantes:  180000,
		AreaKm2:     22.3,
		Densidad:    8072,
		Conectividad: []string{"Norte", "Sur", "Este"},
		TipoZona:    "Popular-Alta Densidad",
	},
	"Villa 1ro de Mayo": {
		Habitantes:  95000,
		AreaKm2:     18.7,
		Densidad:    5080,
		Conectividad: []string{"Oeste", "Centro"},
		TipoZona:    "Residencial",
	},
	"Sur": {
		Habitantes:  125000,
		AreaKm2:     28.4,
		Densidad:    4401,
		Conectividad: []string{"Equipetrol", "Plan Tres Mil", "Centro"},
		TipoZona:    "Residencial-Comercial",
	},
	"Oeste": {
		Habitantes:  75000,
		AreaKm2:     35.2,
		Densidad:    2131,
		Conectividad: []string{"Villa 1ro de Mayo", "Centro"},
		TipoZona:    "Residencial-Periférico",
	},
	"Este": {
		Habitantes:  60000,
		AreaKm2:     42.1,
		Densidad:    1425,
		Conectividad: []string{"Norte", "Plan Tres Mil"},
		TipoZona:    "Periférico-Rural",
	},
	"Centro": {
		Habitantes:  45000,
		AreaKm2:     8.2,
		Densidad:    5488,
		Conectividad: []string{"Equipetrol", "Sur", "Oeste", "Villa 1ro de Mayo"},
		TipoZona:    "Comercial-Histórico",
	},
}

type DensidadDistrito struct {
	Habitantes   int      `json:"habitantes"`
	AreaKm2      float64  `json:"area_km2"`
	Densidad     int      `json:"densidad_hab_km2"`
	Conectividad []string `json:"conectividad"`
	TipoZona     string   `json:"tipo_zona"`
}

type CasoTemporal struct {
	Fecha           time.Time `json:"fecha"`
	Distrito        string    `json:"distrito"`
	TotalCasos      int       `json:"total_casos"`
	CasosContagiosos int      `json:"casos_contagiosos"`
	Coordenadas     Coordenada `json:"coordenadas"`
}

type Coordenada struct {
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
}

type VelocidadPropagacion struct {
	Enfermedad           string                    `json:"enfermedad"`
	PeriodoAnalisis      PeriodoAnalisis          `json:"periodo_analisis"`
	VelocidadPromedio    float64                  `json:"velocidad_promedio_casos_por_dia"`
	VelocidadMaxima      float64                  `json:"velocidad_maxima_casos_por_dia"`
	DistritosAfectados   []DistritoAfectado       `json:"distritos_afectados"`
	RutasPropagacion     []RutaPropagacion        `json:"rutas_propagacion"`
	FactorDensidad       float64                  `json:"factor_densidad"`
	PredictedSpread      []PrediccionPropagacion  `json:"prediccion_propagacion"`
	RecomendacionesAlert []string                 `json:"recomendaciones_alerta"`
}

type PeriodoAnalisis struct {
	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
	DiasTotales int       `json:"dias_totales"`
}

type DistritoAfectado struct {
	Distrito         string  `json:"distrito"`
	PrimerCaso       time.Time `json:"primer_caso"`
	UltimoCaso       time.Time `json:"ultimo_caso"`
	TotalCasos       int     `json:"total_casos"`
	DensidadHab      int     `json:"densidad_habitantes"`
	VelocidadLocal   float64 `json:"velocidad_local_casos_por_dia"`
	RiesgoExpansion  string  `json:"riesgo_expansion"`
}

type RutaPropagacion struct {
	DistritoOrigen  string    `json:"distrito_origen"`
	DistritoDestino string    `json:"distrito_destino"`
	FechaPropagacion time.Time `json:"fecha_propagacion"`
	DiasTransicion  int       `json:"dias_transicion"`
	DistanciaKm     float64   `json:"distancia_km"`
	VelocidadKmDia  float64   `json:"velocidad_km_por_dia"`
}

type PrediccionPropagacion struct {
	Distrito        string    `json:"distrito"`
	FechaPrediccion time.Time `json:"fecha_prediccion"`
	CasosPredichos  int       `json:"casos_predichos"`
	Probabilidad    float64   `json:"probabilidad"`
	NivelRiesgo     string    `json:"nivel_riesgo"`
}

func NewPropagacionService() *PropagacionService {
	return &PropagacionService{
		db: database.GetDB(),
	}
}

// AnalyzeSpreadVelocity analiza la velocidad de propagación de una enfermedad específica
func (s *PropagacionService) AnalyzeSpreadVelocity(enfermedad string, diasAnalisis int) (*VelocidadPropagacion, error) {
	// Calcular período de análisis
	fechaFin := time.Now()
	fechaInicio := fechaFin.AddDate(0, 0, -diasAnalisis)

	// Obtener casos temporales
	casosTemporales, err := s.obtenerCasosTemporales(enfermedad, fechaInicio, fechaFin)
	if err != nil {
		return nil, err
	}

	if len(casosTemporales) == 0 {
		return nil, fmt.Errorf("no se encontraron casos para la enfermedad %s en el período especificado", enfermedad)
	}

	// Analizar distritos afectados
	distritosAfectados := s.analizarDistritosAfectados(casosTemporales)

	// Calcular rutas de propagación
	rutasPropagacion := s.calcularRutasPropagacion(distritosAfectados)

	// Calcular velocidades
	velocidadPromedio, velocidadMaxima := s.calcularVelocidades(casosTemporales, diasAnalisis)

	// Calcular factor de densidad
	factorDensidad := s.calcularFactorDensidad(distritosAfectados)

	// Generar predicciones
	predicciones := s.generarPredicciones(distritosAfectados, velocidadPromedio, factorDensidad)

	// Generar recomendaciones
	recomendaciones := s.generarRecomendaciones(distritosAfectados, velocidadPromedio, factorDensidad)

	resultado := &VelocidadPropagacion{
		Enfermedad: enfermedad,
		PeriodoAnalisis: PeriodoAnalisis{
			FechaInicio: fechaInicio,
			FechaFin:    fechaFin,
			DiasTotales: diasAnalisis,
		},
		VelocidadPromedio:    velocidadPromedio,
		VelocidadMaxima:      velocidadMaxima,
		DistritosAfectados:   distritosAfectados,
		RutasPropagacion:     rutasPropagacion,
		FactorDensidad:       factorDensidad,
		PredictedSpread:      predicciones,
		RecomendacionesAlert: recomendaciones,
	}

	return resultado, nil
}

func (s *PropagacionService) obtenerCasosTemporales(enfermedad string, fechaInicio, fechaFin time.Time) ([]CasoTemporal, error) {
	var resultados []struct {
		Fecha    time.Time `json:"fecha"`
		Distrito string    `json:"distrito"`
		Count    int       `json:"count"`
		ContagiousCount int  `json:"contagious_count"`
		AvgLat   float64   `json:"avg_lat"`
		AvgLng   float64   `json:"avg_lng"`
	}

	err := s.db.Model(&models.HistorialClinico{}).
		Select(`
			consultation_date::date as fecha,
			patient_district as distrito,
			COUNT(*) as count,
			COUNT(CASE WHEN is_contagious = true THEN 1 END) as contagious_count,
			AVG(patient_latitude) as avg_lat,
			AVG(patient_longitude) as avg_lng
		`).
		Where("LOWER(enfermedad) = LOWER(?) AND consultation_date BETWEEN ? AND ?", enfermedad, fechaInicio, fechaFin).
		Group("consultation_date::date, patient_district").
		Order("fecha ASC, distrito").
		Scan(&resultados).Error

	if err != nil {
		return nil, err
	}

	casosTemporales := make([]CasoTemporal, len(resultados))
	for i, resultado := range resultados {
		casosTemporales[i] = CasoTemporal{
			Fecha:           resultado.Fecha,
			Distrito:        resultado.Distrito,
			TotalCasos:      resultado.Count,
			CasosContagiosos: resultado.ContagiousCount,
			Coordenadas: Coordenada{
				Latitud:  resultado.AvgLat,
				Longitud: resultado.AvgLng,
			},
		}
	}

	return casosTemporales, nil
}

func (s *PropagacionService) analizarDistritosAfectados(casos []CasoTemporal) []DistritoAfectado {
	distritoMap := make(map[string]*DistritoAfectado)

	// Agrupar casos por distrito
	for _, caso := range casos {
		if distrito, exists := distritoMap[caso.Distrito]; exists {
			distrito.TotalCasos += caso.TotalCasos
			if caso.Fecha.After(distrito.UltimoCaso) {
				distrito.UltimoCaso = caso.Fecha
			}
			if caso.Fecha.Before(distrito.PrimerCaso) {
				distrito.PrimerCaso = caso.Fecha
			}
		} else {
			densidad := 0
			if info, exists := densidadPoblacionalSantaCruz[caso.Distrito]; exists {
				densidad = info.Densidad
			}

			distritoMap[caso.Distrito] = &DistritoAfectado{
				Distrito:        caso.Distrito,
				PrimerCaso:      caso.Fecha,
				UltimoCaso:      caso.Fecha,
				TotalCasos:      caso.TotalCasos,
				DensidadHab:     densidad,
				VelocidadLocal:  0,
				RiesgoExpansion: "",
			}
		}
	}

	// Convertir a slice y calcular velocidades locales
	var distritos []DistritoAfectado
	for _, distrito := range distritoMap {
		// Calcular velocidad local
		dias := int(distrito.UltimoCaso.Sub(distrito.PrimerCaso).Hours()/24) + 1
		if dias > 0 {
			distrito.VelocidadLocal = float64(distrito.TotalCasos) / float64(dias)
		}

		// Calcular riesgo de expansión
		distrito.RiesgoExpansion = s.calcularRiesgoExpansion(distrito.DensidadHab, distrito.VelocidadLocal)

		distritos = append(distritos, *distrito)
	}

	// Ordenar por total de casos descendente
	sort.Slice(distritos, func(i, j int) bool {
		return distritos[i].TotalCasos > distritos[j].TotalCasos
	})

	return distritos
}

func (s *PropagacionService) calcularRutasPropagacion(distritos []DistritoAfectado) []RutaPropagacion {
	var rutas []RutaPropagacion

	// Ordenar distritos por fecha del primer caso
	sort.Slice(distritos, func(i, j int) bool {
		return distritos[i].PrimerCaso.Before(distritos[j].PrimerCaso)
	})

	// Analizar propagación entre distritos conectados
	for i, origen := range distritos {
		if conectividad, exists := densidadPoblacionalSantaCruz[origen.Distrito]; exists {
			for _, distritoConectado := range conectividad.Conectividad {
				// Buscar el distrito conectado en la lista de afectados
				for j, destino := range distritos {
					if j > i && destino.Distrito == distritoConectado {
						diasTransicion := int(destino.PrimerCaso.Sub(origen.PrimerCaso).Hours() / 24)
						if diasTransicion > 0 && diasTransicion <= 14 { // Máximo 14 días para considerar propagación directa
							distancia := s.calcularDistanciaKm(origen.Distrito, destino.Distrito)
							velocidadKm := 0.0
							if diasTransicion > 0 {
								velocidadKm = distancia / float64(diasTransicion)
							}

							ruta := RutaPropagacion{
								DistritoOrigen:   origen.Distrito,
								DistritoDestino:  destino.Distrito,
								FechaPropagacion: destino.PrimerCaso,
								DiasTransicion:   diasTransicion,
								DistanciaKm:      distancia,
								VelocidadKmDia:   velocidadKm,
							}
							rutas = append(rutas, ruta)
						}
					}
				}
			}
		}
	}

	return rutas
}

func (s *PropagacionService) calcularDistanciaKm(distrito1, distrito2 string) float64 {
	// Coordenadas aproximadas del centro de cada distrito en Santa Cruz
	coordenadas := map[string]Coordenada{
		"Equipetrol":       {-17.7690416, -63.1956686},
		"Norte":            {-17.7987909, -63.210345},
		"Plan Tres Mil":    {-17.798792, -63.210345},
		"Villa 1ro de Mayo": {-17.7379806, -63.2484834},
		"Sur":              {-17.7441931, -63.1801563},
		"Oeste":            {-17.7439533, -63.1756103},
		"Este":             {-17.7728417, -63.2374135},
		"Centro":           {-17.7807346, -63.1890985},
	}

	coord1, exists1 := coordenadas[distrito1]
	coord2, exists2 := coordenadas[distrito2]

	if !exists1 || !exists2 {
		return 0
	}

	return s.calcularDistanciaHaversine(coord1.Latitud, coord1.Longitud, coord2.Latitud, coord2.Longitud)
}

func (s *PropagacionService) calcularDistanciaHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radio de la Tierra en km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := R * c

	return distance
}

func (s *PropagacionService) calcularVelocidades(casos []CasoTemporal, diasAnalisis int) (promedio, maxima float64) {
	if len(casos) == 0 {
		return 0, 0
	}

	// Agrupar casos por día
	casosPorDia := make(map[string]int)
	for _, caso := range casos {
		fecha := caso.Fecha.Format("2006-01-02")
		casosPorDia[fecha] += caso.TotalCasos
	}

	// Calcular velocidades
	totalCasos := 0
	maxCasosDia := 0

	for _, casosDia := range casosPorDia {
		totalCasos += casosDia
		if casosDia > maxCasosDia {
			maxCasosDia = casosDia
		}
	}

	promedio = float64(totalCasos) / float64(diasAnalisis)
	maxima = float64(maxCasosDia)

	return promedio, maxima
}

func (s *PropagacionService) calcularFactorDensidad(distritos []DistritoAfectado) float64 {
	if len(distritos) == 0 {
		return 0
	}

	sumaFactores := 0.0
	for _, distrito := range distritos {
		// Factor basado en densidad poblacional y casos
		factor := float64(distrito.DensidadHab) * distrito.VelocidadLocal / 1000
		sumaFactores += factor
	}

	return sumaFactores / float64(len(distritos))
}

func (s *PropagacionService) calcularRiesgoExpansion(densidad int, velocidad float64) string {
	score := float64(densidad)/1000 + velocidad*2

	switch {
	case score >= 15:
		return "CRÍTICO"
	case score >= 10:
		return "ALTO"
	case score >= 5:
		return "MEDIO"
	default:
		return "BAJO"
	}
}

func (s *PropagacionService) generarPredicciones(distritos []DistritoAfectado, velocidadPromedio, factorDensidad float64) []PrediccionPropagacion {
	var predicciones []PrediccionPropagacion

	// Obtener distritos no afectados o con baja incidencia
	for distritoNombre, info := range densidadPoblacionalSantaCruz {
		afectado := false
		casosActuales := 0

		for _, distrito := range distritos {
			if distrito.Distrito == distritoNombre {
				afectado = true
				casosActuales = distrito.TotalCasos
				break
			}
		}

		// Predecir para distritos no afectados o con pocos casos
		if !afectado || casosActuales < 5 {
			prediccion := s.calcularPrediccion(distritoNombre, info, velocidadPromedio, factorDensidad, casosActuales)
			predicciones = append(predicciones, prediccion)
		}
	}

	// Ordenar por probabilidad descendente
	sort.Slice(predicciones, func(i, j int) bool {
		return predicciones[i].Probabilidad > predicciones[j].Probabilidad
	})

	return predicciones
}

func (s *PropagacionService) calcularPrediccion(distrito string, info DensidadDistrito, velocidadPromedio, factorDensidad float64, casosActuales int) PrediccionPropagacion {
	// Calcular probabilidad basada en densidad, conectividad y velocidad de propagación
	probabilidadBase := float64(info.Densidad) / 10000 // Normalizar densidad
	factorConectividad := float64(len(info.Conectividad)) / 10
	factorVelocidad := velocidadPromedio / 10

	probabilidad := (probabilidadBase + factorConectividad + factorVelocidad + factorDensidad/10) * 100
	if probabilidad > 100 {
		probabilidad = 100
	}

	// Calcular casos predichos
	casosPredichos := int(velocidadPromedio * probabilidad / 100)
	if casosActuales > 0 {
		casosPredichos += casosActuales
	}

	// Determinar nivel de riesgo
	nivelRiesgo := "BAJO"
	switch {
	case probabilidad >= 80:
		nivelRiesgo = "CRÍTICO"
	case probabilidad >= 60:
		nivelRiesgo = "ALTO"
	case probabilidad >= 40:
		nivelRiesgo = "MEDIO"
	}

	return PrediccionPropagacion{
		Distrito:        distrito,
		FechaPrediccion: time.Now().AddDate(0, 0, 7), // Predicción a 7 días
		CasosPredichos:  casosPredichos,
		Probabilidad:    math.Round(probabilidad*100) / 100,
		NivelRiesgo:     nivelRiesgo,
	}
}

func (s *PropagacionService) generarRecomendaciones(distritos []DistritoAfectado, velocidadPromedio, factorDensidad float64) []string {
	var recomendaciones []string

	// Análisis de la situación actual
	if velocidadPromedio > 5 {
		recomendaciones = append(recomendaciones, "⚠️ ALERTA: Velocidad de propagación alta (>5 casos/día). Implementar medidas de contención inmediatas.")
	}

	if factorDensidad > 10 {
		recomendaciones = append(recomendaciones, "🏙️ Enfocar esfuerzos en distritos de alta densidad poblacional como Plan Tres Mil y Norte.")
	}

	// Recomendaciones por distrito de alto riesgo
	for _, distrito := range distritos {
		if distrito.RiesgoExpansion == "CRÍTICO" {
			recomendaciones = append(recomendaciones, fmt.Sprintf("🚨 %s: Riesgo crítico - Establecer cerco epidemiológico y aumentar vigilancia.", distrito.Distrito))
		}
	}

	// Recomendaciones generales
	recomendaciones = append(recomendaciones, "📍 Intensificar vigilancia epidemiológica en distritos conectados a focos activos.")
	recomendaciones = append(recomendaciones, "🏥 Redistribuir recursos médicos según patrones de propagación identificados.")
	recomendaciones = append(recomendaciones, "📊 Actualizar análisis cada 48-72 horas para ajustar estrategias de contención.")

	return recomendaciones
}

// GetSpreadPredictionsByDistrict obtiene predicciones específicas para un distrito
func (s *PropagacionService) GetSpreadPredictionsByDistrict(distrito, enfermedad string, dias int) (*PrediccionPropagacion, error) {
	analisis, err := s.AnalyzeSpreadVelocity(enfermedad, dias)
	if err != nil {
		return nil, err
	}

	for _, prediccion := range analisis.PredictedSpread {
		if prediccion.Distrito == distrito {
			return &prediccion, nil
		}
	}

	return nil, fmt.Errorf("no se encontraron predicciones para el distrito %s", distrito)
}