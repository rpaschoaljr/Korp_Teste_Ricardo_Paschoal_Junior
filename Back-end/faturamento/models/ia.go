package models

type IAInsight struct {
	ID           string `json:"id"`
	Tipo         string `json:"tipo"` // 'ESTOQUE' | 'VENDA'
	Mensagem     string `json:"mensagem"`
	Prioridade   string `json:"prioridade"` // 'ALTA' | 'MEDIA' | 'BAIXA'
	AcaoSugerida string `json:"acao_sugerida"`
}
