package models

type Product struct {
	ID        int     `json:"id"`
	Codigo    string  `json:"codigo"`
	Descricao string  `json:"descricao" binding:"required"`
	Saldo     int     `json:"saldo" binding:"required,min=0"`
	PrecoBase float64 `json:"preco_base" binding:"required"`
}

type StockUpdate struct {
	Codigo     string `json:"codigo" binding:"required"`
	Quantidade int    `json:"quantidade" binding:"required,min=1"`
}
