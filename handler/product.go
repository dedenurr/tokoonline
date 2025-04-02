package handler

import (
	"database/sql"
	"errors"
	"log"

	"github.com/dedenurr/tokoonline/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil data dari request body
		var product model.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(400, gin.H{"error": "Data Product Tidak Valid"})
			return
		}

		//  atur id dari uuid
		product.ID = uuid.New().String()

		// simpan product ke database
		if err := model.InsertProduct(db, product); err != nil {
			c.JSON(500, gin.H{"error": "Terjadi Kesalahan pada server"})
			return
		}

		c.JSON(201, product)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//  ambil id dari produk URL
		id := c.Param("id")

		// ambil data dari request body
		var productReq model.Product
		if err := c.BindJSON(&productReq); err != nil {
			c.JSON(400, gin.H{"error": "Data Product Tidak Valid"})
			return
		}

		// ambil darta product dara database
		product, err := model.SelectProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(400, gin.H{"error": "Product Tidak Ditermukan"})
				return
			}
		}

		// update harga nama jika tidak kosong
		if productReq.Name != "" {
			product.Name = productReq.Name
		}

		// update harga produk jika tidak kosong
		if productReq.Price != 0 {
			product.Price = productReq.Price
		}

		if err := model.UpdateProduct(db, product); err != nil {
			c.JSON(500, gin.H{"error": "Terjadi Kesalahan pada server"})
			return
		}
		// tampilkan data produk yang diupdate
		c.JSON(200, product)
	}
}

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
