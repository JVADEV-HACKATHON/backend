# Hospital API 🏥

Una API REST completa desarrollada con Go y Gin para la gestión de sistemas hospitalarios con capacidades epidemiológicas y mapas de calor.

## 🚀 Características Principales

- **Autenticación JWT** para hospitales
- **CRUD completo** para pacientes y historiales clínicos
- **Geolocalización** para mapas de calor epidemiológicos
- **Estadísticas epidemiológicas** en tiempo real
- **Base de datos PostgreSQL** con GORM
- **Dockerizado** para fácil despliegue
- **Arquitectura escalable** con separación de responsabilidades
- **Validaciones robustas** con struct tags
- **Middleware personalizado** para logging y CORS
- **Paginación** en todos los endpoints de listado

## 🏗️ Arquitectura

```
.
├── cmd/
│   └── server/          # Punto de entrada de la aplicación
├── internal/
│   ├── config/          # Configuración de la aplicación
│   ├── database/        # Conexión y migraciones de BD
│   ├── handlers/        # Controladores HTTP
│   ├── middleware/      # Middleware personalizado
│   ├── models/          # Modelos de datos GORM
│   ├── routes/          # Definición de rutas
│   ├── services/        # Lógica de negocio
│   └── utils/           # Utilidades y helpers
├── scripts/             # Scripts de inicialización
├── docker-compose.yml   # Configuración Docker
├── Dockerfile          # Imagen Docker de la app
└── Makefile           # Comandos de desarrollo
```

## 📊 Modelo de Datos

### Hospital

- ID, nombre, dirección, ciudad, teléfono
- Email y password hasheado para autenticación
- Relación con historiales clínicos

### Paciente

- Información personal (nombre, fecha nacimiento, sexo)
- Datos médicos (tipo sangre, peso, altura)
- Relación con historiales clínicos

### Historial Clínico

- Información médica (motivo, diagnóstico, tratamiento)
- **Geolocalización crítica** (latitud, longitud, dirección, distrito)
- Datos epidemiológicos (fecha síntomas, es contagioso)
- Relaciones con paciente y hospital

## 🐳 Inicio Rápido con Docker

1. **Clona el repositorio:**

```bash
git clone <repo-url>
cd hospital-api
```

2. **Configura las variables de entorno:**

```bash
cp .env.example .env
# Edita .env con tus valores si es necesario
```

3. **Inicia con Docker Compose:**

```bash
make docker-run
# o
docker-compose up -d
```

4. **Verifica que esté funcionando:**

```bash
curl http://localhost:8080/api/v1/health
```

## 💻 Desarrollo Local

### Prerrequisitos

- Go 1.24.4+
- PostgreSQL 16+
- Make (opcional, para comandos automatizados)

### Configuración

1. **Instala dependencias:**

```bash
make deps
# o
go mod download && go mod tidy
```

2. **Inicia PostgreSQL:**

```bash
# Solo la base de datos
docker-compose up -d db
```

3. **Ejecuta la aplicación:**

```bash
make dev
# o
go run ./cmd/server/main.go
```

## 🔑 API Endpoints

### Autenticación

```bash
# Login de hospital
POST /api/v1/auth/login
{
  "email": "admin@hospitalcentral.com",
  "password": "admin123"
}

# Obtener perfil (requiere JWT)
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

### Pacientes

```bash
# Crear paciente
POST /api/v1/pacientes
Authorization: Bearer <token>
{
  "nombre": "Juan Pérez",
  "fecha_nacimiento": "1990-05-15T00:00:00Z",
  "sexo": "M",
  "tipo_sangre": "O+",
  "peso_kg": 75.5,
  "altura_cm": 175
}

# Listar pacientes (paginado)
GET /api/v1/pacientes?page=1&limit=10

# Buscar pacientes
GET /api/v1/pacientes/search?q=Juan&page=1&limit=10

# Obtener paciente por ID
GET /api/v1/pacientes/1

# Actualizar paciente
PUT /api/v1/pacientes/1

