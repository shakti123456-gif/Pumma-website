package handlers

import (
	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct{ Store *db.Database }

func (h ProductHandler) List(c *gin.Context) {
	products, err := h.Store.Products()
	if err != nil {
		c.JSON(500, gin.H{"error": "could not load products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
func (h ProductHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	product, found, err := h.Store.Product(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not load product"})
		return
	}
	if !found {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}
	c.JSON(200, product)
}
func (h ProductHandler) Create(c *gin.Context) {
	var product models.Product
	if c.ShouldBindJSON(&product) != nil || !validProduct(product) {
		c.JSON(400, gin.H{"error": "invalid product"})
		return
	}
	saved, err := h.Store.CreateProduct(product)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not save product"})
		return
	}
	c.JSON(201, saved)
}
func (h ProductHandler) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var product models.Product
	if c.ShouldBindJSON(&product) != nil || !validProduct(product) {
		c.JSON(400, gin.H{"error": "invalid product"})
		return
	}
	saved, found, err := h.Store.UpdateProduct(id, product)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not save product"})
		return
	}
	if !found {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}
	c.JSON(200, saved)
}
func (h ProductHandler) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	found, err := h.Store.DeleteProduct(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not save product"})
		return
	}
	if !found {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}
	c.Status(204)
}
func pathID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid ID"})
		return 0, false
	}
	return id, true
}
func validProduct(product models.Product) bool {
	return strings.TrimSpace(product.Name) != "" && product.Price >= 0 && product.Stock >= 0
}
