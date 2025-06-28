package models

import "time"

// HistorialEnfermedadResponse estructura específica para la respuesta del endpoint de búsqueda por enfermedad
type HistorialEnfermedadResponse struct {
	ID                  uint                   `json:"id"`
	FechaIngreso        time.Time              `json:"fecha_ingreso"`
	MotivoConsulta      string                 `json:"motivo_consulta"`
	Diagnostico         string                 `json:"diagnostico"`
	Tratamiento         string                 `json:"tratamiento"`
	Medicamentos        string                 `json:"medicamentos"`
	Observaciones       string                 `json:"observaciones"`
	PatientLatitude     float64                `json:"patient_latitude"`
	PatientLongitude    float64                `json:"patient_longitude"`
	PatientAddress      string                 `json:"patient_address"`
	PatientDistrict     string                 `json:"patient_district"`
	PatientNeighborhood string                 `json:"patient_neighborhood"`
	ConsultationDate    time.Time              `json:"consultation_date"`
	SymptomsStartDate   *time.Time             `json:"symptoms_start_date"`
	IsContagious        bool                   `json:"is_contagious"`
	CreatedAt           time.Time              `json:"created_at"`
	Paciente            PacienteEnfermedadInfo `json:"paciente"`
	Hospital            HospitalEnfermedadInfo `json:"hospital"`
}

// PacienteEnfermedadInfo información simplificada del paciente
type PacienteEnfermedadInfo struct {
	ID       uint   `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Edad     int    `json:"edad"`
	Sexo     string `json:"sexo"`
}

// HospitalEnfermedadInfo información simplificada del hospital
type HospitalEnfermedadInfo struct {
	ID                uint    `json:"id"`
	Nombre            string  `json:"nombre"`
	Direccion         string  `json:"direccion"`
	HospitalLatitude  float64 `json:"hospital_latitude"`
	HospitalLongitude float64 `json:"hospital_longitude"`
	Ciudad            string  `json:"ciudad"`
	Telefono          string  `json:"telefono"`
}

// EnfermedadSearchResponse estructura para la respuesta completa del endpoint
type EnfermedadSearchResponse struct {
	Message string                        `json:"message"`
	Total   int64                         `json:"total"`
	Data    []HistorialEnfermedadResponse `json:"data"`
}

// ToEnfermedadResponse convierte HistorialClinico a HistorialEnfermedadResponse
func (h *HistorialClinico) ToEnfermedadResponse() HistorialEnfermedadResponse {
	// Procesar nombre del paciente
	nombre := ""
	apellido := ""
	if h.Paciente.Nombre != "" {
		// Dividir el nombre completo en nombre y apellido
		// Asumiendo formato: "Nombre Apellido1 Apellido2"
		parts := parseFullName(h.Paciente.Nombre)
		if len(parts) >= 2 {
			nombre = parts[0]
			apellido = parts[1]
			if len(parts) > 2 {
				apellido += " " + parts[2]
			}
		} else if len(parts) == 1 {
			nombre = parts[0]
			apellido = ""
		}
	}

	return HistorialEnfermedadResponse{
		ID:                  h.ID,
		FechaIngreso:        h.FechaIngreso,
		MotivoConsulta:      h.MotivoConsulta,
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
		Paciente: PacienteEnfermedadInfo{
			ID:       h.Paciente.ID,
			Nombre:   nombre,
			Apellido: apellido,
			Edad:     h.Paciente.GetAge(),
			Sexo:     h.Paciente.Sexo,
		},
		Hospital: HospitalEnfermedadInfo{
			ID:                h.Hospital.ID,
			Nombre:            h.Hospital.Nombre,
			Direccion:         h.Hospital.Direccion,
			HospitalLatitude:  h.Hospital.Latitud,
			HospitalLongitude: h.Hospital.Longitud,
			Ciudad:            h.Hospital.Ciudad,
			Telefono:          h.Hospital.Telefono,
		},
	}
}

// parseFullName divide un nombre completo en partes
func parseFullName(fullName string) []string {
	// Implementación simple para dividir el nombre
	parts := make([]string, 0)
	current := ""

	for _, char := range fullName {
		if char == ' ' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
