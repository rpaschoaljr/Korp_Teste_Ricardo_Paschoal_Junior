package handlers

import (
	"database/sql"
	"estoque_service/database"
	"estoque_service/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, codigo, descricao, saldo, preco_base FROM itens ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar produtos"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Codigo, &p.Descricao, &p.Saldo, &p.PrecoBase); err != nil {
			continue
		}
		products = append(products, p)
	}

	if products == nil {
		products = []models.Product{}
	}

	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// 1. Normalização rigorosa da descrição
	cleanDesc := strings.TrimSpace(strings.ToUpper(p.Descricao))
	if cleanDesc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A descrição do produto é obrigatória"})
		return
	}
	p.Descricao = cleanDesc

	// 2. Verificação de Duplicidade pela DESCRIÇÃO (Independente do código)
	var existingID int
	err := database.DB.QueryRow("SELECT id FROM itens WHERE UPPER(TRIM(descricao)) = $1", cleanDesc).Scan(&existingID)
	
	if err == nil {
		// Se não deu erro, significa que encontrou um registro
		c.JSON(http.StatusConflict, gin.H{"error": "Já existe um produto cadastrado com o nome: " + cleanDesc})
		return
	} else if err != sql.ErrNoRows {
		// Erro técnico de banco
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao validar unicidade: " + err.Error()})
		return
	}

	// 3. Inserção final
	query := "INSERT INTO itens (codigo, descricao, saldo, preco_base) VALUES ($1, $2, $3, $4) RETURNING id"
	err = database.DB.QueryRow(query, p.Codigo, p.Descricao, p.Saldo, p.PrecoBase).Scan(&p.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar produto: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func UpdateStock(c *gin.Context) {
	var updates []models.StockUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lista de baixa inválida"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação"})
		return
	}
	defer tx.Rollback()

	for _, update := range updates {
		var currentSaldo int
		err := tx.QueryRow("SELECT saldo FROM itens WHERE codigo = $1 FOR UPDATE", update.Codigo).Scan(&currentSaldo)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado: " + update.Codigo})
			return
		}

		if currentSaldo < update.Quantidade {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo insuficiente para o produto: " + update.Codigo})
			return
		}

		_, err = tx.Exec("UPDATE itens SET saldo = saldo - $1 WHERE codigo = $2", update.Quantidade, update.Codigo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar saldo"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao confirmar transação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estoque atualizado com sucesso"})
}
