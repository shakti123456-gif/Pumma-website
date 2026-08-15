package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/models"
	"github.com/gin-gonic/gin"
)

var (
	validSports  = map[string]bool{"Running": true, "Training": true, "Football": true, "Basketball": true, "Lifestyle": true}
	validGenders = map[string]bool{"Men": true, "Women": true, "Unisex": true}
	validTypes   = map[string]bool{"Clothing": true, "Shoes": true}
)

type ProductHandler struct{ Store *db.Database }

func (h ProductHandler) ListPublic(c *gin.Context) {
	h.list(c, true)
}

func (h ProductHandler) ListAdmin(c *gin.Context) {
	h.list(c, false)
}

func (h ProductHandler) list(c *gin.Context, activeOnly bool) {
	ctx, cancel := db.WithTimeout()
	defer cancel()

	filter := db.NormalizeFilter(c.Query("sport"), c.Query("gender"), c.Query("type"))
	filter.ActiveOnly = activeOnly

	products, err := h.Store.Products(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h ProductHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	ctx, cancel := db.WithTimeout()
	defer cancel()

	product, found, err := h.Store.Product(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load product"})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if !product.Active {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h ProductHandler) Create(c *gin.Context) {
	var product models.Product
	if c.ShouldBindJSON(&product) != nil || !validProduct(product) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product"})
		return
	}
	ctx, cancel := db.WithTimeout()
	defer cancel()

	saved, err := h.Store.CreateProduct(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save product"})
		return
	}
	c.JSON(http.StatusCreated, saved)
}

func (h ProductHandler) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var product models.Product
	if c.ShouldBindJSON(&product) != nil || !validProduct(product) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product"})
		return
	}
	ctx, cancel := db.WithTimeout()
	defer cancel()

	saved, found, err := h.Store.UpdateProduct(ctx, id, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save product"})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, saved)
}

func (h ProductHandler) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	ctx, cancel := db.WithTimeout()
	defer cancel()

	found, err := h.Store.DeleteProduct(ctx, id)
	if errors.Is(err, db.ErrProductInOrders) {
		c.JSON(http.StatusConflict, gin.H{"error": "product is used in existing orders and cannot be deleted"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product"})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func pathID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return 0, false
	}
	return id, true
}

func validProduct(product models.Product) bool {
	return strings.TrimSpace(product.Name) != "" &&
		validSports[product.Sport] &&
		validGenders[product.Gender] &&
		validTypes[product.Type] &&
		product.Price >= 0 &&
		product.Stock >= 0
}
