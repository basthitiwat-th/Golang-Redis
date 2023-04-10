package repositories

import (
	"gorm.io/gorm"
)

type productRepositoryDB struct {
	db *gorm.DB
}

func NewProductReposittoryDB(db *gorm.DB) ProductRepository {
	db.AutoMigrate(&product{})
	mockData(db)
	return productRepositoryDB{db: db}
}

func (r productRepositoryDB) GetProducts() (products []product, err error) {
	err = r.db.Order("quantity DESC").Limit(30).Find(&products).Error
	return
}
