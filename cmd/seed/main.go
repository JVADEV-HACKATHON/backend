package main

import (
	"flag"
	"log"
	"os"

	"hospital-api/internal/config"
	"hospital-api/internal/database"
	"hospital-api/internal/seeders"
)

func main() {
	// Configurar flags de línea de comandos
	clean := flag.Bool("clean", false, "Limpiar la base de datos antes del seeding")
	help := flag.Bool("help", false, "Mostrar ayuda")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	log.Println("🌱 Iniciando proceso de seeding...")

	// Cargar configuración
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Error cargando configuración: %v", err)
	}

	// Conectar a la base de datos
	database.ConnectDatabase()

	// Crear instancia del seeder
	seeder := seeders.NewSeeder()

	// Limpiar base de datos si se especifica
	if *clean {
		if err := seeder.CleanDatabase(); err != nil {
			log.Fatalf("❌ Error limpiando la base de datos: %v", err)
		}
	}

	// Ejecutar seeding
	if err := seeder.SeedAll(); err != nil {
		log.Fatalf("❌ Error ejecutando seeding: %v", err)
	}

	log.Println("🎉 Proceso de seeding completado exitosamente!")
	log.Println("📋 Datos disponibles:")
	log.Println("   🏥 15 Hospitales de Santa Cruz de la Sierra")
	log.Println("   👥 15 Pacientes")
	log.Println("   📊 12+ Historiales clínicos con datos geográficos")
	log.Println("   🗺️  Datos listos para mapas de calor epidemiológicos")
}

func printHelp() {
	log.Println("🌱 Hospital API Database Seeder")
	log.Println("")
	log.Println("Uso:")
	log.Println("  go run cmd/seed/main.go [flags]")
	log.Println("")
	log.Println("Flags:")
	log.Println("  -clean    Limpiar la base de datos antes del seeding")
	log.Println("  -help     Mostrar esta ayuda")
	log.Println("")
	log.Println("Ejemplos:")
	log.Println("  go run cmd/seed/main.go")
	log.Println("  go run cmd/seed/main.go -clean")
	log.Println("")
	log.Println("Datos que se insertarán:")
	log.Println("  🏥 15 Hospitales de Santa Cruz de la Sierra")
	log.Println("  👥 15 Pacientes con datos realistas")
	log.Println("  📊 12+ Historiales clínicos con coordenadas reales")
	log.Println("  🗺️  Datos geográficos para mapas de calor")
}
