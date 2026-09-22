package main

import "time"

// Lead representa um lead qualificado
type Lead struct {
	ID           string    `json:"id"`
	ScraperID    string    `json:"scraper_id"`
	Nome         string    `json:"nome"`
	Categoria    string    `json:"categoria"`
	Telefone     string    `json:"telefone"`
	Email        string    `json:"email"`
	Website      string    `json:"website"`
	Endereco     string    `json:"endereco"`
	Rating       float64   `json:"rating"`
	ReviewsCount int       `json:"reviews_count"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Qualidade    string    `json:"qualidade"` // high, medium, low
	Perfil       string    `json:"perfil"`   // prime_flow, gerencia, site_vendas
	Status       string    `json:"status"`   // novo, contato, conversao, perdido
	Notas        string    `json:"notas"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
	CampanhaID   *string   `json:"campanha_id"`
}

// Search representa uma busca realizada
type Search struct {
	ID               string    `json:"id"`
	Termo            string    `json:"termo"`
	Categoria        string    `json:"categoria"`
	DataExecucao     time.Time `json:"data_execucao"`
	TotalResultados  int       `json:"total_resultados"`
	LeadsQualificados int      `json:"leads_qualificados"`
	ArquivoRaw       string    `json:"arquivo_raw"`
	Status           string    `json:"status"` // em_andamento, concluido, erro
	Erro             string    `json:"erro"`
}

// Campaign representa uma campanha de prospecção
type Campaign struct {
	ID             string    `json:"id"`
	Nome           string    `json:"nome"`
	Descricao      string    `json:"descricao"`
	Tipo           string    `json:"tipo"` // distribuidora, varejista, ecommerce, etc
	DataInicio     time.Time `json:"data_inicio"`
	DataFim        *time.Time `json:"data_fim"`
	TotalLeads     int       `json:"total_leads"`
	ContatosFeitos int       `json:"contatos_feitos"`
	Conversoes     int       `json:"conversoes"`
	ArquivoExport  string    `json:"arquivo_export"`
	CriadoEm       time.Time `json:"criado_em"`
}

// Stats representa estatísticas gerais
type Stats struct {
	TotalLeads         int     `json:"total_leads"`
	LeadsAltos         int     `json:"leads_altos"`
	LeadsMedios        int     `json:"leads_medios"`
	LeadsBaixos        int     `json:"leads_baixos"`
	TotalCampanhas     int     `json:"total_campanhas"`
	TotalBuscas        int     `json:"total_buscas"`
	TaxaConversao      float64 `json:"taxa_conversao"`
	LeadsPorPerfil     map[string]int `json:"leads_por_perfil"`
	LeadsPorCategoria  map[string]int `json:"leads_por_categoria"`
}

// Filtros para busca avançada
type FiltrosLead struct {
	Status    string `json:"status"`
	Qualidade string `json:"qualidade"`
	Perfil    string `json:"perfil"`
	Categoria string `json:"categoria"`
	Busca     string `json:"busca"` // busca em nome/email/telefone
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}
