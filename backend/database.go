package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

// Init cria as tabelas
func (db *Database) Init() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS leads (
			id TEXT PRIMARY KEY,
			scraper_id TEXT UNIQUE,
			nome TEXT NOT NULL,
			categoria TEXT,
			telefone TEXT,
			email TEXT,
			website TEXT,
			endereco TEXT,
			rating REAL,
			reviews_count INTEGER,
			latitude REAL,
			longitude REAL,
			qualidade TEXT DEFAULT 'medium',
			perfil TEXT,
			status TEXT DEFAULT 'novo',
			notas TEXT,
			campanha_id TEXT,
			criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS buscas (
			id TEXT PRIMARY KEY,
			termo TEXT NOT NULL,
			categoria TEXT,
			data_execucao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			total_resultados INTEGER,
			leads_qualificados INTEGER,
			arquivo_raw TEXT,
			status TEXT DEFAULT 'em_andamento',
			erro TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS campanhas (
			id TEXT PRIMARY KEY,
			nome TEXT NOT NULL,
			descricao TEXT,
			tipo TEXT,
			data_inicio TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			data_fim TIMESTAMP,
			total_leads INTEGER DEFAULT 0,
			contatos_feitos INTEGER DEFAULT 0,
			conversoes INTEGER DEFAULT 0,
			arquivo_export TEXT,
			criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		// Índices para melhor performance
		`CREATE INDEX IF NOT EXISTS idx_leads_status ON leads(status)`,
		`CREATE INDEX IF NOT EXISTS idx_leads_perfil ON leads(perfil)`,
		`CREATE INDEX IF NOT EXISTS idx_leads_campanha ON leads(campanha_id)`,
		`CREATE INDEX IF NOT EXISTS idx_leads_qualidade ON leads(qualidade)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("erro ao criar tabela: %v", err)
		}
	}

	return nil
}

// ===== LEADS =====

func (db *Database) CreateLead(lead *Lead) error {
	lead.ID = uuid.Must(uuid.NewV4()).String()
	lead.CriadoEm = time.Now()
	lead.AtualizadoEm = time.Now()

	query := `INSERT INTO leads (id, scraper_id, nome, categoria, telefone, email, website, 
	endereco, rating, reviews_count, latitude, longitude, qualidade, perfil, status, notas, campanha_id, criado_em, atualizado_em)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query, lead.ID, lead.ScraperID, lead.Nome, lead.Categoria,
		lead.Telefone, lead.Email, lead.Website, lead.Endereco, lead.Rating, lead.ReviewsCount,
		lead.Latitude, lead.Longitude, lead.Qualidade, lead.Perfil, lead.Status, lead.Notas,
		lead.CampanhaID, lead.CriadoEm, lead.AtualizadoEm)

	return err
}

func (db *Database) GetLead(id string) (*Lead, error) {
	lead := &Lead{}
	query := `SELECT id, scraper_id, nome, categoria, telefone, email, website, endereco, 
	rating, reviews_count, latitude, longitude, qualidade, perfil, status, notas, campanha_id, criado_em, atualizado_em 
	FROM leads WHERE id = ?`

	err := db.conn.QueryRow(query, id).Scan(&lead.ID, &lead.ScraperID, &lead.Nome, &lead.Categoria,
		&lead.Telefone, &lead.Email, &lead.Website, &lead.Endereco, &lead.Rating, &lead.ReviewsCount,
		&lead.Latitude, &lead.Longitude, &lead.Qualidade, &lead.Perfil, &lead.Status, &lead.Notas,
		&lead.CampanhaID, &lead.CriadoEm, &lead.AtualizadoEm)

	return lead, err
}

func (db *Database) ListLeads(filtros *FiltrosLead) ([]Lead, error) {
	query := `SELECT id, scraper_id, nome, categoria, telefone, email, website, endereco, 
	rating, reviews_count, latitude, longitude, qualidade, perfil, status, notas, campanha_id, criado_em, atualizado_em 
	FROM leads WHERE 1=1`

	var args []interface{}

	if filtros.Status != "" {
		query += " AND status = ?"
		args = append(args, filtros.Status)
	}
	if filtros.Qualidade != "" {
		query += " AND qualidade = ?"
		args = append(args, filtros.Qualidade)
	}
	if filtros.Perfil != "" {
		query += " AND perfil = ?"
		args = append(args, filtros.Perfil)
	}
	if filtros.Categoria != "" {
		query += " AND categoria = ?"
		args = append(args, filtros.Categoria)
	}
	if filtros.Busca != "" {
		query += " AND (nome LIKE ? OR email LIKE ? OR telefone LIKE ?)"
		busca := "%" + filtros.Busca + "%"
		args = append(args, busca, busca, busca)
	}

	query += " ORDER BY atualizado_em DESC"

	if filtros.Limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, filtros.Limit, filtros.Offset)
	}

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		lead := Lead{}
		err := rows.Scan(&lead.ID, &lead.ScraperID, &lead.Nome, &lead.Categoria,
			&lead.Telefone, &lead.Email, &lead.Website, &lead.Endereco, &lead.Rating, &lead.ReviewsCount,
			&lead.Latitude, &lead.Longitude, &lead.Qualidade, &lead.Perfil, &lead.Status, &lead.Notas,
			&lead.CampanhaID, &lead.CriadoEm, &lead.AtualizadoEm)
		if err != nil {
			return nil, err
		}
		leads = append(leads, lead)
	}

	return leads, rows.Err()
}

