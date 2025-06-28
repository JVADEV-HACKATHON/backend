package seeders

import (
	"fmt"
	"hospital-api/internal/database"
	"hospital-api/internal/models"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seeder estructura principal para el seeding de datos
type Seeder struct {
	db *gorm.DB
}

// NewSeeder crea una nueva instancia del seeder
func NewSeeder() *Seeder {
	return &Seeder{
		db: database.GetDB(),
	}
}

// CleanDatabase limpia todas las tablas de la base de datos
func (s *Seeder) CleanDatabase() error {
	log.Println("🧹 Limpiando base de datos...")

	// Deshabilitar restricciones de clave foránea temporalmente
	if err := s.db.Exec("SET foreign_key_checks = 0").Error; err != nil {
		// Para PostgreSQL usar:
		if err := s.db.Exec("SET CONSTRAINTS ALL DEFERRED").Error; err != nil {
			log.Printf("⚠️ No se pudieron deshabilitar las restricciones: %v", err)
		}
	}

	// Limpiar tablas en orden para evitar problemas de claves foráneas
	tables := []string{"historial_clinicos", "pacientes", "hospitals"}

	for _, table := range tables {
		if err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			log.Printf("⚠️ Error limpiando tabla %s: %v", table, err)
		} else {
			log.Printf("✅ Tabla %s limpiada", table)
		}
	}

	// Rehabilitar restricciones de clave foránea
	if err := s.db.Exec("SET foreign_key_checks = 1").Error; err != nil {
		// Para PostgreSQL usar:
		if err := s.db.Exec("SET CONSTRAINTS ALL IMMEDIATE").Error; err != nil {
			log.Printf("⚠️ No se pudieron rehabilitar las restricciones: %v", err)
		}
	}

	log.Println("✅ Base de datos limpiada exitosamente")
	return nil
}

// SeedAll ejecuta todo el proceso de seeding
func (s *Seeder) SeedAll() error {
	log.Println("🌱 Iniciando proceso completo de seeding...")

	// Ejecutar seeding de datos aleatorios para Santa Cruz
	if err := s.SeedAllRandom(); err != nil {
		return fmt.Errorf("error en seeding aleatorio: %w", err)
	}

	// Mostrar estadísticas
	if err := s.ShowEnfermedadesStats(); err != nil {
		return fmt.Errorf("error mostrando estadísticas: %w", err)
	}

	log.Println("✅ Proceso de seeding completado exitosamente!")
	return nil
}

// Estructura para las direcciones de Santa Cruz
type DireccionSantaCruz struct {
	Direccion string
	Latitud   float64
	Longitud  float64
	Distrito  string
	Barrio    string
}

// Array con 10 direcciones distribuidas por Santa Cruz
var direccionesSantaCruz = []DireccionSantaCruz{
	{
		Direccion: "Av. San Martín 3456, Equipetrol",
		Latitud:   -17.7690416,
		Longitud:  -63.1956686,
		Distrito:  "Equipetrol",
		Barrio:    "Equipetrol Norte",
	},
	{
		Direccion: "Radial 10, Km 6.5, Zona Norte",
		Latitud:   -17.7987909,
		Longitud:  -63.210345,
		Distrito:  "Norte",
		Barrio:    "Las Palmas",
	},
	{
		Direccion: "Av. Grigotá 2890",
		Latitud:   -17.798792,
		Longitud:  -63.210345,
		Distrito:  "Plan Tres Mil",
		Barrio:    "Plan Tres Mil Centro",
	},

	{
		Direccion: "Av. Alemana 1245, Villa 1ro de Mayo",
		Latitud:   -17.7379806,
		Longitud:  -63.2484834,
		Distrito:  "Villa 1ro de Mayo",
		Barrio:    "Villa 1ro de Mayo",
	},
	{
		Direccion: "Av. Banzer Km 8, Zona Norte",
		Latitud:   -17.7379989,
		Longitud:  -63.1866809,
		Distrito:  "Norte",
		Barrio:    "Norte",
	},
	{
		Direccion: "Radial 27, Km 4, Zona Sur",
		Latitud:   -17.7441931,
		Longitud:  -63.1801563,
		Distrito:  "Sur",
		Barrio:    "Zona Sur",
	},
	{
		Direccion: "Av. Cristo Redentor 567, Zona Oeste",
		Latitud:   -17.7439533,
		Longitud:  -63.1756103,
		Distrito:  "Oeste",
		Barrio:    "Pampa de la Isla",
	},
	{
		Direccion: "Doble Vía La Guardia Km 12, Zona Este",
		Latitud:   -17.7728417,
		Longitud:  -63.2374135,
		Distrito:  "Este",
		Barrio:    "La Guardia",
	},
	{
		Direccion: "Av. Roca y Coronado 1890, Equipetrol Sur",
		Latitud:   -17.77286,
		Longitud:  -63.175611,
		Distrito:  "Equipetrol",
		Barrio:    "Equipetrol Sur",
	},
}

