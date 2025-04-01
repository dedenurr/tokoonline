package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
)

func main() {
	// Koneksi Database
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		fmt.Printf("Gagal Membuat koneksi ke database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Verifikasi Koneksi Database
	if err = db.Ping(); err != nil {
		fmt.Printf("Gagal memverifikasi koneksi ke database: %v\n", err)
		os.Exit(1)
	}

	// Migrasi Database
	if _, err = migrate(db); err != nil {
		fmt.Printf("Gagal membuat tabel: %v\n", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	// Server Handler
	if err = server.ListenAndServe(); err != nil {
		fmt.Printf("Gagal menjalankan server: %v\n", err)
		os.Exit(1)
	}
}
