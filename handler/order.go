package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/dedenurr/tokoonline/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

		// TODO: Buat Kata Sandi

		//siapkan passcode
		passcode := generatePasscode(5)

		// TODO: Hash Kata Sandi
		hashPasscode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)
		if err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		hashPasscodeStr := string(hashPasscode)

		// TODO: Buat Order & Detail
		order := model.Order{
			ID:         uuid.New().String(),
			Email:      CheckoutOrder.Email,
			Address:    CheckoutOrder.Address,
			Passcode:   &hashPasscodeStr,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}

		for _, p := range products {
			total := p.Price * int64(orderQty[p.ID])

			detail := model.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQty[p.ID],
				Price:     p.Price,
				Total:     total,
			}

			details = append(details, detail)

			order.GrandTotal += total
		}

		model.CreateOrder(db, order, details)

		orderWithDetail := model.OrderWithDetail{
			Order:   order,
			Details: details,
		}

		orderWithDetail.Passcode = &passcode

		c.JSON(200, orderWithDetail)

	}

}

func generatePasscode(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	random := make([]byte, length)
	for i := range random {
		random[i] = charset[randomGen.Intn(len(charset))]
	}
	return string(random)
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO : Ambil id dari params
		id := c.Param("id")

		// TODO : Baca request Body
		var confirmReq model.Confirm
		if err := c.BindJSON(&confirmReq); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data konfirmasi tidak valid"})
			return
		}

		// TODO : Ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat meembaca data pesanan: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		if order.Passcode == nil {
			log.Println("Passcode tidak valid")
			c.JSON(500, gin.H{"error": "Terjadi Kesalahan Pada Server"})
			return
		}

		// TODO : cocokan kata sandi Pesanan
		if err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(confirmReq.Passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat mencocokan kata sandi: %v\n", err)
			c.JSON(401, gin.H{"error": "tidak di ijinkan mengakses pesanan"})
			return
		}

		// TODO : Pastikan Pesanan Belum Dibayar
		if order.PaidAt != nil {
			log.Println("Pesanan sudah dibayar")
			c.JSON(400, gin.H{"error": "Pesanan sudah dibayar"})
			return
		}

		// TODO : Cocokan Jumlah Pembayaran
		if order.GrandTotal != confirmReq.Amount {
			log.Printf("Jumlah harga tidak sesuai: %d\n", confirmReq.Amount)
			c.JSON(400, gin.H{"error": "Jumlah pembayaran tidak sesuai"})
			return
		}

		// TODO : Update informasi pesanan
		current := time.Now()
		if err = model.UpdateOrderByID(db, id, confirmReq, current); err != nil {
			log.Printf("Terjadi kesalahan saat memperbarui data pesanan: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		order.Passcode = nil

		order.PaidAt = &current
		order.PaidBank = &confirmReq.Bank
		order.PaidAccountNumber = &confirmReq.AccountNumber

		c.JSON(200, order)
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO : Ambil id dari params
		id := c.Param("id")

		// TODO : ambil password dari query param
		passcode := c.Query("passcode")

		// TODO : Ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat meembaca data pesanan: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		if order.Passcode == nil {
			log.Println("Passcode tidak valid")
			c.JSON(500, gin.H{"error": "Terjadi Kesalahan Pada Server"})
			return
		}

		// TODO : cocokan kata sandi Pesanan
		if err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat mencocokan kata sandi: %v\n", err)
			c.JSON(401, gin.H{"error": "tidak di ijinkan mengakses pesanan"})
			return
		}

		order.Passcode = nil

		c.JSON(200, order)
	}

}
