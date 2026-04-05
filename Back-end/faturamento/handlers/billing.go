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
	"strings"

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
		errStr := err.Error()
		if strings.Contains(errStr, "required") {
			if strings.Contains(errStr, "ClienteID") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário selecionar um CLIENTE para abrir a nota."})
				return
			}
			if strings.Contains(errStr, "Itens") || strings.Contains(errStr, "ItemID") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "A nota deve conter pelo menos um PRODUTO com quantidade válida."})
				return
			}
		}
		
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: Verifique se todos os campos da nota foram preenchidos."})
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

	// --- INICIO DA TRAVA DE CONCORRENCIA / IDEMPOTENCIA ---
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transacao de seguranca"})
		return
	}
	defer tx.Rollback()

	var currentStatus string
	// O FOR UPDATE bloqueia a linha no banco ate o Commit/Rollback, impedindo que outra 
	// requisicao simultanea leia o status 'ABERTA' ao mesmo tempo.
	err = tx.QueryRow("SELECT status FROM faturas WHERE id = $1 FOR UPDATE", id).Scan(&currentStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fatura nao encontrada"})
		return
	}

	if currentStatus != "ABERTA" {
		c.JSON(http.StatusConflict, gin.H{"error": "Esta nota ja foi processada ou esta fechada. Operacao cancelada para evitar duplicidade."})
		return
	}
	// --- FIM DA TRAVA ---

	clientHealth := services.NewClientServiceHealth()
	stockClient := services.NewStockClient()
	printClient := services.NewPrintClient()

	// ... verificacoes de health permanecem ...
	if err := clientHealth.CheckHealth(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "CLIENTES", "error": "Servico de CLIENTES indisponivel."})
		return
	}
	if err := printClient.CheckHealth(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "IMPRESSAO", "error": "Servico de IMPRESSAO indisponivel."})
		return
	}

	var fatura models.Fatura
	// Recarregamos os dados (ja estamos dentro da transacao segura)
	err = tx.QueryRow("SELECT id, cliente_id, status, valor_total, data_criacao FROM faturas WHERE id = $1", id).
		Scan(&fatura.ID, &fatura.ClienteID, &fatura.Status, &fatura.ValorTotal, &fatura.DataCriacao)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar dados da fatura"})
		return
	}

	// Busca detalhes do Cliente (externo)
	clientURL := os.Getenv("CLIENTES_URL")
	if clientURL == "" { clientURL = "http://clientes_api:8083" }
	resp, err := http.Get(fmt.Sprintf("%s/clientes/%d", clientURL, fatura.ClienteID))
	var clientInfo models.ClientInfo
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "CLIENTES", "error": "Nao foi possivel obter detalhes do cliente."})
		return
	}
	json.NewDecoder(resp.Body).Decode(&clientInfo)
	resp.Body.Close()

	// Busca Itens (local)
	rows, err := tx.Query(`
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

	// Chamadas externas (PDF e Estoque)
	pdfBytes, err := printClient.GeneratePDF(printData)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"service": "IMPRESSAO", "error": "Falha ao gerar PDF: " + err.Error()})
		return
	}

	if err := stockClient.ReduceStock(stockUpdates); err != nil {
		// Se o erro contém "estoque recusou a baixa", tratamos como erro de negócio (422)
		// caso contrário, tratamos como erro de serviço (503)
		status := http.StatusServiceUnavailable
		if strings.Contains(err.Error(), "estoque recusou a baixa") {
			status = http.StatusUnprocessableEntity
		}

		c.JSON(status, gin.H{
			"service": "ESTOQUE",
			"error":   "Falha na baixa de estoque: " + err.Error(),
		})
		return
	}

	// Se chegou aqui, deu tudo certo. Fechamos a nota.
	_, err = tx.Exec("UPDATE faturas SET status = 'FECHADA' WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar nota no banco"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao confirmar transacao final"})
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

		// Buscar itens para cada fatura
		itemRows, err := database.DB.Query(`
			SELECT it.id, it.item_id, it.quantidade, it.preco_unitario, it.subtotal, i.codigo, i.descricao
			FROM itens_fatura it
			JOIN itens i ON it.item_id = i.id
			WHERE it.fatura_id = $1`, f.ID)
		
		if err == nil {
			var itens []models.ItemFatura
			for itemRows.Next() {
				var item models.ItemFatura
				itemRows.Scan(&item.ID, &item.ItemID, &item.Quantidade, &item.PrecoUnitario, &item.Subtotal, &item.CodigoProduto, &item.Descricao)
				itens = append(itens, item)
			}
			f.Itens = itens
			itemRows.Close()
		}

		faturas = append(faturas, f)
	}

	if faturas == nil { faturas = []models.Fatura{} }
	c.JSON(http.StatusOK, faturas)
}
