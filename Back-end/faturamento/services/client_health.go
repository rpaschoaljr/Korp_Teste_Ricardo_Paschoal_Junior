package services

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type ClientServiceHealth struct {
	BaseURL string
}

func NewClientServiceHealth() *ClientServiceHealth {
	url := os.Getenv("CLIENTES_URL")
	if url == "" {
		url = "http://clientes_api:8083"
	}
	return &ClientServiceHealth{BaseURL: url}
}

// CheckHealth verifica se o serviço de clientes está respondendo
func (c *ClientServiceHealth) CheckHealth() error {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(c.BaseURL + "/clientes")
	if err != nil {
		return fmt.Errorf("serviço de clientes indisponível: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("serviço de clientes retornou erro: %d", resp.StatusCode)
	}

	return nil
}