func (db *Database) UpdateLead(lead *Lead) error {
	lead.AtualizadoEm = time.Now()

	query := `UPDATE leads SET nome=?, categoria=?, telefone=?, email=?, website=?, 
	endereco=?, rating=?, reviews_count=?, qualidade=?, perfil=?, status=?, notas=?, 
	campanha_id=?, atualizado_em=? WHERE id=?`

	_, err := db.conn.Exec(query, lead.Nome, lead.Categoria, lead.Telefone, lead.Email,
		lead.Website, lead.Endereco, lead.Rating, lead.ReviewsCount, lead.Qualidade,
		lead.Perfil, lead.Status, lead.Notas, lead.CampanhaID, lead.AtualizadoEm, lead.ID)

	return err
}

func (db *Database) DeleteLead(id string) error {
	_, err := db.conn.Exec("DELETE FROM leads WHERE id = ?", id)
	return err
}

// ===== SEARCHES =====

func (db *Database) CreateSearch(search *Search) error {
	search.ID = uuid.Must(uuid.NewV4()).String()

	query := `INSERT INTO buscas (id, termo, categoria, data_execucao, total_resultados, 
	leads_qualificados, arquivo_raw, status, erro)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query, search.ID, search.Termo, search.Categoria, search.DataExecucao,
		search.TotalResultados, search.LeadsQualificados, search.ArquivoRaw, search.Status, search.Erro)

	return err
}

func (db *Database) GetSearch(id string) (*Search, error) {
	search := &Search{}
	query := `SELECT id, termo, categoria, data_execucao, total_resultados, leads_qualificados, 
	arquivo_raw, status, erro FROM buscas WHERE id = ?`

	err := db.conn.QueryRow(query, id).Scan(&search.ID, &search.Termo, &search.Categoria,
		&search.DataExecucao, &search.TotalResultados, &search.LeadsQualificados,
		&search.ArquivoRaw, &search.Status, &search.Erro)

	return search, err
}

func (db *Database) ListSearches() ([]Search, error) {
	query := `SELECT id, termo, categoria, data_execucao, total_resultados, leads_qualificados, 
	arquivo_raw, status, erro FROM buscas ORDER BY data_execucao DESC LIMIT 50`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searches []Search
	for rows.Next() {
		search := Search{}
		err := rows.Scan(&search.ID, &search.Termo, &search.Categoria, &search.DataExecucao,
			&search.TotalResultados, &search.LeadsQualificados, &search.ArquivoRaw,
			&search.Status, &search.Erro)
		if err != nil {
			return nil, err
		}
		searches = append(searches, search)
	}

	return searches, rows.Err()
}

// ===== CAMPAIGNS =====

func (db *Database) CreateCampaign(campaign *Campaign) error {
	campaign.ID = uuid.Must(uuid.NewV4()).String()
	campaign.CriadoEm = time.Now()

	query := `INSERT INTO campanhas (id, nome, descricao, tipo, data_inicio, data_fim, 
	total_leads, contatos_feitos, conversoes, arquivo_export, criado_em)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query, campaign.ID, campaign.Nome, campaign.Descricao, campaign.Tipo,
		campaign.DataInicio, campaign.DataFim, campaign.TotalLeads, campaign.ContatosFeitos,
		campaign.Conversoes, campaign.ArquivoExport, campaign.CriadoEm)

	return err
}

