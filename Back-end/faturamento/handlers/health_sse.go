package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemHealth struct {
	Estoque     bool `json:"estoque"`
	Faturamento bool `json:"faturamento"`
	Clientes    bool `json:"clientes"`
	IA          bool `json:"ia"`
}

func HealthStream(c *gin.Context) {
	// Configura headers para SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Access-Control-Allow-Origin", "*")

	// URLs internas (usando nomes dos containers do Docker)
	estoqueURL := os.Getenv("ESTOQUE_URL")
	if estoqueURL == "" { estoqueURL = "http://estoque_api:8081" }
	clientesURL := os.Getenv("CLIENTES_URL")
	if clientesURL == "" { clientesURL = "http://clientes_api:8083" }

	// Loop de monitoramento
	c.Stream(func(w io.Writer) bool {
		health := SystemHealth{
			Faturamento: true, // Se este código está rodando, o faturamento está ok
			IA:          true, // IA roda no mesmo serviço de faturamento
		}

		// Checagem ultra-rápida de Estoque
		client := http.Client{Timeout: 1 * time.Second}
		respE, errE := client.Get(estoqueURL + "/produtos")
		health.Estoque = (errE == nil && respE.StatusCode == http.StatusOK)
		if respE != nil { respE.Body.Close() }

		// Checagem ultra-rápida de Clientes
		respC, errC := client.Get(clientesURL + "/clientes")
		health.Clientes = (errC == nil && respC.StatusCode == http.StatusOK)
		if respC != nil { respC.Body.Close() }

		// Transforma em JSON
		jsonData, _ := json.Marshal(health)

		// Envia para o cliente no formato SSE: "data: {json}\n\n"
		fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
		
		// Espera 3 segundos antes da próxima checagem
		time.Sleep(3 * time.Second)
		
		return true // Continua o stream
	})
}
