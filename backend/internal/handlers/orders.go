package handlers

import (
	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/models"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type OrderHandler struct{ Store *db.Database }

func (h OrderHandler) List(c *gin.Context) {
	orders, err := h.Store.Orders()
	if err != nil {
		c.JSON(500, gin.H{"error": "could not load orders"})
		return
	}
	c.JSON(200, orders)
}
func (h OrderHandler) Create(c *gin.Context) {
	var order models.Order
	if c.ShouldBindJSON(&order) != nil || strings.TrimSpace(order.Customer) == "" || len(order.Items) == 0 {
		c.JSON(400, gin.H{"error": "customer and at least one item are required"})
		return
	}
	order.Status = "pending"
	order.CreatedAt = time.Now().UTC()
	saved, err := h.Store.CreateOrder(order)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not save order"})
		return
	}
	c.JSON(201, saved)
}
func (h OrderHandler) UpdateStatus(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var request struct {
		Status string `json:"status"`
	}
	if c.ShouldBindJSON(&request) != nil || request.Status == "" {
		c.JSON(400, gin.H{"error": "status is required"})
		return
	}
	order, found, err := h.Store.UpdateOrderStatus(id, request.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not save order"})
		return
	}
	if !found {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}
	c.JSON(200, order)
}
