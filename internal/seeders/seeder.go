package seeders

import (
	"log"
	"time"

	"hospital-api/internal/database"
	"hospital-api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

// NewSeeder crea una nueva instancia del seeder
func NewSeeder() *Seeder {
	return &Seeder{
		db: database.GetDB(),
	}
}

// SeedAll ejecuta todos los seeders en orden
func (s *Seeder) SeedAll() error {
	log.Println("🌱 Iniciando seeding de la base de datos...")

	// Orden de seeding (importante por las relaciones)
	if err := s.SeedHospitales(); err != nil {
		return err
	}

	if err := s.SeedPacientes(); err != nil {
		return err
	}

	if err := s.SeedHistorialesClinico(); err != nil {
		return err
	}

	log.Println("✅ Seeding completado exitosamente!")
	return nil
}

// CleanDatabase limpia todas las tablas
func (s *Seeder) CleanDatabase() error {
	log.Println("🧹 Limpiando base de datos...")

	// Orden de limpieza (inverso por las relaciones)
	if err := s.db.Exec("DELETE FROM historial_clinico").Error; err != nil {
		return err
	}

	if err := s.db.Exec("DELETE FROM pacientes").Error; err != nil {
		return err
	}

	if err := s.db.Exec("DELETE FROM hospitales").Error; err != nil {
		return err
	}

	// Reiniciar secuencias
	if err := s.db.Exec("ALTER SEQUENCE historial_clinico_id_seq RESTART WITH 1").Error; err != nil {
		return err
	}

	if err := s.db.Exec("ALTER SEQUENCE pacientes_id_seq RESTART WITH 1").Error; err != nil {
		return err
	}

	if err := s.db.Exec("ALTER SEQUENCE hospitales_id_seq RESTART WITH 1").Error; err != nil {
		return err
	}

	log.Println("✅ Base de datos limpiada!")
	return nil
}

// SeedHospitales inserta datos de hospitales
func (s *Seeder) SeedHospitales() error {
	log.Println("🏥 Seeding hospitales...")

	// Hashear contraseñas
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hospitales := []models.Hospital{
		{
			Nombre:    "Hospital Central de La Paz",
			Direccion: "Av. Saavedra 2302, Miraflores",
			Ciudad:    "La Paz",
			Telefono:  "+591-2-2651100",
			Email:     "admin@hospitalcentral.com",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital del Niño Dr. Ovidio Aliaga Uría",
			Direccion: "Av. Simón Bolívar 1800, Miraflores",
			Ciudad:    "La Paz",
			Telefono:  "+591-2-2456789",
			Email:     "admin@hospitalnino.com",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital de Clínicas",
			Direccion: "Av. Saavedra 2302, Sopocachi",
			Ciudad:    "La Paz",
			Telefono:  "+591-2-2459080",
			Email:     "admin@hospitalclinicas.com",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital San Gabriel",
			Direccion: "Calle Capitán Ravelo 2048, San Pedro",
			Ciudad:    "La Paz",
			Telefono:  "+591-2-2401122",
			Email:     "admin@hospitalsangabriel.com",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital Arco Iris",
			Direccion: "Av. 6 de Agosto 2821, San Miguel",
			Ciudad:    "La Paz",
			Telefono:  "+591-2-2431234",
			Email:     "admin@hospitalarcoiris.com",
			Password:  string(hashedPassword),
		},
	}

	for _, hospital := range hospitales {
		// Verificar si ya existe
		var existingHospital models.Hospital
		result := s.db.Where("email = ?", hospital.Email).First(&existingHospital)

		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			// No existe, crear nuevo
			if err := s.db.Create(&hospital).Error; err != nil {
				return err
			}
			log.Printf("✅ Hospital creado: %s", hospital.Nombre)
		} else {
			log.Printf("⚠️  Hospital ya existe: %s", hospital.Nombre)
		}
	}

	return nil
}

