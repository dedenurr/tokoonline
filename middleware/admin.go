package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil kunci rahasia
		key := os.Getenv("ADMIN_KEY")
		if key == "" {
			c.JSON(500, gin.H{"error": "ADMIN_KEY tidak ditemukan"})
			c.Abort()
			return
		}

		// ambil header Authorization
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Akses ditolak, header Authorization tidak ditemukan"})
			c.Abort()
			return
		}

		// validasi kunci admin dengan header Authorization
		if auth != key {
			c.JSON(403, gin.H{"error": "Akses ditolak, kunci admin tidak valid"})
			c.Abort()
			return
		}

		// melanjutkan ke handler berikutnya
		c.Next()
	}
}
