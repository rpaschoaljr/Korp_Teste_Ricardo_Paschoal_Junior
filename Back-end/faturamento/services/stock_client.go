package services

import (
	"bytes"
	"encoding/json"
	"faturamento_service/models"
	"fmt"
	"net/http"
	"os"
	"time"
)

type StockClient struct {
	BaseURL string
}

func NewStockClient() *StockClient {
	url := os.Getenv("ESTOQUE_URL")
	if url == "" {
		url = "http://estoque:8081"
	}
	return &StockClient{BaseURL: url}
}

func (s *StockClient) ReduceStock(updates []models.StockUpdate) error {
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(s.BaseURL+"/produtos/baixa", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("serviço de estoque indisponível: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("estoque recusou a baixa: %s", errResp["error"])
	}

	return nil
}
