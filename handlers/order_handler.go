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

// ConfirmPayment - PUT /v1/orders/:id/confirm-payment
func (h *OrderHandler) ConfirmPayment(c *gin.Context) {
	userID := getContextUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	order, err := h.orderService.ConfirmPayment(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Order tidak ditemukan"})
		return
	}

	// === KIRIM NOTIFIKASI FCM SETELAH BAYAR BERHASIL ===
	if order.FCMToken != "" {
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		app, fErr := firebase.NewApp(context.Background(), nil, opt)
		if fErr == nil {
			client, mErr := app.Messaging(context.Background())
			if mErr == nil {
				message := &messaging.Message{
					Notification: &messaging.Notification{
						Title: "Pembayaran Berhasil ✅",
						Body:  "Pesanan #" + fmt.Sprint(order.ID) + " telah dibayar. Stok sudah dikurangi, pesanan sedang diproses.",
					},
					Token: order.FCMToken,
				}
				_, sErr := client.Send(context.Background(), message)
				if sErr != nil {
					fmt.Println("Gagal kirim notifikasi:", sErr)
				}
			}
		}
	}
	// ==========================================

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pembayaran dikonfirmasi"})
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
