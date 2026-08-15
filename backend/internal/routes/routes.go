package routes

import (
	"net/http"

	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, repository *db.Database) {
	products := handlers.ProductHandler{Store: repository}
	orders := handlers.OrderHandler{Store: repository}
	dashboard := handlers.DashboardHandler{Store: repository}
	auth := handlers.NewAuthHandler()

	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := router.Group("/api")
	api.POST("/auth/login", auth.Login)
	api.GET("/products", products.ListPublic)
	api.GET("/products/:id", products.Get)
	api.POST("/orders", orders.Create)

	admin := api.Group("/admin", auth.RequireAdmin())
	admin.GET("/products", products.ListAdmin)
	admin.POST("/products", products.Create)
	admin.PUT("/products/:id", products.Update)
	admin.DELETE("/products/:id", products.Delete)
	admin.GET("/orders", orders.List)
	admin.PATCH("/orders/:id", orders.UpdateStatus)
	admin.GET("/dashboard", dashboard.Overview)
}
