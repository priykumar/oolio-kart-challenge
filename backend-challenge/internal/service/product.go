package service

import (
	"github.com/priykumar/oolio-kart-challenge/internal/model"
	"github.com/priykumar/oolio-kart-challenge/internal/repo"
)

type ProductService interface {
	GetAllAvailableProducts() ([]model.Product, error)
	GetProductById(int64) (*model.Product, error)
}

type productService struct {
	db repo.KartRepository
}

func NewProductService(db repo.KartRepository) ProductService {
	return &productService{db}
}

func (p *productService) GetAllAvailableProducts() ([]model.Product, error) {
	products, err := p.db.ListAvailableProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productService) GetProductById(productId int64) (*model.Product, error) {
	products, err := p.db.GetProductById(productId)
	if err != nil {
		return nil, err
	}

	return products, nil
}
