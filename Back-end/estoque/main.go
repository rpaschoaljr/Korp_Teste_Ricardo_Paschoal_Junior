package main

import (
	"estoque_service/database"
	"estoque_service/handlers"
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

	r.GET("/produtos", handlers.GetProducts)
	r.POST("/produtos", handlers.CreateProduct)
	r.POST("/produtos/baixa", handlers.UpdateStock)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Microserviço de Estoque rodando na porta %s", port)
	r.Run(":" + port)
}
