package handlers

import (
	"GOREDIS/services"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type catalogHandlerRedis struct {
	catalogSrv  services.CatalogService
	redisClient *redis.Client
}

func NewCatalogHandlerRedis(catalogSrv services.CatalogService, redisClient *redis.Client) CatalogHandler {
	return catalogHandlerRedis{catalogSrv, redisClient}
}

func (h catalogHandlerRedis) GetProducts(c *fiber.Ctx) error {
	key := "handler::GetProducts"

	//Redis GET
	if responseJson, err := h.redisClient.Get(context.Background(), key).Result(); err == nil {
		fmt.Println("redis")
		c.Set("Content-type", "application/json")
		return c.SendString(responseJson)
	}

	//GET Service
	products, err := h.catalogSrv.GetProducts()
	if err != nil {
		return err
	}

	reponse := fiber.Map{
		"status":  "ok",
		"product": products,
	}

	//Redis SET
	if data, err := json.Marshal(reponse); err == nil {
		h.redisClient.Set(context.Background(), key, string(data), time.Second*10)
	}
	fmt.Println("database")
	fmt.Println(reponse)
	return c.JSON(reponse)
}
