package models

import "time"

type Fatura struct {
	ID          int           `json:"id"`
	ClienteID   int           `json:"cliente_id" binding:"required"`
	Status      string        `json:"status"`
	ValorTotal  float64       `json:"valor_total"`
	DataCriacao time.Time     `json:"data_criacao"`
	Itens       []ItemFatura `json:"itens" binding:"required,min=1"`
}

type ItemFatura struct {
	ID            int     `json:"id"`
	FaturaID      int     `json:"fatura_id"`
	ItemID        int     `json:"item_id" binding:"required"`
	CodigoProduto string  `json:"codigo_produto"`
	Descricao     string  `json:"descricao"` // Adicionado para o PDF
	Quantidade    int     `json:"quantidade" binding:"required,min=1"`
	PrecoUnitario float64 `json:"preco_unitario" binding:"required"`
	Subtotal      float64 `json:"subtotal"`
}

type StockUpdate struct {
	Codigo     string `json:"codigo"`
	Quantidade int    `json:"quantidade"`
}

// Modelos para integração com Clientes e Impressão
type ClientInfo struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Endereco string `json:"endereco"`
	CPF      string `json:"cpf"`
	CNPJ     string `json:"cnpj"`
}

type PrintData struct {
	ID         int          `json:"id"`
	Cliente    ClientInfo   `json:"cliente"`
	ValorTotal float64      `json:"valor_total"`
	Data       string       `json:"data"`
	Itens      []ItemFatura `json:"itens"`
}
