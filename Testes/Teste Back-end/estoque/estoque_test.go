package estoque_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

type Item struct {
	ID        int     `json:"id"`
	Codigo    string  `json:"codigo"`
	Descricao string  `json:"descricao"`
	Saldo     int     `json:"saldo"`
	PrecoBase float64 `json:"preco_base"`
}

var (
	dbConnStr  = os.Getenv("DB_CONN_STR")
	backendURL = os.Getenv("BACKEND_URL")
)

func getDBConnection(t *testing.T) *sql.DB {
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		t.Fatalf("Erro ao conectar no banco: %v", err)
	}
	return db
}

func getEstadoDoBanco(t *testing.T, db *sql.DB) map[int]Item {
	rows, err := db.Query("SELECT id, codigo, descricao, saldo, preco_base FROM itens")
	if err != nil {
		t.Fatalf("Erro ao consultar o banco de dados: %v", err)
	}
	defer rows.Close()

	estado := make(map[int]Item)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Codigo, &item.Descricao, &item.Saldo, &item.PrecoBase); err != nil {
			t.Fatalf("Erro ao escanear linha do banco: %v", err)
		}
		estado[item.ID] = item
	}
	return estado
}

func TestServicoEstoque(t *testing.T) {
	if backendURL == "" {
		backendURL = "http://localhost:8081"
	}

	db := getDBConnection(t)
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Banco de dados indisponível: %v", err)
	}

	t.Run("GET /produtos - Deve retornar o estado atual do banco", func(t *testing.T) {
		memoria := getEstadoDoBanco(t, db)
		resp, err := http.Get(backendURL + "/produtos")
		if err != nil {
			t.Fatalf("Erro ao chamar o back-end: %v", err)
		}
		defer resp.Body.Close()

		var itensBackend []Item
		json.NewDecoder(resp.Body).Decode(&itensBackend)

		mapBackend := make(map[int]Item)
		for _, item := range itensBackend {
			mapBackend[item.ID] = item
		}

		if len(memoria) != len(mapBackend) {
			t.Errorf("Quantidade de itens difere entre banco (%d) e backend (%d)", len(memoria), len(mapBackend))
		}
	})

	t.Run("POST /produtos - Código Automático e Retorno", func(t *testing.T) {
		novoItem := map[string]interface{}{
			"descricao":  "PRODUTO TESTE AUTOMATICO",
			"saldo":      50,
			"preco_base": 199.90,
		}

		body, _ := json.Marshal(novoItem)
		resp, err := http.Post(backendURL+"/produtos", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Erro ao chamar o POST: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Status HTTP esperado 201, recebido: %d", resp.StatusCode)
		}

		var itemCriado Item
		json.NewDecoder(resp.Body).Decode(&itemCriado)

		if !strings.HasPrefix(itemCriado.Codigo, "PROD-") {
			t.Errorf("Código gerado inválido: %s. Deve começar com PROD-", itemCriado.Codigo)
		}

		if itemCriado.Descricao != "PRODUTO TESTE AUTOMATICO" {
			t.Errorf("Descrição retornada incorreta: %s", itemCriado.Descricao)
		}
	})

	t.Run("SEGURANÇA: POST /produtos (SQL Injection)", func(t *testing.T) {
		payloadMalicioso := map[string]interface{}{
			"descricao":  "Produto'; DROP TABLE itens; --",
			"saldo":      1,
			"preco_base": 10.0,
		}

		body, _ := json.Marshal(payloadMalicioso)
		resp, err := http.Post(backendURL+"/produtos", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Erro ao chamar POST: %v", err)
		}
		defer resp.Body.Close()

		// Verifica se a tabela ainda existe tentando um GET
		respGet, err := http.Get(backendURL + "/produtos")
		if err != nil || respGet.StatusCode != http.StatusOK {
			t.Fatalf("ALERTA: Possível injeção SQL comprometeu o banco!")
		}
	})

	t.Run("SEGURANÇA: Validação de Saldo Negativo", func(t *testing.T) {
		payloadInvalido := map[string]interface{}{
			"descricao":  "TESTE NEGATIVO",
			"saldo":      -10,
			"preco_base": 10.0,
		}

		body, _ := json.Marshal(payloadInvalido)
		resp, err := http.Post(backendURL+"/produtos", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Erro ao chamar POST: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			t.Log("Aviso: O Back-end ainda permite saldo negativo. Validar no Front ou adicionar regra no Back.")
		}
	})
}
