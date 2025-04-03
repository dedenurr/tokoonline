package handler

import (
	"database/sql"
	"log"

	"github.com/dedenurr/tokoonline/model"
	"github.com/gin-gonic/gin"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Ambil data pesanan dari request
		var CheckoutOrder model.Checkout
		if err := c.BindJSON(&CheckoutOrder); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data produk tidak valid"})
			return
		}

		ids := []string{}
		orderQty := make(map[string]int32)
		for _, o := range CheckoutOrder.Products {
			ids = append(ids, o.ID)
			orderQty[o.ID] = o.Quantity
		}

		// TODO : Ambil Product data dari database
		products, err := model.SelectProductIn(db, ids)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil data produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan saat mengambil data produk"})
			return
		}

		c.JSON(200, gin.H{"products": products})

		// TODO: Buat Kata Sandi

		// TODO: Hash Kata Sandi

		// TODO: Buat Order & Detail

	}

}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}

}
