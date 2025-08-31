package service

import (
	"github.com/priykumar/oolio-kart-challenge/internal/model"
	"github.com/priykumar/oolio-kart-challenge/internal/repo"
)

type OrderService interface {
	PlaceOrder(model.OrderDetail) (*model.OrderResp, error)
}

type orderService struct {
	db repo.KartRepository
}

func NewOrderService(db repo.KartRepository) OrderService {
	return &orderService{db}
}

func (o *orderService) PlaceOrder(oDetail model.OrderDetail) (*model.OrderResp, error) {
	// check for duplicate productIds
	pId_Count := map[string]int{}
	for _, items := range oDetail.OrderedProduct {
		pId_Count[items.ProductId] += items.Quantity
	}

	items := []model.OrderedProduct{}
	for k, v := range pId_Count {
		items = append(items, model.OrderedProduct{ProductId: k, Quantity: v})
	}
	oDetail.OrderedProduct = items

	order, err := o.db.PlaceOrder(oDetail)
	if err != nil {
		return nil, err
	}

	return order, err
}
