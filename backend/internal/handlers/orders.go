package handlers

import (
	"net/http"
	"strings"
	"time"

	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct{ Store *db.Database }

func (h OrderHandler) List(c *gin.Context) {
	ctx, cancel := db.WithTimeout()
	defer cancel()

	orders, err := h.Store.Orders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h OrderHandler) Create(c *gin.Context) {
	var order models.Order
	if c.ShouldBindJSON(&order) != nil || strings.TrimSpace(order.Customer) == "" || len(order.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer and at least one item are required"})
		return
	}
	for _, item := range order.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "each item needs a valid product and quantity"})
			return
		}
	}

	order.Status = "pending"
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now().UTC()
	}

	ctx, cancel := db.WithTimeout()
	defer cancel()

	saved, err := h.Store.CreateOrder(ctx, order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, saved)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	ctx, cancel := db.WithTimeout()
	defer cancel()

	order, found, err := h.Store.UpdateOrderStatus(ctx, id, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save order"})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}
