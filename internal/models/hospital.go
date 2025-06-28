package models

import (
	"time"

	"gorm.io/gorm"
)

// Hospital representa la tabla de hospitales
type Hospital struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre    string         `json:"nombre" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	Direccion string         `json:"direccion" gorm:"type:varchar(200);not null" validate:"required,min=5,max=200"`
	Ciudad    string         `json:"ciudad" gorm:"type:varchar(50);not null" validate:"required,min=2,max=50"`
	Telefono  string         `json:"telefono" gorm:"type:varchar(20);unique"`
	Email     string         `json:"email" gorm:"type:varchar(100);unique;not null" validate:"required,email"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"` // No se incluye en JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	HistorialesClinico []HistorialClinico `json:"historiales_clinico,omitempty" gorm:"foreignKey:IDHospital"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (Hospital) TableName() string {
	return "hospitales"
}

// HospitalResponse es la estructura para respuestas sin información sensible
type HospitalResponse struct {
	ID        uint      `json:"id"`
	Nombre    string    `json:"nombre"`
	Direccion string    `json:"direccion"`
	Ciudad    string    `json:"ciudad"`
	Telefono  string    `json:"telefono"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse convierte Hospital a HospitalResponse
func (h *Hospital) ToResponse() HospitalResponse {
	return HospitalResponse{
		ID:        h.ID,
		Nombre:    h.Nombre,
		Direccion: h.Direccion,
		Ciudad:    h.Ciudad,
		Telefono:  h.Telefono,
		Email:     h.Email,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}
