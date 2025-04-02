package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/dedenurr/tokoonline/handler"
	"github.com/dedenurr/tokoonline/middleware"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.GET("/api/v1/products", handler.ListProducts(db))
	r.GET("/api/v1/products/:id", handler.GetProduct(db))
	r.POST("/api/v1/checkout")

	r.POST("/api/v1/orders/:id/confirm")
	r.GET("/api/v1/orders/:id")

	r.POST("/admin/v1/products", middleware.AdminOnly(), handler.CreateProduct(db))
	r.PUT("/admin/v1/products/:id", middleware.AdminOnly(), handler.UpdateProduct(db))
	r.DELETE("/admin/v1/products/:id", middleware.AdminOnly(), handler.DeleteProduct(db))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Server Handler
	if err = server.ListenAndServe(); err != nil {
		fmt.Printf("Gagal menjalankan server: %v\n", err)
		os.Exit(1)
	}
}