// Arrays de nombres y apellidos bolivianos
var nombresMasculinos = []string{
	"Carlos", "José", "Luis", "Miguel", "Juan", "Roberto", "Fernando", "Eduardo", "Diego", "Antonio",
	"Alejandro", "Francisco", "Manuel", "Rafael", "Ricardo", "Sergio", "Jorge", "Pedro", "Daniel", "Alberto",
	"Andrés", "Guillermo", "Mauricio", "Rodrigo", "Javier", "Óscar", "Víctor", "Raúl", "Pablo", "Álvaro",
	"Gonzalo", "Marcelo", "Rubén", "Sebastián", "Adrián", "Leonardo", "Martín", "Hugo", "Iván", "Cristian",
	"Nelson", "Wilson", "Ronald", "Ramiro", "Freddy", "Johnny", "Henry", "Jimmy", "Kevin", "Alex",
}

var nombresFemeninos = []string{
	"María", "Ana", "Carmen", "Rosa", "Elena", "Patricia", "Claudia", "Silvia", "Verónica", "Mónica",
	"Gabriela", "Andrea", "Paola", "Vanessa", "Roxana", "Carla", "Daniela", "Alejandra", "Fernanda", "Lucía",
	"Isabel", "Teresa", "Beatriz", "Esperanza", "Gloria", "Mirian", "Karina", "Lourdes", "Sandra", "Nancy",
	"Yolanda", "Sonia", "Lidia", "Graciela", "Delia", "Martha", "Julia", "Cristina", "Viviana", "Marcela",
	"Lorena", "Susana", "Irma", "Nora", "Laura", "Jessica", "Karen", "Evelyn", "Daysi", "Wendy",
}

var apellidos = []string{
	"Suárez", "Mendoza", "Gutiérrez", "Rodríguez", "González", "Martínez", "López", "García", "Pérez", "Sánchez",
	"Rocha", "Terceros", "Peña", "Rivero", "Soliz", "Antelo", "Barbery", "Justiniano", "Vaca", "Diez",
	"Salvatierra", "Morón", "Ribera", "Landivar", "Saavedra", "Parada", "Burgos", "Cronenbold", "Richter", "Roca",
	"Aguilar", "Monasterio", "Claure", "Añez", "Pedraza", "Melgar", "Hurtado", "Flores", "Vargas", "Mamani",
	"Quispe", "Choque", "Condori", "Torrez", "Ramos", "Cruz", "Huanca", "Arroyo", "Marca", "Morales",
	"Poma", "Silva", "Herrera", "Jiménez", "Castro", "Romero", "Fernández", "Ruiz", "Díaz", "Moreno",
	"Muñoz", "Álvarez", "Ramírez", "Torres", "Domínguez", "Vásquez", "Ramos", "Gil", "Serrano", "Blanco",
	"Molina", "Medina", "Guerrero", "Cortés", "Ibáñez", "Campos", "Rubio", "Vega", "Delgado", "Reyes",
}

var tiposSangre = []string{"O+", "O-", "A+", "A-", "B+", "B-", "AB+", "AB-"}

// Estructura para definir cada enfermedad con sus características específicas
type EnfermedadInfo struct {
	Nombre         string
	MotivoConsulta []string
	Diagnosticos   []string
	Tratamientos   []string
	Medicamentos   []string
	EsContagiosa   bool
	Observaciones  []string
}

