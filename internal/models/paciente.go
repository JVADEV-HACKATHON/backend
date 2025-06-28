package models

import (
	"time"

	"gorm.io/gorm"
)

// Paciente representa la tabla de pacientes
type Paciente struct {
	ID              uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre          string         `json:"nombre" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	FechaNacimiento time.Time      `json:"fecha_nacimiento" gorm:"type:date;not null" validate:"required"`
	Sexo            string         `json:"sexo" gorm:"type:varchar(1);not null;check:sexo IN ('M','F','O')" validate:"required,oneof=M F O"`
	TipoSangre      string         `json:"tipo_sangre" gorm:"type:varchar(4)" validate:"omitempty,max=4"`
	PesoKg          float64        `json:"peso_kg" gorm:"type:decimal(5,2);check:peso_kg > 0" validate:"omitempty,gt=0"`
	AlturaCm        int            `json:"altura_cm" gorm:"type:int;check:altura_cm > 0" validate:"omitempty,gt=0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	HistorialesClinico []HistorialClinico `json:"historiales_clinico,omitempty" gorm:"foreignKey:IDPaciente"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (Paciente) TableName() string {
	return "pacientes"
}

// GetAge calcula la edad del paciente
func (p *Paciente) GetAge() int {
	now := time.Now()
	age := now.Year() - p.FechaNacimiento.Year()

	// Ajustar si el cumpleaños no ha pasado este año
	if now.YearDay() < p.FechaNacimiento.YearDay() {
		age--
	}

	return age
}
