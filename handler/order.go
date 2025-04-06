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

	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}

}
