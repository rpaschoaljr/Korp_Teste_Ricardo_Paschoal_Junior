package handlers

import (
	"clientes/database"
	"clientes/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func GetClientes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query("SELECT id, nome, telefone, endereco, cpf, cnpj FROM clientes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var c models.Cliente
		var telefone, endereco, cpf, cnpj sql.NullString
		if err := rows.Scan(&c.ID, &c.Nome, &telefone, &endereco, &cpf, &cnpj); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c.Telefone = telefone.String
		c.Endereco = endereco.String
		c.CPF = cpf.String
		c.CNPJ = cnpj.String
		clientes = append(clientes, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func GetClienteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/clientes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var c models.Cliente
	var telefone, endereco, cpf, cnpj sql.NullString
	err = database.DB.QueryRow("SELECT id, nome, telefone, endereco, cpf, cnpj FROM clientes WHERE id = $1", id).
		Scan(&c.ID, &c.Nome, &telefone, &endereco, &cpf, &cnpj)

	if err == sql.ErrNoRows {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Telefone = telefone.String
	c.Endereco = endereco.String
	c.CPF = cpf.String
	c.CNPJ = cnpj.String

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}