// SeedPacientes inserta datos de pacientes
func (s *Seeder) SeedPacientes() error {
	log.Println("👥 Seeding pacientes...")

	pacientes := []models.Paciente{
		{
			Nombre:          "Juan Carlos Mendoza",
			FechaNacimiento: time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "O+",
			PesoKg:          75.5,
			AlturaCm:        175,
		},
		{
			Nombre:          "María Elena Quispe",
			FechaNacimiento: time.Date(1992, 7, 22, 0, 0, 0, 0, time.UTC),
			Sexo:            "F",
			TipoSangre:      "A+",
			PesoKg:          62.3,
			AlturaCm:        160,
		},
		{
			Nombre:          "Carlos Alberto Mamani",
			FechaNacimiento: time.Date(1978, 11, 8, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "B+",
			PesoKg:          82.1,
			AlturaCm:        180,
		},
		{
			Nombre:          "Ana Lucía Vargas",
			FechaNacimiento: time.Date(1995, 1, 30, 0, 0, 0, 0, time.UTC),
			Sexo:            "F",
			TipoSangre:      "AB+",
			PesoKg:          58.7,
			AlturaCm:        165,
		},
		{
			Nombre:          "Pedro Antonio Choque",
			FechaNacimiento: time.Date(1967, 9, 14, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "O-",
			PesoKg:          78.9,
			AlturaCm:        172,
		},
		{
			Nombre:          "Rosa Elena Condori",
			FechaNacimiento: time.Date(1988, 4, 12, 0, 0, 0, 0, time.UTC),
			Sexo:            "F",
			TipoSangre:      "A-",
			PesoKg:          65.2,
			AlturaCm:        158,
		},
		{
			Nombre:          "Miguel Ángel Torrez",
			FechaNacimiento: time.Date(1982, 6, 25, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "B-",
			PesoKg:          71.4,
			AlturaCm:        177,
		},
		{
			Nombre:          "Claudia Patricia Flores",
			FechaNacimiento: time.Date(1990, 12, 3, 0, 0, 0, 0, time.UTC),
			Sexo:            "F",
			TipoSangre:      "AB-",
			PesoKg:          59.8,
			AlturaCm:        162,
		},
		{
			Nombre:          "Luis Fernando Ramos",
			FechaNacimiento: time.Date(1975, 8, 17, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "O+",
			PesoKg:          85.3,
			AlturaCm:        183,
		},
		{
			Nombre:          "Silvia Mónica Cruz",
			FechaNacimiento: time.Date(1993, 5, 9, 0, 0, 0, 0, time.UTC),
			Sexo:            "F",
			TipoSangre:      "A+",
			PesoKg:          63.7,
			AlturaCm:        167,
		},
		{
			Nombre:          "Roberto Inca Huanca",
			FechaNacimiento: time.Date(1980, 10, 21, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "B+",
			PesoKg:          79.6,
			AlturaCm:        174,
		},
		{
			Nombre:          "Patricia Luz Arroyo",
			FechaNacimiento: time.Date(1987, 2, 28, 0, 0, 0, 0, time.UTC),
			Sexo:            "F",
			TipoSangre:      "O-",
			PesoKg:          61.9,
			AlturaCm:        159,
		},
		{
			Nombre:          "Fernando José Marca",
			FechaNacimiento: time.Date(1972, 12, 7, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "A-",
			PesoKg:          73.8,
			AlturaCm:        176,
		},
		{
			Nombre:          "Verónica Isabel Morales",
			FechaNacimiento: time.Date(1991, 9, 16, 0, 0, 0, 0, time.UTC),
			Sexo:            "F",
			TipoSangre:      "AB+",
			PesoKg:          56.4,
			AlturaCm:        163,
		},
		{
			Nombre:          "Diego Alejandro Poma",
			FechaNacimiento: time.Date(1984, 1, 11, 0, 0, 0, 0, time.UTC),
			Sexo:            "M",
			TipoSangre:      "B-",
			PesoKg:          81.2,
			AlturaCm:        179,
		},
	}

	for _, paciente := range pacientes {
		// Verificar si ya existe
		var existingPaciente models.Paciente
		result := s.db.Where("nombre = ? AND fecha_nacimiento = ?",
			paciente.Nombre, paciente.FechaNacimiento).First(&existingPaciente)

		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			// No existe, crear nuevo
			if err := s.db.Create(&paciente).Error; err != nil {
				return err
			}
			log.Printf("✅ Paciente creado: %s", paciente.Nombre)
		} else {
			log.Printf("⚠️  Paciente ya existe: %s", paciente.Nombre)
		}
	}

	return nil
}

// SeedHistorialesClinico inserta datos de historiales clínicos
func (s *Seeder) SeedHistorialesClinico() error {
	log.Println("📋 Seeding historiales clínicos...")

	// Obtener algunos IDs existentes de hospitales y pacientes
	var hospitales []models.Hospital
	if err := s.db.Find(&hospitales).Error; err != nil {
		return err
	}

	var pacientes []models.Paciente
	if err := s.db.Find(&pacientes).Error; err != nil {
		return err
	}

	if len(hospitales) == 0 || len(pacientes) == 0 {
		log.Println("⚠️  No hay hospitales o pacientes para crear historiales")
		return nil
	}

	// Datos de historiales clínicos con coordenadas reales de La Paz
	historialesData := []struct {
		IDPaciente          uint
		IDHospital          uint
		FechaIngreso        time.Time
		MotivoConsulta      string
		Diagnostico         string
		Tratamiento         string
		Medicamentos        string
		Observaciones       string
		PatientLatitude     float64
		PatientLongitude    float64
		PatientAddress      string
		PatientDistrict     string
		PatientNeighborhood string
		ConsultationDate    time.Time
		SymptomsStartDate   *time.Time
		IsContagious        bool
	}{
		{
			IDPaciente:          1,
			IDHospital:          1,
			FechaIngreso:        time.Now().AddDate(0, 0, -5),
			MotivoConsulta:      "Dolor de cabeza intenso y fiebre",
			Diagnostico:         "Migraña tensional",
			Tratamiento:         "Reposo y analgésicos",
			Medicamentos:        "Ibuprofeno 400mg cada 8 horas",
			Observaciones:       "Paciente estable, mejora progresiva",
			PatientLatitude:     -16.5000,
			PatientLongitude:    -68.1193,
			PatientAddress:      "Av. El Prado 1425, Centro",
			PatientDistrict:     "Centro",
			PatientNeighborhood: "El Prado",
			ConsultationDate:    time.Now().AddDate(0, 0, -5),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -7)),
			IsContagious:        false,
		},
		{
			IDPaciente:          2,
			IDHospital:          1,
			FechaIngreso:        time.Now().AddDate(0, 0, -3),
			MotivoConsulta:      "Tos persistente y dificultad respiratoria",
			Diagnostico:         "Bronquitis aguda",
			Tratamiento:         "Broncodilatadores y antibióticos",
			Medicamentos:        "Salbutamol, Amoxicilina 500mg",
			Observaciones:       "Evolución favorable, control en 7 días",
			PatientLatitude:     -16.5097,
			PatientLongitude:    -68.1192,
			PatientAddress:      "Calle Sagárnaga 318, Centro",
			PatientDistrict:     "Centro",
			PatientNeighborhood: "Rosario",
			ConsultationDate:    time.Now().AddDate(0, 0, -3),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -5)),
			IsContagious:        true,
		},
		{
			IDPaciente:          3,
			IDHospital:          2,
			FechaIngreso:        time.Now().AddDate(0, 0, -10),
			MotivoConsulta:      "Dolor abdominal severo",
			Diagnostico:         "Gastritis aguda",
			Tratamiento:         "Dieta blanda y protectores gástricos",
			Medicamentos:        "Omeprazol 20mg, Sucralfato",
			Observaciones:       "Paciente en ayunas, mejora notable",
			PatientLatitude:     -16.5322,
			PatientLongitude:    -68.0753,
			PatientAddress:      "Av. 6 de Agosto 2420, San Miguel",
			PatientDistrict:     "San Miguel",
			PatientNeighborhood: "Villa Fátima",
			ConsultationDate:    time.Now().AddDate(0, 0, -10),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -12)),
			IsContagious:        false,
		},
		{
			IDPaciente:          4,
			IDHospital:          2,
			FechaIngreso:        time.Now().AddDate(0, 0, -8),
			MotivoConsulta:      "Erupción cutánea y picazón",
			Diagnostico:         "Dermatitis alérgica",
			Tratamiento:         "Antihistamínicos y corticoides tópicos",
			Medicamentos:        "Loratadina 10mg, Betametasona crema",
			Observaciones:       "Reacción alérgica a alimento, mejoría evidente",
			PatientLatitude:     -16.5203,
			PatientLongitude:    -68.1127,
			PatientAddress:      "Av. Saavedra 1950, Sopocachi",
			PatientDistrict:     "Sopocachi",
			PatientNeighborhood: "Sopocachi",
			ConsultationDate:    time.Now().AddDate(0, 0, -8),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -9)),
			IsContagious:        false,
		},
		{
			IDPaciente:          5,
			IDHospital:          3,
			FechaIngreso:        time.Now().AddDate(0, 0, -15),
			MotivoConsulta:      "Fiebre alta y malestar general",
			Diagnostico:         "Síndrome gripal",
			Tratamiento:         "Reposo absoluto y sintomáticos",
			Medicamentos:        "Paracetamol 500mg, abundantes líquidos",
			Observaciones:       "Cuadro viral, evolución satisfactoria",
			PatientLatitude:     -16.4955,
			PatientLongitude:    -68.1336,
			PatientAddress:      "Calle 21 de Calacoto 1234, Calacoto",
			PatientDistrict:     "Calacoto",
			PatientNeighborhood: "Calacoto",
			ConsultationDate:    time.Now().AddDate(0, 0, -15),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -17)),
			IsContagious:        true,
		},
		{
			IDPaciente:          6,
			IDHospital:          3,
			FechaIngreso:        time.Now().AddDate(0, 0, -2),
			MotivoConsulta:      "Lesión en tobillo por caída",
			Diagnostico:         "Esguince de tobillo grado II",
			Tratamiento:         "Inmovilización y fisioterapia",
			Medicamentos:        "Diclofenaco 50mg, hielo local",
			Observaciones:       "Inflamación moderada, pronóstico favorable",
			PatientLatitude:     -16.5408,
			PatientLongitude:    -68.0619,
			PatientAddress:      "Av. Costanera s/n, Achachicala",
			PatientDistrict:     "Max Paredes",
			PatientNeighborhood: "Achachicala",
			ConsultationDate:    time.Now().AddDate(0, 0, -2),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -2)),
			IsContagious:        false,
		},
		{
			IDPaciente:          7,
			IDHospital:          4,
			FechaIngreso:        time.Now().AddDate(0, 0, -6),
			MotivoConsulta:      "Conjuntivitis bilateral",
			Diagnostico:         "Conjuntivitis viral",
			Tratamiento:         "Lágrimas artificiales y compresas frías",
			Medicamentos:        "Tobramicina gotas oftálmicas",
			Observaciones:       "Proceso viral autolimitado, buena evolución",
			PatientLatitude:     -16.4894,
			PatientLongitude:    -68.1317,
			PatientAddress:      "Calle 10 de Obrajes 890, Obrajes",
			PatientDistrict:     "Obrajes",
			PatientNeighborhood: "Obrajes",
			ConsultationDate:    time.Now().AddDate(0, 0, -6),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -8)),
			IsContagious:        true,
		},
		{
			IDPaciente:          8,
			IDHospital:          4,
			FechaIngreso:        time.Now().AddDate(0, 0, -12),
			MotivoConsulta:      "Hipertensión arterial descontrolada",
			Diagnostico:         "Crisis hipertensiva",
			Tratamiento:         "Control estricto y ajuste de medicación",
			Medicamentos:        "Enalapril 10mg, Amlodipino 5mg",
			Observaciones:       "Presión controlada, seguimiento ambulatorio",
			PatientLatitude:     -16.5189,
			PatientLongitude:    -68.0888,
			PatientAddress:      "Av. Buenos Aires 1567, Miraflores",
			PatientDistrict:     "Miraflores",
			PatientNeighborhood: "Miraflores",
			ConsultationDate:    time.Now().AddDate(0, 0, -12),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -14)),
			IsContagious:        false,
		},
		{
			IDPaciente:          9,
			IDHospital:          5,
			FechaIngreso:        time.Now().AddDate(0, 0, -1),
			MotivoConsulta:      "Gastroenteritis aguda",
			Diagnostico:         "Gastroenteritis viral",
			Tratamiento:         "Hidratación oral y dieta astringente",
			Medicamentos:        "Suero oral, Loperamida 2mg",
			Observaciones:       "Deshidratación leve, recuperación rápida",
			PatientLatitude:     -16.5075,
			PatientLongitude:    -68.1064,
			PatientAddress:      "Calle Landaeta 754, San Pedro",
			PatientDistrict:     "San Pedro",
			PatientNeighborhood: "San Pedro",
			ConsultationDate:    time.Now().AddDate(0, 0, -1),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -2)),
			IsContagious:        true,
		},
		{
			IDPaciente:          10,
			IDHospital:          5,
			FechaIngreso:        time.Now().AddDate(0, 0, -20),
			MotivoConsulta:      "Control prenatal rutinario",
			Diagnostico:         "Embarazo de 28 semanas normal",
			Tratamiento:         "Continuación de vitaminas prenatales",
			Medicamentos:        "Ácido fólico 5mg, Sulfato ferroso",
			Observaciones:       "Evolución normal del embarazo, próximo control en 4 semanas",
			PatientLatitude:     -16.4973,
			PatientLongitude:    -68.1245,
			PatientAddress:      "Av. Arce 2450, San Jorge",
			PatientDistrict:     "San Jorge",
			PatientNeighborhood: "San Jorge",
			ConsultationDate:    time.Now().AddDate(0, 0, -20),
			SymptomsStartDate:   nil,
			IsContagious:        false,
		},
		// Casos contagiosos adicionales para estadísticas epidemiológicas
		{
			IDPaciente:          11,
			IDHospital:          1,
			FechaIngreso:        time.Now().AddDate(0, 0, -4),
			MotivoConsulta:      "Síntomas de COVID-19",
			Diagnostico:         "COVID-19 leve",
			Tratamiento:         "Aislamiento domiciliario y sintomáticos",
			Medicamentos:        "Paracetamol, Ibuprofeno, Vitamina D",
			Observaciones:       "Caso confirmado por PCR, contactos rastreados",
			PatientLatitude:     -16.5245,
			PatientLongitude:    -68.0516,
			PatientAddress:      "Calle Final Landaeta 123, Villa San Antonio",
			PatientDistrict:     "Villa San Antonio",
			PatientNeighborhood: "Villa San Antonio",
			ConsultationDate:    time.Now().AddDate(0, 0, -4),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -6)),
			IsContagious:        true,
		},
		{
			IDPaciente:          12,
			IDHospital:          2,
			FechaIngreso:        time.Now().AddDate(0, 0, -7),
			MotivoConsulta:      "Diarrea y vómitos persistentes",
			Diagnostico:         "Intoxicación alimentaria",
			Tratamiento:         "Hidratación endovenosa",
			Medicamentos:        "Suero fisiológico, Metoclopramida",
			Observaciones:       "Posible brote familiar, notificado a epidemiología",
			PatientLatitude:     -16.4856,
			PatientLongitude:    -68.1589,
			PatientAddress:      "Av. Hernando Siles 5890, Zona Sur",
			PatientDistrict:     "Zona Sur",
			PatientNeighborhood: "La Florida",
			ConsultationDate:    time.Now().AddDate(0, 0, -7),
			SymptomsStartDate:   timePtr(time.Now().AddDate(0, 0, -8)),
			IsContagious:        true,
		},
	}

	for i, historialData := range historialesData {
		// Verificar que existan los IDs de paciente y hospital
		if int(historialData.IDPaciente) > len(pacientes) || int(historialData.IDHospital) > len(hospitales) {
			log.Printf("⚠️  Saltando historial %d: IDs inválidos", i+1)
			continue
		}

		historial := models.HistorialClinico{
			IDPaciente:          historialData.IDPaciente,
			IDHospital:          historialData.IDHospital,
			FechaIngreso:        historialData.FechaIngreso,
			MotivoConsulta:      historialData.MotivoConsulta,
			Diagnostico:         historialData.Diagnostico,
			Tratamiento:         historialData.Tratamiento,
			Medicamentos:        historialData.Medicamentos,
			Observaciones:       historialData.Observaciones,
			PatientLatitude:     historialData.PatientLatitude,
			PatientLongitude:    historialData.PatientLongitude,
			PatientAddress:      historialData.PatientAddress,
			PatientDistrict:     historialData.PatientDistrict,
			PatientNeighborhood: historialData.PatientNeighborhood,
			ConsultationDate:    historialData.ConsultationDate,
			SymptomsStartDate:   historialData.SymptomsStartDate,
			IsContagious:        historialData.IsContagious,
		}

		// Verificar si ya existe un historial similar
		var existingHistorial models.HistorialClinico
		result := s.db.Where("id_paciente = ? AND id_hospital = ? AND fecha_ingreso = ?",
			historial.IDPaciente, historial.IDHospital, historial.FechaIngreso).First(&existingHistorial)

		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			// No existe, crear nuevo
			if err := s.db.Create(&historial).Error; err != nil {
				return err
			}
			log.Printf("✅ Historial clínico creado: Paciente %d - %s", historial.IDPaciente, historial.MotivoConsulta)
		} else {
			log.Printf("⚠️  Historial ya existe: Paciente %d", historial.IDPaciente)
		}
	}

	return nil
}

// timePtr es una función helper para crear punteros a time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}
