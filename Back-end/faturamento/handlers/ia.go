package handlers

import (
	"faturamento_service/database"
	"faturamento_service/models"
	"faturamento_service/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetIAInsights(c *gin.Context) {
	stockClient := services.NewStockClient()
	
	// 1. Buscar saldos atuais no serviço de estoque
	products, err := stockClient.GetProducts()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Nao foi possivel obter dados de estoque para a IA: " + err.Error()})
		return
	}

	insights := []models.IAInsight{}

	// 2. Analise de Estoque (Velocidade de Saida)
	// Vamos ver quanto cada produto vendeu nos ultimos 30 dias
	rows, err := database.DB.Query(`
		SELECT i.id, i.codigo, i.descricao, SUM(it.quantidade) as total_vendido
		FROM itens_fatura it
		JOIN faturas f ON it.fatura_id = f.id
		JOIN itens i ON it.item_id = i.id
		WHERE f.status = 'FECHADA' 
		AND f.data_criacao > NOW() - INTERVAL '30 days'
		GROUP BY i.id, i.codigo, i.descricao
	`)
	
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var codigo, descricao string
			var totalVendido int
			rows.Scan(&id, &codigo, &descricao, &totalVendido)

			// Velocidade media por dia (30 dias)
			velocity := float64(totalVendido) / 30.0
			
			// Encontra o saldo atual deste produto
			var currentSaldo int
			for _, p := range products {
				if p.Codigo == codigo {
					currentSaldo = p.Saldo
					break
				}
			}

			// Se a velocidade for alta e o saldo baixo (ex: dura menos de 7 dias)
			if velocity > 0 {
				diasRestantes := float64(currentSaldo) / velocity
				if diasRestantes < 7 {
					prioridade := "MEDIA"
					if diasRestantes < 3 { prioridade = "ALTA" }

					insights = append(insights, models.IAInsight{
						ID:           fmt.Sprintf("EST-%s", codigo),
						Tipo:         "ESTOQUE",
						Prioridade:   prioridade,
						Mensagem:     fmt.Sprintf("O produto %s esta vendendo rapido (%.1f un/dia).", descricao, velocity),
						AcaoSugerida: fmt.Sprintf("Seu estoque de %d unidades dura apenas %.0f dias. Compre mais agora!", currentSaldo, diasRestantes),
					})
				}
			}
		}
	}

	// 3. Analise de Clientes (Churn / Fidelidade)
	// Clientes que nao compram ha mais de 15 dias
	clientRows, err := database.DB.Query(`
		SELECT c.id, c.nome, MAX(f.data_criacao) as ultima_compra
		FROM clientes c
		JOIN faturas f ON f.cliente_id = c.id
		GROUP BY c.id, c.nome
		HAVING MAX(f.data_criacao) < NOW() - INTERVAL '15 days'
		ORDER BY ultima_compra ASC
		LIMIT 5
	`)

	if err == nil {
		defer clientRows.Close()
		for clientRows.Next() {
			var id int
			var nome string
			var ultimaCompra time.Time
			clientRows.Scan(&id, &nome, &ultimaCompra)

			diasSemComprar := int(time.Since(ultimaCompra).Hours() / 24)

			insights = append(insights, models.IAInsight{
				ID:           fmt.Sprintf("CLI-%d", id),
				Tipo:         "VENDA",
				Prioridade:   "BAIXA",
				Mensagem:     fmt.Sprintf("O cliente %s nao realiza compras ha %d dias.", nome, diasSemComprar),
				AcaoSugerida: "Que tal enviar um cupom de desconto para reativa-lo?",
			})
		}
	}

	// Caso nao tenha nada
	if len(insights) == 0 {
		insights = append(insights, models.IAInsight{
			ID: "ZERO",
			Tipo: "ESTOQUE",
			Prioridade: "BAIXA",
			Mensagem: "Tudo sob controle no seu estoque!",
			AcaoSugerida: "Continue monitorando suas vendas.",
		})
	}

	c.JSON(http.StatusOK, insights)
}
