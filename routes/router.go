package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hartamatamatama/gin-firebase-backend/handlers"
	"github.com/hartamatamatama/gin-firebase-backend/middleware"
)

func SetupRouter() *gin.Engine {
	// Gunakan gin.New() agar kita bisa kontrol penuh urutan middleware
	// (tidak pakai gin.Default() yang auto-include logger bawaan gin)
	r := gin.New()
	r.Use(gin.Recovery()) // panic recovery tetap diperlukan
	r.Use(middleware.HTTPLogger())

	// ─── CORS Middleware ───────────────────────────────────────
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

	// ─── Init handlers ────────────────────────────────────────
	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()
	cartHandler := handlers.NewCartHandler()
	orderHandler := handlers.NewOrderHandler()

	// ─── API v1 group ─────────────────────────────────────────
	v1 := r.Group("/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "mycatalog-backend"})
		})

		// ── Auth routes (public) ──────────────────────────────
		auth := v1.Group("/auth")
		{
			auth.POST("/verify-token", authHandler.VerifyToken)
		}

		// ── Products (public — bisa lihat tanpa login) ─────────────
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
		}

		// ── Protected routes (butuh JWT) ──────────────────────
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Admin products
			protected.POST("/products", productHandler.Create, middleware.AdminOnly())
			protected.PUT("/products/:id", productHandler.Update, middleware.AdminOnly())
			protected.DELETE("/products/:id", productHandler.Delete, middleware.AdminOnly())

			// Cart
			cart := protected.Group("/cart")
			{
				cart.GET("", cartHandler.GetCart)           // GET    /v1/cart
				cart.POST("", cartHandler.AddToCart)        // POST   /v1/cart
				cart.PUT("/:id", cartHandler.UpdateCartItem) // PUT    /v1/cart/:id
				cart.DELETE("/:id", cartHandler.RemoveCartItem) // DELETE /v1/cart/:id
				cart.DELETE("", cartHandler.ClearCart)      // DELETE /v1/cart
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.POST("/checkout", orderHandler.Checkout)          // POST   /v1/orders/checkout
				orders.GET("", orderHandler.GetMyOrders)                 // GET    /v1/orders
				orders.GET("/:id", orderHandler.GetOrderByID)            // GET    /v1/orders/:id
				orders.PUT("/:id/confirm-payment", orderHandler.ConfirmPayment) // PUT /v1/orders/:id/confirm-payment
			}

			// Admin — order management
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/orders", orderHandler.GetAllOrders)                     // GET /v1/admin/orders
				admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)     // PUT /v1/admin/orders/:id/status
			}
		}
	}

	return r
}
