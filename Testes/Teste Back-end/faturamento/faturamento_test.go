package faturamento_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

type ItemFatura struct {
	ItemID        int     `json:"item_id"`
	CodigoProduto string  `json:"codigo_produto"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Subtotal      float64 `json:"subtotal"`
}

type Fatura struct {
	ID         int          `json:"id"`
	ClienteID  int          `json:"cliente_id"`
	Status     string       `json:"status"`
	ValorTotal float64      `json:"valor_total"`
	Itens      []ItemFatura `json:"itens"`
}

type Item struct {
	ID     int    `json:"id"`
	Codigo string `json:"codigo"`
	Saldo  int    `json:"saldo"`
}

var (
	dbFaturamentoConn = os.Getenv("DB_FATURAMENTO_CONN")
	faturamentoURL    = os.Getenv("FATURAMENTO_URL")
	estoqueURL        = os.Getenv("ESTOQUE_URL")
)

func getDBConnection(t *testing.T) *sql.DB {
	if dbFaturamentoConn == "" {
		dbFaturamentoConn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", dbFaturamentoConn)
	if err != nil {
		t.Fatalf("Erro ao conectar no banco: %v", err)
	}
	return db
}

func TestServicoFaturamento(t *testing.T) {
	if faturamentoURL == "" {
		faturamentoURL = "http://localhost:8082"
	}
	if estoqueURL == "" {
		estoqueURL = "http://localhost:8081"
	}

	db := getDBConnection(t)
	defer db.Close()

	t.Run("Fluxo Completo: Criar -> Imprimir -> Baixar Estoque", func(t *testing.T) {
		ts := time.Now().UnixNano()
		// 1. Criar um produto dinâmico para o teste
		novoItem := map[string]interface{}{
			"codigo":     fmt.Sprintf("TEST-OK-%d", ts),
			"descricao":  fmt.Sprintf("PRODUTO TESTE OK %d", ts),
			"saldo":      10,
			"preco_base": 10.0,
		}
		bodyItem, _ := json.Marshal(novoItem)
		respItem, err := http.Post(estoqueURL+"/produtos", "application/json", bytes.NewBuffer(bodyItem))
		if err != nil {
			t.Fatalf("Erro ao criar produto para teste: %v", err)
		}
		defer respItem.Body.Close()
		
		var itemCriado Item
		json.NewDecoder(respItem.Body).Decode(&itemCriado)

		// 2. Criar Fatura usando o ID e Código gerados
		novaFatura := Fatura{
			ClienteID: 1,
			Itens: []ItemFatura{
				{ItemID: itemCriado.ID, CodigoProduto: itemCriado.Codigo, Quantidade: 3, PrecoUnitario: 20.0},
			},
		}

		bodyFatura, _ := json.Marshal(novaFatura)
		respFatura, err := http.Post(faturamentoURL+"/faturas", "application/json", bytes.NewBuffer(bodyFatura))
		if err != nil {
			t.Fatalf("Erro ao criar fatura: %v", err)
		}
		defer respFatura.Body.Close()

		var faturaCriada Fatura
		json.NewDecoder(respFatura.Body).Decode(&faturaCriada)

		// 3. Imprimir (Efetivar baixa)
		respImp, err := http.Post(fmt.Sprintf("%s/faturas/%d/imprimir", faturamentoURL, faturaCriada.ID), "application/json", nil)
		if err != nil {
			t.Fatalf("Erro ao imprimir fatura: %v", err)
		}
		defer respImp.Body.Close()

		if respImp.StatusCode != http.StatusOK {
			t.Fatalf("Erro ao imprimir fatura, status: %d", respImp.StatusCode)
		}

		// 4. Validar se estoque foi reduzido para exatamente 7 (10 - 3)
		var saldoFinal int
		err = db.QueryRow("SELECT saldo FROM itens WHERE id = $1", itemCriado.ID).Scan(&saldoFinal)
		if err != nil {
			t.Fatalf("Erro ao consultar saldo final: %v", err)
		}

		if saldoFinal != 7 {
			t.Errorf("Estoque não foi reduzido corretamente. Esperado: 7, Recebido: %d", saldoFinal)
		}
	})

	t.Run("Fluxo Erro: Saldo Insuficiente ao Fechar Nota", func(t *testing.T) {
		ts := time.Now().UnixNano()
		// 1. Criar um produto com saldo 2
		novoItem := map[string]interface{}{
			"codigo":     fmt.Sprintf("TEST-FAIL-%d", ts),
			"descricao":  fmt.Sprintf("PRODUTO TESTE FAIL %d", ts),
			"saldo":      2,
			"preco_base": 10.0,
		}
		bodyItem, _ := json.Marshal(novoItem)
		respItem, err := http.Post(estoqueURL+"/produtos", "application/json", bytes.NewBuffer(bodyItem))
		if err != nil {
			t.Fatalf("Erro ao criar produto para teste: %v", err)
		}
		defer respItem.Body.Close()
		
		var itemCriado Item
		json.NewDecoder(respItem.Body).Decode(&itemCriado)

		// 2. Criar Fatura pedindo 5 unidades (mais que o saldo 2)
		novaFatura := Fatura{
			ClienteID: 1,
			Itens: []ItemFatura{
				{ItemID: itemCriado.ID, CodigoProduto: itemCriado.Codigo, Quantidade: 5, PrecoUnitario: 20.0},
			},
		}

		bodyFatura, _ := json.Marshal(novaFatura)
		respFatura, err := http.Post(faturamentoURL+"/faturas", "application/json", bytes.NewBuffer(bodyFatura))
		if err != nil {
			t.Fatalf("Erro ao criar fatura: %v", err)
		}
		defer respFatura.Body.Close()

		var faturaCriada Fatura
		json.NewDecoder(respFatura.Body).Decode(&faturaCriada)

		// 3. Tentar Imprimir (Deve falhar)
		respImp, err := http.Post(fmt.Sprintf("%s/faturas/%d/imprimir", faturamentoURL, faturaCriada.ID), "application/json", nil)
		if err != nil {
			t.Fatalf("Erro ao chamar endpoint de impressão: %v", err)
		}
		defer respImp.Body.Close()

		if respImp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Esperado status 422 (Unprocessable Entity) por saldo insuficiente, recebido: %d", respImp.StatusCode)
		}

		// 4. Validar se o status da nota continua 'ABERTA' no banco
		var status string
		err = db.QueryRow("SELECT status FROM faturas WHERE id = $1", faturaCriada.ID).Scan(&status)
		if err != nil {
			t.Fatalf("Erro ao consultar status da fatura: %v", err)
		}

		if status != "ABERTA" {
			t.Errorf("A fatura deveria continuar ABERTA após falha no estoque, mas está: %s", status)
		}
	})
}
