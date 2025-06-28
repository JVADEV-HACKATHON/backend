package services

import (
	"errors"
	"time"

	"hospital-api/internal/database"
	"hospital-api/internal/models"

	"gorm.io/gorm"
)

type HistorialService struct {
	db *gorm.DB
}

type HeatMapData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Count     int64   `json:"count"`
	District  string  `json:"district"`
}

type EpidemiologicalStats struct {
	TotalCases      int64           `json:"total_cases"`
	ContagiousCases int64           `json:"contagious_cases"`
	ByDistrict      []DistrictStats `json:"by_district"`
	ByDate          []DateStats     `json:"by_date"`
	HeatMapData     []HeatMapData   `json:"heat_map_data"`
}

type DistrictStats struct {
	District        string `json:"district"`
	TotalCases      int64  `json:"total_cases"`
	ContagiousCases int64  `json:"contagious_cases"`
}

type DateStats struct {
	Date            string `json:"date"`
	TotalCases      int64  `json:"total_cases"`
	ContagiousCases int64  `json:"contagious_cases"`
}

// NewHistorialService crea una nueva instancia del servicio de historial clínico
func NewHistorialService() *HistorialService {
	return &HistorialService{
		db: database.GetDB(),
	}
}

// CreateHistorial crea un nuevo registro de historial clínico
func (s *HistorialService) CreateHistorial(historial *models.HistorialClinico) error {
	return s.db.Create(historial).Error
}

// GetHistorialByID obtiene un historial por ID con información relacionada
func (s *HistorialService) GetHistorialByID(id uint) (*models.HistorialClinico, error) {
	var historial models.HistorialClinico
	err := s.db.Preload("Paciente").Preload("Hospital").First(&historial, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("historial clínico no encontrado")
		}
		return nil, err
	}
	return &historial, nil
}

// GetHistorialByPaciente obtiene el historial clínico de un paciente
func (s *HistorialService) GetHistorialByPaciente(pacienteID uint, page, limit int) ([]models.HistorialClinico, int64, error) {
	var historiales []models.HistorialClinico
	var total int64

	query := s.db.Where("id_paciente = ?", pacienteID)

	// Contar total
	query.Model(&models.HistorialClinico{}).Count(&total)

	// Obtener registros con paginación y relaciones
	offset := (page - 1) * limit
	err := query.Preload("Hospital").
		Offset(offset).
		Limit(limit).
		Order("fecha_ingreso DESC").
		Find(&historiales).Error

	return historiales, total, err
}

// GetHistorialByHospital obtiene el historial clínico de un hospital
func (s *HistorialService) GetHistorialByHospital(hospitalID uint, page, limit int) ([]models.HistorialClinico, int64, error) {
	var historiales []models.HistorialClinico
	var total int64

	query := s.db.Where("id_hospital = ?", hospitalID)

	// Contar total
	query.Model(&models.HistorialClinico{}).Count(&total)

	// Obtener registros con paginación y relaciones
	offset := (page - 1) * limit
	err := query.Preload("Paciente").
		Offset(offset).
		Limit(limit).
		Order("fecha_ingreso DESC").
		Find(&historiales).Error

	return historiales, total, err
}

// UpdateHistorial actualiza un historial clínico
func (s *HistorialService) UpdateHistorial(id uint, updates *models.HistorialClinico) error {
	return s.db.Model(&models.HistorialClinico{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteHistorial elimina un historial clínico (soft delete)
func (s *HistorialService) DeleteHistorial(id uint) error {
	return s.db.Delete(&models.HistorialClinico{}, id).Error
}

// GetEpidemiologicalStats obtiene estadísticas epidemiológicas para mapas de calor
func (s *HistorialService) GetEpidemiologicalStats(startDate, endDate time.Time) (*EpidemiologicalStats, error) {
	stats := &EpidemiologicalStats{}

	// Total de casos
	s.db.Model(&models.HistorialClinico{}).
		Where("consultation_date BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalCases)

	// Casos contagiosos
	s.db.Model(&models.HistorialClinico{}).
		Where("consultation_date BETWEEN ? AND ? AND is_contagious = ?", startDate, endDate, true).
		Count(&stats.ContagiousCases)

	// Estadísticas por distrito
	var districtStats []DistrictStats
	s.db.Model(&models.HistorialClinico{}).
		Select("patient_district as district, COUNT(*) as total_cases, COUNT(CASE WHEN is_contagious = true THEN 1 END) as contagious_cases").
		Where("consultation_date BETWEEN ? AND ?", startDate, endDate).
		Group("patient_district").
		Scan(&districtStats)
	stats.ByDistrict = districtStats

	// Estadísticas por fecha
	var dateStats []DateStats
	s.db.Model(&models.HistorialClinico{}).
		Select("consultation_date::date as date, COUNT(*) as total_cases, COUNT(CASE WHEN is_contagious = true THEN 1 END) as contagious_cases").
		Where("consultation_date BETWEEN ? AND ?", startDate, endDate).
		Group("consultation_date::date").
		Order("date").
		Scan(&dateStats)
	stats.ByDate = dateStats

	// Datos para mapa de calor (agrupado por coordenadas aproximadas)
	var heatMapData []HeatMapData
	s.db.Model(&models.HistorialClinico{}).
		Select("ROUND(patient_latitude::numeric, 4) as latitude, ROUND(patient_longitude::numeric, 4) as longitude, patient_district as district, COUNT(*) as count").
		Where("consultation_date BETWEEN ? AND ?", startDate, endDate).
		Group("ROUND(patient_latitude::numeric, 4), ROUND(patient_longitude::numeric, 4), patient_district").
		Scan(&heatMapData)
	stats.HeatMapData = heatMapData

	return stats, nil
}

// GetContagiousHistorial obtiene historiales de casos contagiosos
func (s *HistorialService) GetContagiousHistorial(page, limit int) ([]models.HistorialClinico, int64, error) {
	var historiales []models.HistorialClinico
	var total int64

	query := s.db.Where("is_contagious = ?", true)

	// Contar total
	query.Model(&models.HistorialClinico{}).Count(&total)

	// Obtener registros con paginación y relaciones
	offset := (page - 1) * limit
	err := query.Preload("Paciente").
		Preload("Hospital").
		Offset(offset).
		Limit(limit).
		Order("fecha_ingreso DESC").
		Find(&historiales).Error

	return historiales, total, err
}
