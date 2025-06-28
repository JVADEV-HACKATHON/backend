package services

import (
	"errors"

	"hospital-api/internal/database"
	"hospital-api/internal/models"

	"gorm.io/gorm"
)

type PacienteService struct {
	db *gorm.DB
}

// NewPacienteService crea una nueva instancia del servicio de pacientes
func NewPacienteService() *PacienteService {
	return &PacienteService{
		db: database.GetDB(),
	}
}

// CreatePaciente crea un nuevo paciente
func (s *PacienteService) CreatePaciente(paciente *models.Paciente) error {
	return s.db.Create(paciente).Error
}

// GetPacienteByID obtiene un paciente por ID
func (s *PacienteService) GetPacienteByID(id uint) (*models.Paciente, error) {
	var paciente models.Paciente
	err := s.db.First(&paciente, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("paciente no encontrado")
		}
		return nil, err
	}
	return &paciente, nil
}

// GetAllPacientes obtiene todos los pacientes con paginación
func (s *PacienteService) GetAllPacientes(page, limit int) ([]models.Paciente, int64, error) {
	var pacientes []models.Paciente
	var total int64

	// Contar total de registros
	s.db.Model(&models.Paciente{}).Count(&total)

	// Obtener registros con paginación
	offset := (page - 1) * limit
	err := s.db.Offset(offset).Limit(limit).Find(&pacientes).Error

	return pacientes, total, err
}

// UpdatePaciente actualiza un paciente
func (s *PacienteService) UpdatePaciente(id uint, updates *models.Paciente) error {
	return s.db.Model(&models.Paciente{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePaciente elimina un paciente (soft delete)
func (s *PacienteService) DeletePaciente(id uint) error {
	return s.db.Delete(&models.Paciente{}, id).Error
}

// SearchPacientes busca pacientes por nombre
func (s *PacienteService) SearchPacientes(query string, page, limit int) ([]models.Paciente, int64, error) {
	var pacientes []models.Paciente
	var total int64

	// Buscar por nombre (insensible a mayúsculas)
	searchQuery := s.db.Where("LOWER(nombre) LIKE LOWER(?)", "%"+query+"%")

	// Contar total
	searchQuery.Model(&models.Paciente{}).Count(&total)

	// Obtener resultados con paginación
	offset := (page - 1) * limit
	err := searchQuery.Offset(offset).Limit(limit).Find(&pacientes).Error

	return pacientes, total, err
}
