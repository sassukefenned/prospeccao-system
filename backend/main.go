package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port            string
	Host            string
	DBPath          string
	ScraperPath     string
	ContactEmail    string
	ContactPhone    string
}

var config AppConfig
var db *Database

func main() {
	// Carregar .env
	godotenv.Load()

	// Config
	config = AppConfig{
		Port:         getEnv("PORT", "8080"),
		Host:         getEnv("HOST", "localhost"),
		DBPath:       getEnv("DB_PATH", "./data/prospeccao.db"),
		ScraperPath:  getEnv("SCRAPER_PATH", "./gmapsscraper"),
		ContactEmail: getEnv("CONTACT_EMAIL", "stfenned@gmail.com"),
		ContactPhone: getEnv("CONTACT_PHONE", "+55 62 8178-9507"),
	}

	// Inicializar BD
	var err error
	db, err = NewDatabase(config.DBPath)
	if err != nil {
		log.Fatalf("Erro ao conectar BD: %v", err)
	}
	defer db.Close()

	// Criar tabelas se não existirem
	if err := db.Init(); err != nil {
		log.Fatalf("Erro ao criar tabelas: %v", err)
	}

	// Fiber app
	app := fiber.New()

	// Middleware
	app.Use(fiber.Logger())

	// Static files
	app.Static("/", "./frontend")

	// API Routes
	api := app.Group("/api")

	// Leads endpoints
	api.Get("/leads", handleGetLeads)
	api.Post("/leads", handleCreateLead)
	api.Put("/leads/:id", handleUpdateLead)
	api.Delete("/leads/:id", handleDeleteLead)

	// Search endpoints
	api.Post("/search", handleNewSearch)
	api.Get("/search/:id", handleGetSearch)
	api.Get("/searches", handleListSearches)

	// Campaign endpoints
	api.Post("/campaigns", handleCreateCampaign)
	api.Get("/campaigns", handleListCampaigns)
	api.Get("/campaigns/:id", handleGetCampaign)

	// Export endpoint
	api.Get("/export/csv", handleExportCSV)
	api.Get("/export/json", handleExportJSON)

	// Stats endpoint
	api.Get("/stats", handleStats)

	// Start server
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Printf("🚀 Servidor iniciado em http://%s", addr)
	log.Fatal(app.Listen(addr))
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
