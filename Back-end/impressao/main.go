package main

import (
	"impressao_service/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", handlers.HealthCheck)
	r.POST("/gerar-pdf", handlers.GeneratePDF)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("Microserviço de Impressão rodando na porta %s", port)
	r.Run(":" + port)
}
