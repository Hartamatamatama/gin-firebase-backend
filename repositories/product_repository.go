package repositories

import (
	"github.com/hartamatamatama/gin-firebase-backend/config"
	"github.com/hartamatamatama/gin-firebase-backend/models"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// FindAll mengambil semua produk aktif dengan pagination
func (r *ProductRepository) FindAll(page, limit int, category string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	query := config.DB.Model(&models.Product{}).Where("is_active = ?", true)

	// Filter by category jika ada
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Hitung total untuk pagination
	query.Count(&total)

	// Ambil data dengan offset & limit
	offset := (page - 1) * limit
	result := query.Offset(offset).Limit(limit).Find(&products)

	return products, total, result.Error
}

// FindByID mengambil satu produk berdasarkan ID
func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	result := config.DB.First(&product, id)
	return &product, result.Error
}