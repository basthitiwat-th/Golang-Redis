package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type productRepositoryRedis struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewProductReposittoryRedis(db *gorm.DB, redisClient *redis.Client) ProductRepository {
	db.AutoMigrate(&product{})
	mockData(db)
	return productRepositoryRedis{db, redisClient}
}

func (r productRepositoryRedis) GetProducts() (products []product, err error) {
	key := "repository:GetProducts"
	//Redis Get
	produsctJson, err := r.redisClient.Get(context.Background(), key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(produsctJson), &products)
		if err == nil {
			fmt.Println("Redis")
			return products, nil
		}
	}
	//Database
	err = r.db.Order("quantity desc").Limit(30).Find(&products).Error
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}
	//Redis Set
	err = r.redisClient.Set(context.Background(), key, string(data), time.Second*50).Err()
	if err != nil {
		return nil, err
	}
	fmt.Println("database")
	return products, nil
}
