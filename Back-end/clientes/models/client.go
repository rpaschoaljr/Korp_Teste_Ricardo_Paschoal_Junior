package models

type Cliente struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Endereco string `json:"endereco"`
	CPF      string `json:"cpf"`
	CNPJ     string `json:"cnpj"`
}
