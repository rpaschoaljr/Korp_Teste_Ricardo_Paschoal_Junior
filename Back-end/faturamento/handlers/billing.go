package handlers

import (
	"encoding/json"
	"faturamento_service/database"
	"faturamento_service/models"
	"faturamento_service/services"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateFatura(c *gin.Context) {
	clientHealth := services.NewClientServiceHealth()
	if err := clientHealth.CheckHealth(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"service": "CLIENTES",
			"error":   "Nao e possivel abrir faturas: o servico de CLIENTES esta offline.",
		})
		return
	}

	var f models.Fatura
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados invalidos: " + err.Error()})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transacao"})
		return
	}
	defer tx.Rollback()

	var faturaID int
	query := "INSERT INTO faturas (cliente_id, status, valor_total) VALUES ($1, 'ABERTA', 0) RETURNING id"
	err = tx.QueryRow(query, f.ClienteID).Scan(&faturaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar fatura"})
		return
	}

	var valorTotal float64
	for _, item := range f.Itens {
		subtotal := float64(item.Quantidade) * item.PrecoUnitario
		valorTotal += subtotal

		_, err = tx.Exec("INSERT INTO itens_fatura (fatura_id, item_id, quantidade, preco_unitario, subtotal) VALUES ($1, $2, $3, $4, $5)",
			faturaID, item.ItemID, item.Quantidade, item.PrecoUnitario, subtotal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao inserir itens da fatura"})
			return
		}
	}

	_, err = tx.Exec("UPDATE faturas SET valor_total = $1 WHERE id = $2", valorTotal, faturaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar total da fatura"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao confirmar fatura"})
		return
	}

	f.ID = faturaID
	f.Status = "ABERTA"
	f.ValorTotal = valorTotal
	c.JSON(http.StatusCreated, f)
}

func PrintFatura(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	fmt.Printf("[DEBUG-FAT] Iniciando impressao da fatura ID: %d\n", id)

	clientHealth := services.NewClientServiceHealth()
	stockClient := services.NewStockClient()
	printClient := services.NewPrintClient()

	if err := clientHealth.CheckHealth(); err != nil {
		fmt.Printf("[DEBUG-FAT] Erro Health Clientes: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "CLIENTES", "error": "Servico de CLIENTES indisponivel."})
		return
	}
	if err := printClient.CheckHealth(); err != nil {
		fmt.Printf("[DEBUG-FAT] Erro Health Impressao: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "IMPRESSAO", "error": "Servico de IMPRESSAO indisponivel."})
		return
	}

	var fatura models.Fatura
	err := database.DB.QueryRow("SELECT id, cliente_id, status, valor_total, data_criacao FROM faturas WHERE id = $1", id).
		Scan(&fatura.ID, &fatura.ClienteID, &fatura.Status, &fatura.ValorTotal, &fatura.DataCriacao)
	if err != nil {
		fmt.Printf("[DEBUG-FAT] Erro ao buscar fatura no banco: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Fatura nao encontrada"})
		return
	}

	// Busca detalhes do Cliente
	clientURL := os.Getenv("CLIENTES_URL")
	if clientURL == "" { clientURL = "http://clientes_api:8083" }
	fullClientURL := fmt.Sprintf("%s/clientes/%d", clientURL, fatura.ClienteID)
	fmt.Printf("[DEBUG-FAT] Chamando Cliente: %s\n", fullClientURL)

	resp, err := http.Get(fullClientURL)
	var clientInfo models.ClientInfo
	if err != nil || resp.StatusCode != http.StatusOK {
		status := "ERRO CONEXAO"
		if resp != nil { status = strconv.Itoa(resp.StatusCode) }
		fmt.Printf("[DEBUG-FAT] Erro ao obter detalhes do cliente: %s\n", status)
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "CLIENTES", "error": "Nao foi possivel obter os detalhes cadastrais do cliente para a nota."})
		return
	}
	
	err = json.NewDecoder(resp.Body).Decode(&clientInfo)
	if err != nil {
		fmt.Printf("[DEBUG-FAT] Erro ao decodificar Cliente: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro de processamento nos dados do cliente"})
		resp.Body.Close()
		return
	}
	resp.Body.Close()
	fmt.Printf("[DEBUG-FAT] Cliente recuperado: %+v\n", clientInfo)

	// Busca Itens
	rows, err := database.DB.Query(`
		SELECT i.codigo, i.descricao, it.quantidade, it.preco_unitario, it.subtotal 
		FROM itens_fatura it 
		JOIN itens i ON it.item_id = i.id 
		WHERE it.fatura_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar itens"})
		return
	}
	defer rows.Close()

	var itens []models.ItemFatura
	var stockUpdates []models.StockUpdate
	for rows.Next() {
		var item models.ItemFatura
		rows.Scan(&item.CodigoProduto, &item.Descricao, &item.Quantidade, &item.PrecoUnitario, &item.Subtotal)
		itens = append(itens, item)
		stockUpdates = append(stockUpdates, models.StockUpdate{Codigo: item.CodigoProduto, Quantidade: item.Quantidade})
	}

	printData := models.PrintData{
		ID:         fatura.ID,
		Cliente:    clientInfo,
		ValorTotal: fatura.ValorTotal,
		Data:       fatura.DataCriacao.Format("02/01/2006 15:04"),
		Itens:      itens,
	}
	fmt.Printf("[DEBUG-FAT] Payload para Impressao: ID=%d, NomeCli='%s', Itens=%d\n", printData.ID, printData.Cliente.Nome, len(printData.Itens))

	pdfBytes, err := printClient.GeneratePDF(printData)
	if err != nil {
		fmt.Printf("[DEBUG-FAT] Erro ao chamar Impressao: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "IMPRESSAO", "error": "Falha ao gerar documento: " + err.Error()})
		return
	}

	if err := stockClient.ReduceStock(stockUpdates); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "ESTOQUE", "error": "Falha na baixa de estoque: " + err.Error()})
		return
	}

	_, err = database.DB.Exec("UPDATE faturas SET status = 'FECHADA' WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fechar fatura no banco local"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func ListFaturas(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, cliente_id, status, valor_total, data_criacao FROM faturas ORDER BY data_criacao DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar faturas"})
		return
	}
	defer rows.Close()

	var faturas []models.Fatura
	for rows.Next() {
		var f models.Fatura
		rows.Scan(&f.ID, &f.ClienteID, &f.Status, &f.ValorTotal, &f.DataCriacao)
		faturas = append(faturas, f)
	}

	if faturas == nil { faturas = []models.Fatura{} }
	c.JSON(http.StatusOK, faturas)
}
