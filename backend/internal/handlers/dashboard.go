package handlers

import (
	"ecommerce-backend/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DashboardHandler struct{ Store *db.Database }

func (h DashboardHandler) Overview(c *gin.Context) {
	products, err := h.Store.Products()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load dashboard"})
		return
	}
	orders, err := h.Store.Orders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load dashboard"})
		return
	}
	var revenue float64
	for _, order := range orders {
		revenue += order.Total
	}
	lowStock := 0
	for _, product := range products {
		if product.Stock < 10 {
			lowStock++
		}
	}
	c.JSON(http.StatusOK, gin.H{"products": len(products), "orders": len(orders), "revenue": revenue, "lowStock": lowStock})
}
