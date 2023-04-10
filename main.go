package main

import (
	"GOREDIS/handlers"
	"GOREDIS/repositories"
	"GOREDIS/services"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDatabase()
	redisClient := initRedis()
	_ = redisClient

	productRepo := repositories.NewProductReposittoryDB(db)
	productService := services.NewCatalogServiceRedis(productRepo, redisClient)
	productHandler := handlers.NewCatalogHandler(productService)

	app := fiber.New()
	app.Get("/products", productHandler.GetProducts)
	app.Listen(":8000")

}

func initDatabase() *gorm.DB {
	dial := mysql.Open("root:Basbm031197#@tcp(localhost:3306)/infinitas")
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		panic((err))
	}
	return db

}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