# Eliminar paciente
DELETE /api/v1/pacientes/1
```

### Historial Clínico

```bash
# Crear historial clínico
POST /api/v1/historial
Authorization: Bearer <token>
{
  "id_paciente": 1,
  "fecha_ingreso": "2023-12-01T10:00:00Z",
  "motivo_consulta": "Dolor de cabeza severo",
  "diagnostico": "Migraña",
  "tratamiento": "Reposo y analgésicos",
  "medicamentos": "Ibuprofeno 400mg",
  "observaciones": "Paciente estable",
  "patient_latitude": -12.0464,
  "patient_longitude": -77.0428,
  "patient_address": "Av. Lima 123, San Isidro",
  "patient_district": "San Isidro",
  "patient_neighborhood": "Centro",
  "consultation_date": "2023-12-01",
  "symptoms_start_date": "2023-11-30",
  "is_contagious": false
}

# Historial por paciente
GET /api/v1/historial/paciente/1?page=1&limit=10

# Historial por hospital actual
GET /api/v1/historial/hospital?page=1&limit=10

# Obtener historial específico
GET /api/v1/historial/1

# Actualizar historial
PUT /api/v1/historial/1

# Eliminar historial
DELETE /api/v1/historial/1
```

### Epidemiología

```bash
# Estadísticas epidemiológicas para mapas de calor
GET /api/v1/epidemiologia/stats?start_date=2023-11-01&end_date=2023-12-01

# Casos contagiosos
GET /api/v1/epidemiologia/contagious?page=1&limit=10
```

## 🗺️ Mapas de Calor

La API proporciona datos georreferenciados para crear mapas de calor:

```json
{
  "success": true,
  "data": {
    "total_cases": 150,
    "contagious_cases": 23,
    "heat_map_data": [
      {
        "latitude": -12.0464,
        "longitude": -77.0428,
        "count": 15,
        "district": "San Isidro"
      }
    ],
    "by_district": [
      {
        "district": "San Isidro",
        "total_cases": 45,
        "contagious_cases": 8
      }
    ],
    "by_date": [
      {
        "date": "2023-12-01",
        "total_cases": 12,
        "contagious_cases": 3
      }
    ]
  }
}
```

## 🛠️ Comandos Útiles

```bash
# Desarrollo
make dev                 # Inicia entorno de desarrollo
make run                 # Ejecuta la aplicación
make build              # Compila la aplicación
make test               # Ejecuta pruebas

# Docker
make docker-run         # Inicia con Docker Compose
make docker-stop        # Detiene contenedores
make docker-logs        # Muestra logs
make docker-rebuild     # Reconstruye contenedores

# Base de datos
make db-reset          # Reinicia la base de datos

# Limpieza
make clean             # Limpia archivos compilados
```

## 🔒 Seguridad

- **JWT Tokens** con expiración de 24 horas
- **Contraseñas hasheadas** con bcrypt
- **Validación de entrada** en todos los endpoints
- **CORS configurado** para producción
- **Middleware de autenticación** en rutas protegidas

## 📝 Variables de Entorno

```bash
# Base de datos
DB_HOST=db
DB_PORT=5432
DB_USER=hospital_user
DB_PASSWORD=hospital_password
DB_NAME=hospital_db
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Servidor
PORT=8080
GIN_MODE=debug
API_VERSION=v1
```

## 🏥 Usuario de Prueba

El sistema incluye un hospital de prueba:

- **Email:** `admin@hospitalcentral.com`
- **Password:** `admin123`

## 📈 Próximas Características

- [ ] Swagger/OpenAPI documentation
- [ ] Rate limiting
- [ ] Caché con Redis
- [ ] Notificaciones en tiempo real
- [ ] Exportación de reportes
- [ ] Integración con servicios de mapas
- [ ] Dashboard web

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 👥 Autores

- **Tu Nombre** - _Desarrollo inicial_ - [TuGitHub](https://github.com/tuusername)

## 🙏 Agradecimientos

- Gin Web Framework
- GORM ORM
- PostgreSQL
- Docker
- Toda la comunidad de Go

---

⭐ Si este proyecto te fue útil, no olvides darle una estrella!
