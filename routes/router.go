package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hartamatamatama/gin-firebase-backend/handlers"
	"github.com/hartamatamatama/gin-firebase-backend/middleware"
)

func SetupRouter() *gin.Engine {
	// gin.Default() sudah include Logger & Recovery middleware
	r := gin.Default()

	// CORS Middleware (izinkan request dari aplikasi Flutter nantinya)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Init handlers
	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()

	// API v1 group
	v1 := r.Group("/v1")
	{
		// Health check (tidak perlu auth, untuk tes server hidup atau tidak)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "gin-firebase-backend"})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		// Terima Firebase token, return Backend JWT
		auth.POST("/verify-token", authHandler.VerifyToken)

		// Protected routes (butuh Backend JWT)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware()) // Pasang satpam di sini

		// Products (semua user terautentikasi bisa GET)
		products := protected.Group("/products")
		products.GET("", productHandler.GetAll)      // GET /v1/products
		products.GET("/:id", productHandler.GetByID) // GET /v1/products/:id

		// Create, Update, Delete (hanya admin)
		adminProducts := products.Group("")
		adminProducts.Use(middleware.AdminOnly()) // Pasang satpam admin di sini
		adminProducts.POST("", productHandler.Create)       // POST /v1/products
		adminProducts.PUT("/:id", productHandler.Update)    // PUT /v1/products/:id
		adminProducts.DELETE("/:id", productHandler.Delete) // DELETE /v1/products/:id
	}

	return r
}