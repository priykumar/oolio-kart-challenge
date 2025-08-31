package service

import (
	"testing"

	myerror "github.com/priykumar/oolio-kart-challenge/internal/error"
	"github.com/priykumar/oolio-kart-challenge/internal/model"
)

// Mock Repository
type mockKartRepository struct {
	products map[int64]*model.Product
	order    *model.OrderResp
	err      error
}

func (m *mockKartRepository) ListAvailableProducts() ([]model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	var products []model.Product
	for _, p := range m.products {
		products = append(products, *p)
	}
	return products, nil
}

func (m *mockKartRepository) GetProductById(id int64) (*model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	if product, exists := m.products[id]; exists {
		return product, nil
	}
	return nil, myerror.KartError{Code: 404, Msg: "Product not found"}
}

func (m *mockKartRepository) PlaceOrder(oDetail model.OrderDetail) (*model.OrderResp, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.order, nil
}

func (m *mockKartRepository) PopulateCoupons(string) {}

// GetAllAvailableProducts Success Tests
func TestGetAllAvailableProducts_Success(t *testing.T) {
	mockRepo := &mockKartRepository{
		products: map[int64]*model.Product{
			1: {Id: "1", Name: "Test Product", Price: 100.0},
		},
	}
	svc := NewProductService(mockRepo)

	products, err := svc.GetAllAvailableProducts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(products))
	}
}

// GetAllAvailableProducts Failure Test
func TestGetAllAvailableProducts_Failure(t *testing.T) {
	// Test success
	mockRepo := &mockKartRepository{
		err: myerror.KartError{Code: 500, Msg: "DB error"},
	}
	svc := NewProductService(mockRepo)

	_, err := svc.GetAllAvailableProducts()
	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
}

func TestGetProductById_Success(t *testing.T) {
	// Test success
	mockRepo := &mockKartRepository{
		products: map[int64]*model.Product{
			1: {Id: "1", Name: "Test Product", Price: 100.0},
		},
	}
	svc := NewProductService(mockRepo)

	product, err := svc.GetProductById(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if product.Id != "1" {
		t.Errorf("Expected product ID '1', got %s", product.Id)
	}

	// Test product not found
	_, err = svc.GetProductById(999)
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}
}

func TestGetProductById_Failure(t *testing.T) {
	mockRepo := &mockKartRepository{
		products: map[int64]*model.Product{
			1: {Id: "1", Name: "Test Product", Price: 100.0},
		},
	}
	svc := NewProductService(mockRepo)

	// Test product not found
	_, err := svc.GetProductById(999)
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}
}

// OrderService Tests
func TestPlaceOrder_Success(t *testing.T) {
	// Test success with duplicate product consolidation
	mockRepo := &mockKartRepository{
		order: &model.OrderResp{Id: "order-123", Total: 300.0},
	}
	svc := NewOrderService(mockRepo)

	orderDetail := model.OrderDetail{
		OrderedProduct: []model.OrderedProduct{
			{ProductId: "1", Quantity: 2},
		},
	}

	order, err := svc.PlaceOrder(orderDetail)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if order.Id != "order-123" {
		t.Errorf("Expected order ID 'order-123', got %s", order.Id)
	}
}

func TestPlaceOrder_Failure(t *testing.T) {
	// Test success with duplicate product consolidation
	mockRepo := &mockKartRepository{
		err: myerror.KartError{Code: 400, Msg: "Invalid product"},
	}
	svc := NewOrderService(mockRepo)

	orderDetail := model.OrderDetail{
		OrderedProduct: []model.OrderedProduct{
			{ProductId: "1", Quantity: 2},
		},
	}
	_, err := svc.PlaceOrder(orderDetail)
	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
}
