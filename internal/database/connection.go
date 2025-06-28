package database

import (
	"fmt"
	"log"
	"os"

	"hospital-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase establece la conexión con la base de datos PostgreSQL
func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSL_MODE")

	// Configurar el DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Configurar el logger de GORM
	var gormLogger logger.Interface
	if os.Getenv("GIN_MODE") == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Establecer conexión
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Configurar pool de conexiones
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error al obtener la instancia de la base de datos: %v", err)
	}

	// Configurar parámetros del pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Conexión exitosa con la base de datos PostgreSQL")

	// Ejecutar migraciones automáticas
	err = AutoMigrate()
	if err != nil {
		log.Fatalf("Error en las migraciones: %v", err)
	}
}

// AutoMigrate ejecuta las migraciones automáticas de GORM
func AutoMigrate() error {
	log.Println("Ejecutando migraciones automáticas...")

	err := DB.AutoMigrate(
		&models.Hospital{},
		&models.Paciente{},
		&models.HistorialClinico{},
	)

	if err != nil {
		return fmt.Errorf("error en migración automática: %w", err)
	}

	log.Println("Migraciones completadas exitosamente")
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}
