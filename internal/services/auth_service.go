package services

import (
	"errors"
	"os"
	"time"

	"hospital-api/internal/database"
	"hospital-api/internal/middleware"
	"hospital-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Nombre    string  `json:"nombre" validate:"required,min=2,max=100"`
	Direccion string  `json:"direccion" validate:"required,min=5,max=200"`
	Latitud   float64 `json:"latitud" validate:"required,latitude"`
	Longitud  float64 `json:"longitud" validate:"required,longitude"`
	Ciudad    string  `json:"ciudad" validate:"required,min=2,max=50"`
	Telefono  string  `json:"telefono" validate:"omitempty,max=20"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token    string                  `json:"token"`
	Hospital models.HospitalResponse `json:"hospital"`
	Success  bool                    `json:"success"`
	Message  string                  `json:"message"`
}

type RegisterResponse struct {
	Hospital models.HospitalResponse `json:"hospital"`
	Success  bool                    `json:"success"`
	Message  string                  `json:"message"`
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

// Login autentica un hospital y retorna un JWT
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	var hospital models.Hospital

	// Buscar hospital por email
	err := s.db.Where("email = ?", req.Email).First(&hospital).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("credenciales inválidas")
		}
		return nil, err
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(hospital.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := s.generateJWT(hospital.ID, hospital.Email)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:    token,
		Hospital: hospital.ToResponse(),
		Success:  true,
		Message:  "Login exitoso",
	}, nil
}

// Register registra un nuevo hospital en el sistema
func (s *AuthService) Register(req RegisterRequest) (*RegisterResponse, error) {
	// Verificar si el email ya existe
	var existingHospital models.Hospital
	err := s.db.Where("email = ?", req.Email).First(&existingHospital).Error
	if err == nil {
		return nil, errors.New("el email ya está registrado")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Verificar si el teléfono ya existe (si se proporciona)
	if req.Telefono != "" {
		err = s.db.Where("telefono = ?", req.Telefono).First(&existingHospital).Error
		if err == nil {
			return nil, errors.New("el teléfono ya está registrado")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Hashear la contraseña
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	// Crear el nuevo hospital
	hospital := models.Hospital{
		Nombre:    req.Nombre,
		Direccion: req.Direccion,
		Latitud:   req.Latitud,
		Longitud:  req.Longitud,
		Ciudad:    req.Ciudad,
		Telefono:  req.Telefono,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	// Guardar en la base de datos
	err = s.db.Create(&hospital).Error
	if err != nil {
		return nil, errors.New("error al crear el hospital")
	}

	return &RegisterResponse{
		Hospital: hospital.ToResponse(),
		Success:  true,
		Message:  "Hospital registrado exitosamente",
	}, nil
}

// generateJWT genera un token JWT para el hospital
func (s *AuthService) generateJWT(hospitalID uint, email string) (string, error) {
	claims := &middleware.JWTClaims{
		HospitalID: hospitalID,
		Email:      email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// HashPassword hashea una contraseña usando bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
