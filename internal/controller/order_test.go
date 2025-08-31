package controller

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/priykumar/oolio-kart-challenge/internal/model"
)

// Mock OrderService
type mockOrderService struct {
	order *model.OrderResp
	err   error
}

func (m *mockOrderService) PlaceOrder(oDetail model.OrderDetail) (*model.OrderResp, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.order, nil
}

func TestValidateOrder_Success(t *testing.T) {
	// Test valid order
	validOrder := model.OrderDetail{
		OrderedProduct: []model.OrderedProduct{
			{ProductId: "1", Quantity: 2},
		},
	}
	if err := validateOrder(validOrder); err != nil {
		t.Errorf("Expected no error for valid order, got: %v", err)
	}
}

func TestValidateOrder_Failure(t *testing.T) {
	// Test invalid order - empty products
	invalidOrder := model.OrderDetail{OrderedProduct: []model.OrderedProduct{}}
	if err := validateOrder(invalidOrder); err == nil {
		t.Error("Expected error for empty products, got nil")
	}
}

// Test Success
func TestPlaceOrderHandler_Success(t *testing.T) {
	// Test success
	mockSvc := &mockOrderService{
		order: &model.OrderResp{Id: "order-123", Total: 200.0},
	}
	controller := NewOrderController(mockSvc)
	orderDetail := model.OrderDetail{
		OrderedProduct: []model.OrderedProduct{
			{ProductId: "1", Quantity: 2},
		},
	}
	body, _ := json.Marshal(orderDetail)
	req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controller.PlaceOrderHandler(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test no request body
	req = httptest.NewRequest("POST", "/order", nil)
	w = httptest.NewRecorder()

	controller.PlaceOrderHandler(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test Failure
func TestPlaceOrderHandler_Failure(t *testing.T) {
	mockSvc := &mockOrderService{
		order: &model.OrderResp{},
	}
	controller := NewOrderController(mockSvc)

	// Test no request body
	req := httptest.NewRequest("POST", "/order", nil)
	w := httptest.NewRecorder()

	controller.PlaceOrderHandler(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
