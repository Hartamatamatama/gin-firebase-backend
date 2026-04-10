package main

import (
	"log"
	"os"

	"github.com/hartamatamatama/gin-firebase-backend/config"
	"github.com/hartamatamatama/gin-firebase-backend/routes"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables dari file .env
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// 2. Inisialisasi Firebase Admin SDK
	config.InitFirebase()

	// 3. Inisialisasi database + AutoMigrate GORM
	config.InitDatabase()

	// 4. Setup Gin router dengan semua routes
	router := routes.SetupRouter()

	// 5. Jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/v1/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}