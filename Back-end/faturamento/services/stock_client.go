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

func (s *StockClient) GetProducts() ([]models.ProductInfo, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(s.BaseURL + "/produtos")
	if err != nil {
		return nil, fmt.Errorf("servico de estoque offline: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("estoque retornou erro %d", resp.StatusCode)
	}

	var products []models.ProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}
	return products, nil
}
