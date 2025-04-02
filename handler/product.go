package handler

import (
	"database/sql"

	"github.com/dedenurr/tokoonline/model"
	"github.com/gin-gonic/gin"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	//  TODO: ambil dari database
	return func(c *gin.Context) {
		products, err := model.SelectProduct(db)
		if err != nil {
			c.JSON(500, gin.H{"error": "Gagal mengambil data produk"})
			return
		}

		// TODO: Berikan response ke client
		c.JSON(200, gin.H{"products": products})
	}
}

func GetProduct(c *gin.Context) {
	// TODO: ambil dari database

	//  TODO: Berikan response ke client

}