// Definición de las 6 enfermedades específicas
var enfermedadesEspecificas = []EnfermedadInfo{
	{
		Nombre: "Dengue",
		MotivoConsulta: []string{
			"Fiebre alta y dolor de cabeza intenso",
			"Dolor muscular y articular severo",
			"Erupción cutánea y fiebre",
			"Malestar general y dolor retroocular",
		},
		Diagnosticos: []string{
			"Dengue clásico sin signos de alarma",
			"Dengue con signos de alarma",
			"Fiebre dengue típica",
		},
		Tratamientos: []string{
			"Reposo absoluto e hidratación oral",
			"Control de fiebre y monitoreo de signos vitales",
			"Hidratación endovenosa si es necesario",
		},
		Medicamentos: []string{
			"Paracetamol 500mg cada 6 horas",
			"Suero oral abundante",
			"Paracetamol 1g cada 8 horas (adultos)",
		},
		EsContagiosa: true,
		Observaciones: []string{
			"Paciente en vigilancia epidemiológica",
			"Control de plaquetas cada 24 horas",
			"Notificado a epidemiología departamental",
			"Familiar orientado sobre signos de alarma",
		},
	},
	{
		Nombre: "Sarampión",
		MotivoConsulta: []string{
			"Erupción cutánea generalizada y fiebre",
			"Tos, fiebre y manchas en la piel",
			"Conjuntivitis y erupción maculopapular",
			"Fiebre alta con exantema característico",
		},
		Diagnosticos: []string{
			"Sarampión confirmado por clínica",
			"Sarampión típico con exantema",
			"Sarampión con complicaciones menores",
		},
		Tratamientos: []string{
			"Aislamiento respiratorio y sintomáticos",
			"Soporte nutricional y vitamina A",
			"Manejo de complicaciones según evolución",
		},
		Medicamentos: []string{
			"Paracetamol para fiebre",
			"Vitamina A 200,000 UI dosis única",
			"Suero fisiológico para hidratación ocular",
		},
		EsContagiosa: true,
		Observaciones: []string{
			"Caso notificado inmediatamente a epidemiología",
			"Aislamiento respiratorio estricto",
			"Investigación epidemiológica de contactos",
			"Seguimiento por 21 días",
		},
	},
	{
		Nombre: "Zika",
		MotivoConsulta: []string{
			"Erupción cutánea con picazón leve",
			"Fiebre baja y dolor articular",
			"Conjuntivitis y exantema",
			"Dolor de cabeza y malestar general",
		},
		Diagnosticos: []string{
			"Zika virus confirmado",
			"Síndrome febril compatible con Zika",
			"Zika con manifestaciones típicas",
		},
		Tratamientos: []string{
			"Reposo y sintomáticos",
			"Hidratación adecuada",
			"Antihistamínicos para prurito",
		},
		Medicamentos: []string{
			"Paracetamol 500mg cada 8 horas",
			"Loratadina 10mg para picazón",
			"Abundantes líquidos",
		},
		EsContagiosa: true,
		Observaciones: []string{
			"Orientación sobre prevención de vectores",
			"Caso notificado a vigilancia epidemiológica",
			"Seguimiento especial si paciente embarazada",
			"Control de evolución a los 7 días",
		},
	},
	{
		Nombre: "Influenza",
		MotivoConsulta: []string{
			"Fiebre alta de inicio súbito",
			"Tos seca y dolor muscular",
			"Malestar general y cefalea intensa",
			"Síntomas respiratorios y fiebre",
		},
		Diagnosticos: []string{
			"Influenza A estacional",
			"Síndrome gripal por Influenza",
			"Influenza con complicaciones menores",
		},
		Tratamientos: []string{
			"Antivirales si se inicia temprano",
			"Reposo y sintomáticos",
			"Hidratación y control de fiebre",
		},
		Medicamentos: []string{
			"Oseltamivir 75mg cada 12 horas por 5 días",
			"Paracetamol 1g cada 8 horas",
			"Ibuprofeno 400mg cada 8 horas",
		},
		EsContagiosa: true,
		Observaciones: []string{
			"Aislamiento respiratorio por 7 días",
			"Vigilancia de complicaciones respiratorias",
			"Orientación sobre medidas preventivas",
			"Control si no mejora en 72 horas",
		},
	},
	{
		Nombre: "Gripe AH1N1",
		MotivoConsulta: []string{
			"Fiebre alta y dificultad respiratoria",
			"Tos persistente y malestar severo",
			"Síntomas gripales intensos",
			"Fiebre, tos y dolor muscular intenso",
		},
		Diagnosticos: []string{
			"Influenza AH1N1 confirmada",
			"Gripe AH1N1 con síntomas respiratorios",
			"Influenza pandémica AH1N1",
		},
		Tratamientos: []string{
			"Oseltamivir inmediato",
			"Aislamiento y monitoreo respiratorio",
			"Soporte ventilatorio si es necesario",
		},
		Medicamentos: []string{
			"Oseltamivir 75mg cada 12 horas por 5 días",
			"Paracetamol para control de fiebre",
			"Broncodilatadores si hay broncoespasmo",
		},
		EsContagiosa: true,
		Observaciones: []string{
			"Notificación inmediata obligatoria",
			"Aislamiento estricto por 7-10 días",
			"Monitoreo de saturación de oxígeno",
			"Seguimiento evolutivo diario",
		},
	},
	{
		Nombre: "Bronquitis",
		MotivoConsulta: []string{
			"Tos persistente con expectoración",
			"Dificultad respiratoria y tos",
			"Tos con flemas y malestar",
			"Dolor torácico y tos productiva",
		},
		Diagnosticos: []string{
			"Bronquitis aguda viral",
			"Bronquitis bacteriana",
			"Bronquitis con componente alérgico",
		},
		Tratamientos: []string{
			"Broncodilatadores y expectorantes",
			"Antibióticos si hay sobreinfección",
			"Fisioterapia respiratoria",
		},
		Medicamentos: []string{
			"Salbutamol inhalador cada 6 horas",
			"Ambroxol 30mg cada 8 horas",
			"Amoxicilina 500mg cada 8 horas si bacteriana",
		},
		EsContagiosa: false,
		Observaciones: []string{
			"Evitar irritantes respiratorios",
			"Hidratación abundante",
			"Control en 7 días si no mejora",
			"Educación sobre factores desencadenantes",
		},
	},
}

