package main

import (
	"context"
	"log"
	"os"

	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("could not load .env: ", err)
	}
	repository, err := db.Connect(context.Background())
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer repository.Close()
	router := gin.Default()
	router.Use(middleware.CORS())
	routes.Register(router, repository)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
