package models

type InvoiceItem struct {
	ID            int     `json:"id"`
	Descricao     string  `json:"descricao"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Subtotal      float64 `json:"subtotal"`
}

type Client struct {
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Endereco string `json:"endereco"`
	CPF      string `json:"cpf"`
	CNPJ     string `json:"cnpj"`
}

type Invoice struct {
	ID         int           `json:"id"`
	Cliente    Client        `json:"cliente"`
	ValorTotal float64       `json:"valor_total"`
	Data       string        `json:"data"`
	Itens      []InvoiceItem `json:"itens"`
}
