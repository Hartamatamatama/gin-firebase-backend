package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"github.com/hartamatamatama/gin-firebase-backend/models"
	"github.com/hartamatamatama/gin-firebase-backend/services"
	"google.golang.org/api/option"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{orderService: services.NewOrderService()}
}

// Checkout - POST /v1/orders/checkout
func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := getContextUserID(c)

	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	order, err := h.orderService.Checkout(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// === OPERASI PELUNCURAN NOTIFIKASI FCM ===
	if req.FCMToken != "" {
		// 1. Inisialisasi Firebase Admin
		opt := option.WithCredentialsFile("serviceAccountKey.json") // Pastikan file ini ada!
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err == nil {
			client, err := app.Messaging(context.Background())
			if err == nil {
				// 2. Siapkan Amunisi Pesan
				message := &messaging.Message{
					Notification: &messaging.Notification{
						Title: "Order Diterima 🛍️",
						Body:  "Pesanan jam tangan mewahmu sedang kami proses ke alamat: " + req.ShippingAddress,
					},
					Token: req.FCMToken, // Tembak ke HP yang tepat
				}

				// 3. Luncurkan!
				_, sendErr := client.Send(context.Background(), message)
				if sendErr != nil {
					fmt.Println("Gagal mengirim notifikasi:", sendErr)
				} else {
					fmt.Println("🔔 Notifikasi sukses ditembakkan ke HP user!")
				}
			}
		} else {
			fmt.Println("Gagal inisialisasi Firebase Admin:", err)
		}
	}
	// ==========================================

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}

// GetMyOrders - GET /v1/orders
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := getContextUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.orderService.GetMyOrders(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetOrderByID - GET /v1/orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID := getContextUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	order, err := h.orderService.GetOrderByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Order tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})
}

// GetAllOrders - GET /v1/admin/orders (admin only)
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.orderService.GetAllOrders(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// UpdateOrderStatus - PUT /v1/admin/orders/:id/status (admin only)
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	var req struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.orderService.UpdateOrderStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status order diperbarui"})
}
