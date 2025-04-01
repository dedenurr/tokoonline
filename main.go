package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
)

func main() {
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))

	if err != nil {
		fmt.Printf("Gagal Membuat koneksi ke database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Gagal memverifikasi koneksi ke database: %v\n", err)
		os.Exit(1)
	}
}
