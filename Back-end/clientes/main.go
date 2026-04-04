package main

import (
	"clientes/database"
	"clientes/handlers"
	"fmt"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	database.InitDB()
	defer database.DB.Close()

	mux := http.NewServeMux()

	// Handler unificado para tratar /clientes e /clientes/{id}
	mux.HandleFunc("/clientes", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetClientes(w, r)
	})

	mux.HandleFunc("/clientes/", func(w http.ResponseWriter, r *http.Request) {
		// Se for apenas "/clientes/", redireciona para a lista (sem 301)
		if r.URL.Path == "/clientes/" {
			handlers.GetClientes(w, r)
			return
		}
		// Caso contrário, trata como ID
		handlers.GetClienteByID(w, r)
	})

	port := ":8083"
	fmt.Printf("Serviço de Clientes rodando na porta %s\n", port)
	
	// Aplica o CORS em todo o Mux
	log.Fatal(http.ListenAndServe(port, enableCORS(mux)))
}