// SeedHospitalesSantaCruz inserta datos de hospitales de Santa Cruz de la Sierra
func (s *Seeder) SeedHospitalesSantaCruz() error {
	log.Println("🏥 Seeding hospitales de Santa Cruz de la Sierra...")

	// Hashear contraseñas
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hospitales := []models.Hospital{
		{
			Nombre:    "Hospital Japonés",
			Direccion: "Av. Japón s/n, 3er Anillo Externo",
			Latitud:   -17.7725285,
			Longitud:  -63.153871,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3460101",
			Email:     "admin@huj.org.bo",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital de Niños Dr. Mario Ortiz Suárez",
			Direccion: "Calle René Moreno 171, Centro",
			Latitud:   -17.7807346,
			Longitud:  -63.1890985,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3346969",
			Email:     "admin@hospitalninos.gob.bo",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital San Juan de Dios",
			Direccion: "Calle Junín 1248, Centro",
			Latitud:   -17.779344,
			Longitud:  -63.1887634,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3342777",
			Email:     "admin@hospitalsanjuan.org.bo",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital Percy Boland",
			Direccion: "Av. Santos Dumont 2do Anillo",
			Latitud:   -17.7783784,
			Longitud:  -63.1897871,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3462031",
			Email:     "admin@percyboland.com",
			Password:  string(hashedPassword),
		},

		{
			Nombre:    "Hospital Municipal Francés",
			Direccion: "Av. Grigotá 1946, Plan Tres Mil",
			Latitud:   -17.8518622,
			Longitud:  -63.2225207,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3480200",
			Email:     "admin@hospitalfrances.gob.bo",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital del Norte",
			Direccion: "Av. Banzer 6to Anillo, Zona Norte",
			Latitud:   -17.3487718,
			Longitud:  -66.1773225,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3555100",
			Email:     "admin@hospitalnorte.com.bo",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Clínica Foianini",
			Direccion: "Av. Alemana 6to Anillo",
			Latitud:   -17.7916862,
			Longitud:  -63.1824279,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3462100",
			Email:     "admin@foianini.org",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital La Católica",
			Direccion: "Calle Cristóbal de Mendoza 297, Centro",
			Latitud:   -17.7374565,
			Longitud:  -63.1923283,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3336633",
			Email:     "admin@lacatolica.com.bo",
			Password:  string(hashedPassword),
		},

		{
			Nombre:    "Hospital de la Mujer Dr. Percy Boland",
			Direccion: "Av. Alemana 341, Villa 1ro de Mayo",
			Latitud:   -17.7783784,
			Longitud:  -63.1897871,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3462800",
			Email:     "admin@hospitaldelamujer.gob.bo",
			Password:  string(hashedPassword),
		},
		{
			Nombre:    "Hospital General San Juan de Dios",
			Direccion: "Barrio San Juan, Villa 1ro de Mayo",
			Latitud:   -17.9757477,
			Longitud:  -67.1164299,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3480300",
			Email:     "admin@hospitalgeneralsanjuan.gob.bo",
			Password:  string(hashedPassword),
		},

		{
			Nombre:    "Hospital Corazón de Jesús",
			Direccion: "Calle Beni 738, Centro",
			Latitud:   -16.5669841,
			Longitud:  -68.226954,
			Ciudad:    "Santa Cruz de la Sierra",
			Telefono:  "+591-3-3337700",
			Email:     "admin@corazondejesus.org.bo",
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
			log.Printf("⚠  Hospital ya existe: %s", hospital.Nombre)
		}
	}

	return nil
}

