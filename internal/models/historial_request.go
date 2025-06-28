package models

import "time"

// HistorialClinicoRequest estructura para recibir datos del frontend
type HistorialClinicoRequest struct {
	IDPaciente     uint      `json:"id_paciente" validate:"required"`
	FechaIngreso   time.Time `json:"fecha_ingreso" validate:"required"`
	MotivoConsulta string    `json:"motivo_consulta" validate:"required,min=3,max=200"`
	Diagnostico    string    `json:"diagnostico"`
	Tratamiento    string    `json:"tratamiento"`
	Medicamentos   string    `json:"medicamentos"`
	Observaciones  string    `json:"observaciones"`

	// Solo la dirección del frontend
	PatientAddress string `json:"patient_address" validate:"required,min=5,max=500"`

	// Campos opcionales si el frontend los envía
	PatientDistrict     string `json:"patient_district,omitempty"`
	PatientNeighborhood string `json:"patient_neighborhood,omitempty"`

	ConsultationDate  time.Time  `json:"consultation_date"`
	SymptomsStartDate *time.Time `json:"symptoms_start_date,omitempty"`
	IsContagious      bool       `json:"is_contagious"`
}

// ToHistorialClinico convierte el request a modelo de base de datos
func (r *HistorialClinicoRequest) ToHistorialClinico() *HistorialClinico {
	return &HistorialClinico{
		IDPaciente:          r.IDPaciente,
		FechaIngreso:        r.FechaIngreso,
		MotivoConsulta:      r.MotivoConsulta,
		Diagnostico:         r.Diagnostico,
		Tratamiento:         r.Tratamiento,
		Medicamentos:        r.Medicamentos,
		Observaciones:       r.Observaciones,
		PatientAddress:      r.PatientAddress,
		PatientDistrict:     r.PatientDistrict,
		PatientNeighborhood: r.PatientNeighborhood,
		ConsultationDate:    r.ConsultationDate,
		SymptomsStartDate:   r.SymptomsStartDate,
		IsContagious:        r.IsContagious,
	}
}
