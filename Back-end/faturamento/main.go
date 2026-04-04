package main

import (
	"faturamento_service/database"
	"faturamento_service/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer database.DB.Close()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/faturas", handlers.ListFaturas)
	r.POST("/faturas", handlers.CreateFatura)
	r.POST("/faturas/:id/imprimir", handlers.PrintFatura)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Microserviço de Faturamento rodando na porta %s", port)
	r.Run(":" + port)
}
