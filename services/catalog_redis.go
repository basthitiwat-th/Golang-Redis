package services

import (
	"GOREDIS/repositories"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type catalogServiceRedis struct {
	productRepo repositories.ProductRepository
	redisClient *redis.Client
}

func NewCatalogServiceRedis(productRepo repositories.ProductRepository, redisClient *redis.Client) CatalogService {
	return catalogServiceRedis{productRepo, redisClient}
}

func (s catalogServiceRedis) GetProducts() (products []Product, err error) {
	key := "service::GetProducts"

	//Redis GET
	if productJson, err := s.redisClient.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(productJson), &products) == nil {
			fmt.Println("Redis")
			return products, nil
		}
	}

	//Repository
	producsDB, err := s.productRepo.GetProducts()
	if err != nil {
		return nil, err
	}

	for _, p := range producsDB {
		products = append(products, Product{
			ID:       p.ID,
			Name:     p.Name,
			Quantity: p.Quantity,
		})
	}

	//Redis SET
	if data, err := json.Marshal(products); err == nil {
		s.redisClient.Set(context.Background(), key, string(data), time.Second*20)
	}
	fmt.Println("Database")
	return products, nil
}
