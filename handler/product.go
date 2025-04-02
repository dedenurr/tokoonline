package handler

import (
	"database/sql"
	"errors"
	"log"

	"github.com/dedenurr/tokoonline/model"
	"github.com/gin-gonic/gin"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//  TODO: ambil dari database
		products, err := model.SelectProduct(db)
		if err != nil {
			log.Printf("terjadi kesalahan saat mengambil data produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Gagal mengambil data produk"})
			return
		}

		// TODO: Berikan response ke client
		c.JSON(200, gin.H{"products": products})
	}
}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO Baca id dari parameter
		id := c.Param("id")

		// TODO: ambil dari database
		product, err := model.SelectProductByID(db, id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("terjadi kesalahan saat mengambil data produk: %v\n", err)
				c.JSON(404, gin.H{"error": "Produck Tidak Ditermukan"})
				return
			}
			log.Printf("terjadi kesalahan saat mengambil data produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi Kesalahan pada server"})
			return
		}

		//  TODO: Berikan response ke client
		c.JSON(200, gin.H{"product": product})
	}

}