// SeedRandomPacientes genera 500 pacientes aleatorios
func (s *Seeder) SeedRandomPacientes() error {
	log.Println("👥 Generando 500 pacientes aleatorios...")

	// Inicializar generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 500; i++ {
		// Generar sexo aleatorio
		sexo := "M"
		var nombre string
		if rand.Intn(2) == 0 {
			sexo = "F"
			nombre = nombresFemeninos[rand.Intn(len(nombresFemeninos))]
		} else {
			nombre = nombresMasculinos[rand.Intn(len(nombresMasculinos))]
		}

		// Generar nombre completo
		apellido1 := apellidos[rand.Intn(len(apellidos))]
		apellido2 := apellidos[rand.Intn(len(apellidos))]
		nombreCompleto := fmt.Sprintf("%s %s %s", nombre, apellido1, apellido2)

		// Generar fecha de nacimiento (entre 1950 y 2020)
		añoNacimiento := 1950 + rand.Intn(70)
		mesNacimiento := 1 + rand.Intn(12)
		diaNacimiento := 1 + rand.Intn(28)
		fechaNacimiento := time.Date(añoNacimiento, time.Month(mesNacimiento), diaNacimiento, 0, 0, 0, 0, time.UTC)

		// Generar datos físicos aleatorios
		peso := 45.0 + rand.Float64()*55.0 // Entre 45 y 100 kg
		altura := 150 + rand.Intn(50)      // Entre 150 y 200 cm
		tipoSangre := tiposSangre[rand.Intn(len(tiposSangre))]

		paciente := models.Paciente{
			Nombre:          nombreCompleto,
			FechaNacimiento: fechaNacimiento,
			Sexo:            sexo,
			TipoSangre:      tipoSangre,
			PesoKg:          float64(int(peso*10)) / 10, // Redondear a 1 decimal
			AlturaCm:        altura,
		}

		// Verificar si ya existe (muy improbable con nombres aleatorios)
		var existingPaciente models.Paciente
		result := s.db.Where("nombre = ?", paciente.Nombre).First(&existingPaciente)

		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			if err := s.db.Create(&paciente).Error; err != nil {
				return err
			}
			if (i+1)%50 == 0 {
				log.Printf("✅ Creados %d pacientes...", i+1)
			}
		}
	}

	log.Println("✅ 500 pacientes aleatorios creados exitosamente!")
	return nil
}

