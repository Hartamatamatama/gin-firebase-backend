package services

import (
	"github.com/hartamatamatama/gin-firebase-backend/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{productRepo: repositories.NewProductRepository()}
}