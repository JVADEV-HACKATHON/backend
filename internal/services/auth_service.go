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

type LoginResponse struct {
	Token    string                  `json:"token"`
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