func (db *Database) ListCampaigns() ([]Campaign, error) {
	query := `SELECT id, nome, descricao, tipo, data_inicio, data_fim, total_leads, 
	contatos_feitos, conversoes, arquivo_export, criado_em FROM campanhas ORDER BY criado_em DESC`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []Campaign
	for rows.Next() {
		campaign := Campaign{}
		err := rows.Scan(&campaign.ID, &campaign.Nome, &campaign.Descricao, &campaign.Tipo,
			&campaign.DataInicio, &campaign.DataFim, &campaign.TotalLeads, &campaign.ContatosFeitos,
			&campaign.Conversoes, &campaign.ArquivoExport, &campaign.CriadoEm)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, rows.Err()
}

func (db *Database) GetCampaign(id string) (*Campaign, error) {
	campaign := &Campaign{}
	query := `SELECT id, nome, descricao, tipo, data_inicio, data_fim, total_leads, 
	contatos_feitos, conversoes, arquivo_export, criado_em FROM campanhas WHERE id = ?`

	err := db.conn.QueryRow(query, id).Scan(&campaign.ID, &campaign.Nome, &campaign.Descricao,
		&campaign.Tipo, &campaign.DataInicio, &campaign.DataFim, &campaign.TotalLeads,
		&campaign.ContatosFeitos, &campaign.Conversoes, &campaign.ArquivoExport, &campaign.CriadoEm)

	return campaign, err
}

// ===== STATS =====

func (db *Database) GetStats() (*Stats, error) {
	stats := &Stats{
		LeadsPorPerfil:    make(map[string]int),
		LeadsPorCategoria: make(map[string]int),
	}

	// Total leads
	db.conn.QueryRow("SELECT COUNT(*) FROM leads").Scan(&stats.TotalLeads)

	// Leads por qualidade
	db.conn.QueryRow("SELECT COUNT(*) FROM leads WHERE qualidade = 'high'").Scan(&stats.LeadsAltos)
	db.conn.QueryRow("SELECT COUNT(*) FROM leads WHERE qualidade = 'medium'").Scan(&stats.LeadsMedios)
	db.conn.QueryRow("SELECT COUNT(*) FROM leads WHERE qualidade = 'low'").Scan(&stats.LeadsBaixos)

	// Total campanhas e buscas
	db.conn.QueryRow("SELECT COUNT(*) FROM campanhas").Scan(&stats.TotalCampanhas)
	db.conn.QueryRow("SELECT COUNT(*) FROM buscas").Scan(&stats.TotalBuscas)

	// Taxa de conversão
	var conversoes, contatosFeitos int
	db.conn.QueryRow("SELECT COALESCE(SUM(conversoes), 0), COALESCE(SUM(contatos_feitos), 0) FROM campanhas").
		Scan(&conversoes, &contatosFeitos)
	if contatosFeitos > 0 {
		stats.TaxaConversao = float64(conversoes) / float64(contatosFeitos) * 100
	}

	// Leads por perfil
	rows, _ := db.conn.Query("SELECT perfil, COUNT(*) FROM leads WHERE perfil IS NOT NULL GROUP BY perfil")
	defer rows.Close()
	for rows.Next() {
		var perfil string
		var count int
		rows.Scan(&perfil, &count)
		stats.LeadsPorPerfil[perfil] = count
	}

	// Leads por categoria
	rows, _ = db.conn.Query("SELECT categoria, COUNT(*) FROM leads WHERE categoria IS NOT NULL GROUP BY categoria")
	defer rows.Close()
	for rows.Next() {
		var categoria string
		var count int
		rows.Scan(&categoria, &count)
		stats.LeadsPorCategoria[categoria] = count
	}

	return stats, nil
}
