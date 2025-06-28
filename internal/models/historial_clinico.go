package models

import (
	"time"

	"gorm.io/gorm"
)

// HistorialClinico representa la tabla de historial clínico
type HistorialClinico struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDPaciente     uint      `json:"id_paciente" gorm:"not null" validate:"required"`
	IDHospital     uint      `json:"id_hospital" gorm:"not null" validate:"required"`
	FechaIngreso   time.Time `json:"fecha_ingreso" gorm:"type:timestamp;not null" validate:"required"`
	MotivoConsulta string    `json:"motivo_consulta" gorm:"type:varchar(200);not null" validate:"required,min=3,max=200"`
	Enfermedad     string    `json:"enfermedad" gorm:"type:varchar(150);not null" validate:"required,min=2,max=150"`
	Diagnostico    string    `json:"diagnostico" gorm:"type:text"`
	Tratamiento    string    `json:"tratamiento" gorm:"type:text"`
	Medicamentos   string    `json:"medicamentos" gorm:"type:text"`
	Observaciones  string    `json:"observaciones" gorm:"type:text"`

	// Geolocalización crítica para mapas de calor
	PatientLatitude     float64 `json:"patient_latitude" gorm:"type:decimal(10,8);not null" validate:"required,latitude"`
	PatientLongitude    float64 `json:"patient_longitude" gorm:"type:decimal(11,8);not null" validate:"required,longitude"`
	PatientAddress      string  `json:"patient_address" gorm:"type:varchar(500);not null" validate:"required,min=5,max=500"`
	PatientDistrict     string  `json:"patient_district" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	PatientNeighborhood string  `json:"patient_neighborhood" gorm:"type:varchar(100)"`

	// Datos temporales
	ConsultationDate  time.Time  `json:"consultation_date" gorm:"type:date;not null;default:CURRENT_DATE"`
	SymptomsStartDate *time.Time `json:"symptoms_start_date" gorm:"type:date"`

	// Contexto epidemiológico
	IsContagious bool `json:"is_contagious" gorm:"default:false"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	Paciente Paciente `json:"paciente,omitempty" gorm:"foreignKey:IDPaciente"`
	Hospital Hospital `json:"hospital,omitempty" gorm:"foreignKey:IDHospital"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (HistorialClinico) TableName() string {
	return "historial_clinico"
}

// HistorialClinicoResponse estructura para respuestas con información relacionada
type HistorialClinicoResponse struct {
	ID                  uint       `json:"id"`
	FechaIngreso        time.Time  `json:"fecha_ingreso"`
	MotivoConsulta      string     `json:"motivo_consulta"`
	Enfermedad          string     `json:"enfermedad"`
	Diagnostico         string     `json:"diagnostico"`
	Tratamiento         string     `json:"tratamiento"`
	Medicamentos        string     `json:"medicamentos"`
	Observaciones       string     `json:"observaciones"`
	PatientLatitude     float64    `json:"patient_latitude"`
	PatientLongitude    float64    `json:"patient_longitude"`
	PatientAddress      string     `json:"patient_address"`
	PatientDistrict     string     `json:"patient_district"`
	PatientNeighborhood string     `json:"patient_neighborhood"`
	ConsultationDate    time.Time  `json:"consultation_date"`
	SymptomsStartDate   *time.Time `json:"symptoms_start_date"`
	IsContagious        bool       `json:"is_contagious"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// Información relacionada
	Paciente *Paciente `json:"paciente,omitempty"`
	Hospital *Hospital `json:"hospital,omitempty"`
}

// ToResponse convierte HistorialClinico a HistorialClinicoResponse
func (h *HistorialClinico) ToResponse() HistorialClinicoResponse {
	return HistorialClinicoResponse{
		ID:                  h.ID,
		FechaIngreso:        h.FechaIngreso,
		MotivoConsulta:      h.MotivoConsulta,
		Enfermedad:          h.Enfermedad,
		Diagnostico:         h.Diagnostico,
		Tratamiento:         h.Tratamiento,
		Medicamentos:        h.Medicamentos,
		Observaciones:       h.Observaciones,
		PatientLatitude:     h.PatientLatitude,
		PatientLongitude:    h.PatientLongitude,
		PatientAddress:      h.PatientAddress,
		PatientDistrict:     h.PatientDistrict,
		PatientNeighborhood: h.PatientNeighborhood,
		ConsultationDate:    h.ConsultationDate,
		SymptomsStartDate:   h.SymptomsStartDate,
		IsContagious:        h.IsContagious,
		CreatedAt:           h.CreatedAt,
		UpdatedAt:           h.UpdatedAt,
	}
}