// SeedRandomHistoriales genera 100 historiales clínicos aleatorios
func (s *Seeder) SeedRandomHistoriales() error {
	log.Println("📋 Generando 100 historiales clínicos aleatorios...")

	// Obtener IDs de hospitales y pacientes existentes
	var hospitales []models.Hospital
	if err := s.db.Find(&hospitales).Error; err != nil {
		return err
	}

	var pacientes []models.Paciente
	if err := s.db.Find(&pacientes).Error; err != nil {
		return err
	}

	if len(hospitales) == 0 || len(pacientes) == 0 {
		return fmt.Errorf("no hay hospitales o pacientes suficientes para crear historiales")
	}

	// Inicializar generador
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		// Seleccionar paciente y hospital aleatorios
		paciente := pacientes[rand.Intn(len(pacientes))]
		hospital := hospitales[rand.Intn(len(hospitales))]

		// Seleccionar dirección aleatoria
		direccion := direccionesSantaCruz[rand.Intn(len(direccionesSantaCruz))]

		// Generar fecha de ingreso aleatoria (últimos 30 días)
		diasAtras := rand.Intn(30) + 1
		fechaIngreso := time.Now().AddDate(0, 0, -diasAtras)
		fechaConsulta := fechaIngreso

		// Fecha de inicio de síntomas (1-5 días antes de la consulta)
		diasSintomas := rand.Intn(5) + 1
		fechaSintomas := fechaConsulta.AddDate(0, 0, -diasSintomas)

		// Seleccionar enfermedad aleatoria de las 6 específicas
		enfermedadInfo := enfermedadesEspecificas[rand.Intn(len(enfermedadesEspecificas))]

		// Seleccionar datos específicos de la enfermedad
		motivoConsulta := enfermedadInfo.MotivoConsulta[rand.Intn(len(enfermedadInfo.MotivoConsulta))]
		diagnostico := enfermedadInfo.Diagnosticos[rand.Intn(len(enfermedadInfo.Diagnosticos))]
		tratamiento := enfermedadInfo.Tratamientos[rand.Intn(len(enfermedadInfo.Tratamientos))]
		medicamento := enfermedadInfo.Medicamentos[rand.Intn(len(enfermedadInfo.Medicamentos))]
		observacion := enfermedadInfo.Observaciones[rand.Intn(len(enfermedadInfo.Observaciones))]

		// La contagiosidad depende de la enfermedad
		esContagioso := enfermedadInfo.EsContagiosa

		historial := models.HistorialClinico{
			IDPaciente:          paciente.ID,
			IDHospital:          hospital.ID,
			FechaIngreso:        fechaIngreso,
			MotivoConsulta:      motivoConsulta,
			Enfermedad:          enfermedadInfo.Nombre,
			Diagnostico:         diagnostico,
			Tratamiento:         tratamiento,
			Medicamentos:        medicamento,
			Observaciones:       observacion,
			PatientLatitude:     direccion.Latitud,
			PatientLongitude:    direccion.Longitud,
			PatientAddress:      direccion.Direccion,
			PatientDistrict:     direccion.Distrito,
			PatientNeighborhood: direccion.Barrio,
			ConsultationDate:    fechaConsulta,
			SymptomsStartDate:   &fechaSintomas,
			IsContagious:        esContagioso,
		}

		if err := s.db.Create(&historial).Error; err != nil {
			return err
		}

		if (i+1)%20 == 0 {
			log.Printf("✅ Creados %d historiales...", i+1)
		}
	}

	log.Println("✅ 100 historiales clínicos aleatorios creados exitosamente!")
	return nil
}

// SeedAllRandom ejecuta la generación de datos aleatorios
func (s *Seeder) SeedAllRandom() error {
	log.Println("🌱 Iniciando generación de datos aleatorios para Santa Cruz...")

	// Primero sembrar hospitales si no existen
	if err := s.SeedHospitalesSantaCruz(); err != nil {
		return err
	}

	// Generar pacientes aleatorios
	if err := s.SeedRandomPacientes(); err != nil {
		return err
	}

	// Generar historiales aleatorios
	if err := s.SeedRandomHistoriales(); err != nil {
		return err
	}

	log.Println("✅ Generación de datos aleatorios completada!")
	return nil
}

// ShowEnfermedadesStats muestra estadísticas de las enfermedades generadas
func (s *Seeder) ShowEnfermedadesStats() error {
	log.Println("📊 Estadísticas de enfermedades generadas:")

	for _, enfermedad := range enfermedadesEspecificas {
		var count int64
		if err := s.db.Model(&models.HistorialClinico{}).Where("enfermedad = ?", enfermedad.Nombre).Count(&count).Error; err != nil {
			return err
		}

		contagiosaStr := "No contagiosa"
		if enfermedad.EsContagiosa {
			contagiosaStr = "Contagiosa"
		}

		log.Printf("   - %s: %d casos (%s)", enfermedad.Nombre, count, contagiosaStr)
	}

	// Estadísticas por distrito
	log.Println("\n📍 Distribución por distritos:")
	for _, direccion := range direccionesSantaCruz {
		var count int64
		if err := s.db.Model(&models.HistorialClinico{}).Where("patient_district = ?", direccion.Distrito).Count(&count).Error; err != nil {
			return err
		}
		log.Printf("   - %s: %d casos", direccion.Distrito, count)
	}

	// Casos contagiosos vs no contagiosos
	var contagiosos, noContagiosos int64
	s.db.Model(&models.HistorialClinico{}).Where("is_contagious = ?", true).Count(&contagiosos)
	s.db.Model(&models.HistorialClinico{}).Where("is_contagious = ?", false).Count(&noContagiosos)

	log.Printf("\n🦠 Casos contagiosos: %d", contagiosos)
	log.Printf("🏥 Casos no contagiosos: %d", noContagiosos)

	return nil
}
