package services

import (
	"github.com/hartamatamatama/gin-firebase-backend/models"
	"github.com/hartamatamatama/gin-firebase-backend/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{productRepo: repositories.NewProductRepository()}
}

func (s *ProductService) GetAll(page, limit int, category string) ([]models.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return s.productRepo.FindAll(page, limit, category)
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) Create(req models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}
	err := s.productRepo.Create(product)
	return product, err
}