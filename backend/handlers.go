package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

// ===== LEADS =====

func handleGetLeads(c fiber.Ctx) error {
	filtros := &FiltrosLead{
		Status:   c.Query("status", ""),
		Qualidade: c.Query("qualidade", ""),
		Perfil:   c.Query("perfil", ""),
		Categoria: c.Query("categoria", ""),
		Busca:    c.Query("busca", ""),
		Limit:    getIntQuery(c, "limit", 50),
		Offset:   getIntQuery(c, "offset", 0),
	}

	leads, err := db.ListLeads(filtros)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.JSON(leads)
}

func handleCreateLead(c fiber.Ctx) error {
	lead := new(Lead)
	if err := c.BindJSON(lead); err != nil {
		return c.Status(400).JSON(fiber.Map{"erro": err.Error()})
	}

	if lead.Nome == "" || lead.Telefone == "" {
		return c.Status(400).JSON(fiber.Map{"erro": "Nome e telefone são obrigatórios"})
	}

	if err := db.CreateLead(lead); err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.Status(201).JSON(lead)
}

func handleUpdateLead(c fiber.Ctx) error {
	id := c.Params("id")
	lead := new(Lead)
	if err := c.BindJSON(lead); err != nil {
		return c.Status(400).JSON(fiber.Map{"erro": err.Error()})
	}

	lead.ID = id
	if err := db.UpdateLead(lead); err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.JSON(lead)
}

func handleDeleteLead(c fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DeleteLead(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.JSON(fiber.Map{"mensagem": "Lead deletado com sucesso"})
}

// ===== SEARCH =====

func handleNewSearch(c fiber.Ctx) error {
	req := struct {
		Termo     string `json:"termo"`
		Categoria string `json:"categoria"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"erro": err.Error()})
	}

	if req.Termo == "" {
		return c.Status(400).JSON(fiber.Map{"erro": "Termo é obrigatório"})
	}

	// Criar registro de busca
	search := &Search{
		Termo:     req.Termo,
		Categoria: req.Categoria,
		Status:    "em_andamento",
	}

	if err := db.CreateSearch(search); err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	// Aqui você integraria com o Google Maps Scraper
	// Por enquanto, retornamos o ID da busca
	go executeSearch(search)

	return c.Status(201).JSON(search)
}

func handleGetSearch(c fiber.Ctx) error {
	id := c.Params("id")
	search, err := db.GetSearch(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"erro": "Busca não encontrada"})
	}

	return c.JSON(search)
}

func handleListSearches(c fiber.Ctx) error {
	searches, err := db.ListSearches()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.JSON(searches)
}

// ===== CAMPAIGNS =====

func handleCreateCampaign(c fiber.Ctx) error {
	campaign := new(Campaign)
	if err := c.BindJSON(campaign); err != nil {
		return c.Status(400).JSON(fiber.Map{"erro": err.Error()})
	}

	if campaign.Nome == "" {
		return c.Status(400).JSON(fiber.Map{"erro": "Nome da campanha é obrigatório"})
	}

	campaign.DataInicio = time.Now()
	if err := db.CreateCampaign(campaign); err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.Status(201).JSON(campaign)
}

func handleListCampaigns(c fiber.Ctx) error {
	campaigns, err := db.ListCampaigns()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.JSON(campaigns)
}

func handleGetCampaign(c fiber.Ctx) error {
	id := c.Params("id")
	campaign, err := db.GetCampaign(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"erro": "Campanha não encontrada"})
	}

	return c.JSON(campaign)
}

// ===== EXPORT =====

func handleExportCSV(c fiber.Ctx) error {
	filtros := &FiltrosLead{
		Status:   c.Query("status", ""),
		Qualidade: c.Query("qualidade", ""),
		Perfil:   c.Query("perfil", ""),
		Categoria: c.Query("categoria", ""),
		Limit:    999999,
		Offset:   0,
	}

	leads, err := db.ListLeads(filtros)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	// Criar arquivo CSV
	filename := fmt.Sprintf("leads_%d.csv", time.Now().Unix())
	filepath := fmt.Sprintf("./data/%s", filename)

	file, err := os.Create(filepath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Headers
	headers := []string{"ID", "Nome", "Categoria", "Telefone", "Email", "Website", 
		"Endereço", "Rating", "Qualidade", "Perfil", "Status", "Data Criação"}
	writer.Write(headers)

	// Dados
	for _, lead := range leads {
		row := []string{
			lead.ID,
			lead.Nome,
			lead.Categoria,
			lead.Telefone,
			lead.Email,
			lead.Website,
			lead.Endereco,
			fmt.Sprintf("%.1f", lead.Rating),
			lead.Qualidade,
			lead.Perfil,
			lead.Status,
			lead.CriadoEm.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}

	// Enviar arquivo
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.SendFile(filepath)
}

func handleExportJSON(c fiber.Ctx) error {
	filtros := &FiltrosLead{
		Status:   c.Query("status", ""),
		Qualidade: c.Query("qualidade", ""),
		Perfil:   c.Query("perfil", ""),
		Categoria: c.Query("categoria", ""),
		Limit:    999999,
		Offset:   0,
	}

	leads, err := db.ListLeads(filtros)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	// Criar objeto com contatos inclusos
	export := fiber.Map{
		"timestamp": time.Now(),
		"total":     len(leads),
		"contato": fiber.Map{
			"email": config.ContactEmail,
			"phone": config.ContactPhone,
		},
		"leads": leads,
	}

	filename := fmt.Sprintf("leads_%d.json", time.Now().Unix())
	filepath := fmt.Sprintf("./data/%s", filename)

	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.SendFile(filepath)
}

// ===== STATS =====

func handleStats(c fiber.Ctx) error {
	stats, err := db.GetStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"erro": err.Error()})
	}

	return c.JSON(stats)
}

// ===== HELPERS =====

func getIntQuery(c fiber.Ctx, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return num
}

// Placeholder para integração com Google Maps Scraper
func executeSearch(search *Search) {
	// Aqui você integraria com o scraper real
	// Por enquanto, apenas simulamos uma busca
	search.Status = "concluido"
	search.TotalResultados = 0
	search.LeadsQualificados = 0
	db.CreateSearch(search)
}
