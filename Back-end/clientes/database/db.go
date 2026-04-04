package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Aguardando banco de dados: %v", err)
	} else {
		fmt.Println("Conectado ao banco de dados com sucesso (Clientes)")
	}
}
