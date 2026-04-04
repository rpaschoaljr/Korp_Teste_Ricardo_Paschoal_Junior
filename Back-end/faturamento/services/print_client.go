package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PrintClient struct {
	BaseURL string
}

func NewPrintClient() *PrintClient {
	url := os.Getenv("IMPRESSAO_URL")
	if url == "" {
		url = "http://impressao_api:8084"
	}
	return &PrintClient{BaseURL: url}
}

// GeneratePDF solicita a geração do PDF e retorna o binário
func (p *PrintClient) GeneratePDF(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(p.BaseURL+"/gerar-pdf", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("serviço de impressão indisponível: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("impressão falhou com status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *PrintClient) CheckHealth() error {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(p.BaseURL + "/health")
	if err != nil {
		return fmt.Errorf("serviço de impressão offline")
	}
	defer resp.Body.Close()
	return nil
}
