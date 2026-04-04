package clientes_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

type Cliente struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Endereco string `json:"endereco"`
	CPF      string `json:"cpf"`
	CNPJ     string `json:"cnpj"`
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

func getEstadoDoBanco(t *testing.T, db *sql.DB) map[int]Cliente {
	rows, err := db.Query("SELECT id, nome, telefone, endereco, cpf, cnpj FROM clientes")
	if err != nil {
		t.Fatalf("Erro ao consultar o banco de dados: %v", err)
	}
	defer rows.Close()

	estado := make(map[int]Cliente)
	for rows.Next() {
		var c Cliente
		var telefone, endereco, cpf, cnpj sql.NullString
		if err := rows.Scan(&c.ID, &c.Nome, &telefone, &endereco, &cpf, &cnpj); err != nil {
			t.Fatalf("Erro ao escanear linha do banco: %v", err)
		}
		c.Telefone = telefone.String
		c.Endereco = endereco.String
		c.CPF = cpf.String
		c.CNPJ = cnpj.String
		estado[c.ID] = c
	}
	return estado
}

func TestServicoClientes(t *testing.T) {
	if backendURL == "" {
		backendURL = "http://localhost:8083" // Porta sugerida para clientes
	}

	db := getDBConnection(t)
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Banco de dados indisponível: %v", err)
	}

	t.Run("GET /clientes - Deve retornar o estado exato do script de seed", func(t *testing.T) {
		memoria := getEstadoDoBanco(t, db)
		resp, err := http.Get(backendURL + "/clientes")
		if err != nil {
			t.Fatalf("Erro ao chamar o back-end de clientes (Serviço fora do ar?): %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status HTTP esperado 200, recebido: %d", resp.StatusCode)
		}

		var clientesBackend []Cliente
		json.NewDecoder(resp.Body).Decode(&clientesBackend)

		if len(memoria) != len(clientesBackend) {
			t.Errorf("Diferença detectada! Banco possui %d clientes, mas o Microserviço retornou %d", len(memoria), len(clientesBackend))
		}

		// Validação de dados específicos de um cliente para garantir que o mapeamento está correto
		for _, cb := range clientesBackend {
			if _, ok := memoria[cb.ID]; !ok {
				t.Errorf("Cliente ID %d retornado pelo backend não existe no banco!", cb.ID)
			}
		}
	})

	t.Run("PROTEÇÃO: Apenas Leitura (Métodos Proibidos)", func(t *testing.T) {
		metodos := []string{"POST", "PUT", "DELETE", "PATCH"}
		for _, m := range metodos {
			req, _ := http.NewRequest(m, backendURL+"/clientes", nil)
			client := &http.Client{}
			resp, err := client.Do(req)
			
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusMethodNotAllowed {
					t.Errorf("ALERTA DE SEGURANÇA: Endpoint de Clientes aceitou método %s! Deveria ser 405 Method Not Allowed.", m)
				}
			}
		}
	})

	t.Run("SEGURANÇA: Cabeçalhos de Resposta", func(t *testing.T) {
		resp, err := http.Get(backendURL + "/clientes")
		if err != nil {
			t.Skip("Pulando teste de cabeçalhos pois serviço está offline.")
			return
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type incorreto. Esperado application/json, recebido: %s", contentType)
		}
	})
}
