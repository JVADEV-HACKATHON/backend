package services

import (
	"errors"

	"hospital-api/internal/database"
	"hospital-api/internal/models"
	"hospital-api/internal/utils"

	"gorm.io/gorm"
)

type HospitalService struct {
	db *gorm.DB
}

// NewHospitalService crea una nueva instancia del servicio de hospitales
func NewHospitalService() *HospitalService {
	return &HospitalService{
		db: database.GetDB(),
	}
}

// GetAllHospitales obtiene todos los hospitales con paginación
func (s *HospitalService) GetAllHospitales(page, limit int) ([]models.Hospital, int64, error) {
	var hospitales []models.Hospital
	var total int64

	// Contar total de registros
	if err := s.db.Model(&models.Hospital{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener registros con paginación
	offset := (page - 1) * limit
	if err := s.db.Offset(offset).Limit(limit).Find(&hospitales).Error; err != nil {
		return nil, 0, err
	}

	return hospitales, total, nil
}

// GetHospitalByID obtiene un hospital por su ID
func (s *HospitalService) GetHospitalByID(id uint) (*models.Hospital, error) {
	var hospital models.Hospital

	if err := s.db.First(&hospital, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hospital no encontrado")
		}
		return nil, err
	}

	return &hospital, nil
}

// GetHospitalesNearby obtiene hospitales cercanos a unas coordenadas usando la fórmula de Haversine
func (s *HospitalService) GetHospitalesNearby(lat, lng, radius float64) ([]models.HospitalResponse, error) {
	var hospitales []models.Hospital

	// Obtener todos los hospitales
	if err := s.db.Find(&hospitales).Error; err != nil {
		return nil, err
	}

	var hospitalesCercanos []models.HospitalResponse

	// Filtrar hospitales dentro del radio especificado
	for _, hospital := range hospitales {
		distancia := utils.CalcularDistanciaHaversine(lat, lng, hospital.Latitud, hospital.Longitud)

		if distancia <= radius {
			hospitalResponse := hospital.ToResponse()
			// Agregar la distancia como información adicional
			hospitalesCercanos = append(hospitalesCercanos, hospitalResponse)
		}
	}

	return hospitalesCercanos, nil
}

// GetHospitalesWithDistances obtiene todos los hospitales con sus distancias a un punto
func (s *HospitalService) GetHospitalesWithDistances(lat, lng float64) ([]HospitalWithDistance, error) {
	var hospitales []models.Hospital

	if err := s.db.Find(&hospitales).Error; err != nil {
		return nil, err
	}

	var hospitalesConDistancia []HospitalWithDistance

	for _, hospital := range hospitales {
		distancia := utils.CalcularDistanciaHaversine(lat, lng, hospital.Latitud, hospital.Longitud)

		hospitalConDistancia := HospitalWithDistance{
			Hospital:  hospital.ToResponse(),
			Distancia: distancia,
		}

		hospitalesConDistancia = append(hospitalesConDistancia, hospitalConDistancia)
	}

	return hospitalesConDistancia, nil
}

// HospitalWithDistance estructura para incluir distancia
type HospitalWithDistance struct {
	Hospital  models.HospitalResponse `json:"hospital"`
	Distancia float64                 `json:"distancia_km"`
}

// ValidateHospitalCoordinates valida que las coordenadas del hospital estén en Santa Cruz
func (s *HospitalService) ValidateHospitalCoordinates(lat, lng float64) error {
	// Límites aproximados de Santa Cruz de la Sierra, Bolivia
	// Latitud: -17.9 a -17.7
	// Longitud: -63.3 a -63.0
	if lat < -17.9 || lat > -17.7 || lng < -63.3 || lng > -63.0 {
		return errors.New("las coordenadas del hospital deben estar ubicadas en Santa Cruz de la Sierra, Bolivia")
	}
	return nil
}

// UpdateHospitalLocation actualiza la ubicación de un hospital
func (s *HospitalService) UpdateHospitalLocation(hospitalID uint, lat, lng float64, direccion string) error {
	// Validar coordenadas
	if err := s.ValidateHospitalCoordinates(lat, lng); err != nil {
		return err
	}

	// Actualizar en la base de datos
	result := s.db.Model(&models.Hospital{}).
		Where("id = ?", hospitalID).
		Updates(map[string]interface{}{
			"latitud":   lat,
			"longitud":  lng,
			"direccion": direccion,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("hospital no encontrado")
	}

	return nil
}
func (s *HospitalService) GetAllHospitalesSinPaginacion() ([]models.HospitalResponse, error) {
	var hospitales []models.Hospital

	if err := s.db.Find(&hospitales).Error; err != nil {
		return nil, err
	}

	// Convertir a HospitalResponse para no exponer información sensible
	var hospitalesResponse []models.HospitalResponse
	for _, hospital := range hospitales {
		hospitalesResponse = append(hospitalesResponse, hospital.ToResponse())
	}

	return hospitalesResponse, nil
}

// HospitalWithPatientsCount estructura para incluir el conteo de pacientes
type HospitalWithPatientsCount struct {
	Hospital      models.HospitalResponse `json:"hospital"`
	TotalPacientes int64                  `json:"total_pacientes"`
}

// GetAllHospitalesWithPatientsCount obtiene todos los hospitales con el conteo de pacientes únicos
func (s *HospitalService) GetAllHospitalesWithPatientsCount() ([]HospitalWithPatientsCount, error) {
	var hospitales []models.Hospital

	// Obtener todos los hospitales
	if err := s.db.Find(&hospitales).Error; err != nil {
		return nil, err
	}

	var hospitalesConConteo []HospitalWithPatientsCount

	// Para cada hospital, contar los pacientes únicos que han tenido historial clínico
	for _, hospital := range hospitales {
		var totalPacientes int64
		
		// Contar pacientes únicos que han tenido historial clínico en este hospital
		err := s.db.Model(&models.HistorialClinico{}).
			Where("id_hospital = ?", hospital.ID).
			Distinct("id_paciente").
			Count(&totalPacientes).Error

		if err != nil {
			return nil, err
		}

		hospitalConConteo := HospitalWithPatientsCount{
			Hospital:       hospital.ToResponse(),
			TotalPacientes: totalPacientes,
		}

		hospitalesConConteo = append(hospitalesConConteo, hospitalConConteo)
	}

	return hospitalesConConteo, nil
}

// GetAllHospitalesWithPatientsCountPaginated obtiene todos los hospitales con el conteo de pacientes únicos con paginación
func (s *HospitalService) GetAllHospitalesWithPatientsCountPaginated(page, limit int) ([]HospitalWithPatientsCount, int64, error) {
	var hospitales []models.Hospital
	var totalHospitales int64

	// Contar total de hospitales
	if err := s.db.Model(&models.Hospital{}).Count(&totalHospitales).Error; err != nil {
		return nil, 0, err
	}

	// Obtener hospitales con paginación
	offset := (page - 1) * limit
	if err := s.db.Offset(offset).Limit(limit).Find(&hospitales).Error; err != nil {
		return nil, 0, err
	}

	var hospitalesConConteo []HospitalWithPatientsCount

	// Para cada hospital, contar los pacientes únicos
	for _, hospital := range hospitales {
		var totalPacientes int64
		
		err := s.db.Model(&models.HistorialClinico{}).
			Where("id_hospital = ?", hospital.ID).
			Distinct("id_paciente").
			Count(&totalPacientes).Error

		if err != nil {
			return nil, 0, err
		}

		hospitalConConteo := HospitalWithPatientsCount{
			Hospital:       hospital.ToResponse(),
			TotalPacientes: totalPacientes,
		}

		hospitalesConConteo = append(hospitalesConConteo, hospitalConConteo)
	}

	return hospitalesConConteo, totalHospitales, nil
}

// GetHospitalWithPatientsCountByID obtiene un hospital específico con el conteo de pacientes
func (s *HospitalService) GetHospitalWithPatientsCountByID(hospitalID uint) (*HospitalWithPatientsCount, error) {
	var hospital models.Hospital

	// Obtener el hospital
	if err := s.db.First(&hospital, hospitalID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hospital no encontrado")
		}
		return nil, err
	}

	// Contar pacientes únicos
	var totalPacientes int64
	err := s.db.Model(&models.HistorialClinico{}).
		Where("id_hospital = ?", hospital.ID).
		Distinct("id_paciente").
		Count(&totalPacientes).Error

	if err != nil {
		return nil, err
	}

	hospitalConConteo := &HospitalWithPatientsCount{
		Hospital:       hospital.ToResponse(),
		TotalPacientes: totalPacientes,
	}

	return hospitalConConteo, nil
}