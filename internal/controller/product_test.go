package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	myerror "github.com/priykumar/oolio-kart-challenge/internal/error"
	"github.com/priykumar/oolio-kart-challenge/internal/model"
)

// Mock ProductService
type mockProductService struct {
	products map[int64]*model.Product
	err      error
}

func (m *mockProductService) GetAllAvailableProducts() ([]model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	var products []model.Product
	for _, p := range m.products {
		products = append(products, *p)
	}
	return products, nil
}

func (m *mockProductService) GetProductById(id int64) (*model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	if product, exists := m.products[id]; exists {
		return product, nil
	}
	return nil, myerror.KartError{Code: 404, Msg: "Product not found"}
}

// Test success
func TestGetProductHandler_Success(t *testing.T) {
	// Test success
	mockSvc := &mockProductService{
		products: map[int64]*model.Product{
			1: {Id: "1", Name: "Test Product", Price: 100.0},
		},
	}
	controller := NewProductController(mockSvc)
	req := httptest.NewRequest("GET", "/product", nil)
	w := httptest.NewRecorder()

	controller.GetProductHandler(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test service error
	mockSvc.err = myerror.KartError{Code: 500, Msg: "Internal error"}
	w = httptest.NewRecorder()
	controller.GetProductHandler(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// Test failure
func TestGetProductHandler_Failure(t *testing.T) {

	mockSvc := &mockProductService{
		err: myerror.KartError{Code: 500, Msg: "Internal error"},
	}
	controller := NewProductController(mockSvc)
	req := httptest.NewRequest("GET", "/product", nil)
	w := httptest.NewRecorder()

	controller.GetProductHandler(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// Test success
func TestGetProductByIdHandler_Success(t *testing.T) {
	mockSvc := &mockProductService{
		products: map[int64]*model.Product{
			1: {Id: "1", Name: "Test Product", Price: 100.0},
		},
	}
	controller := NewProductController(mockSvc)
	req := httptest.NewRequest("GET", "/product/1", nil)
	req = mux.SetURLVars(req, map[string]string{"productId": "1"})
	w := httptest.NewRecorder()

	controller.GetProductByIdHandler(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// Test failure
func TestGetProductByIdHandler_Failure(t *testing.T) {
	mockSvc := &mockProductService{
		products: map[int64]*model.Product{
			1: {Id: "1", Name: "Test Product", Price: 100.0},
		},
	}
	controller := NewProductController(mockSvc)
	w := httptest.NewRecorder()

	// Test invalid product ID
	req := httptest.NewRequest("GET", "/product/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"productId": "abc"})
	controller.GetProductByIdHandler(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// Test invalid product ID
	req = httptest.NewRequest("GET", "/product", nil)
	req = mux.SetURLVars(req, map[string]string{"abc": "pqr"})
	controller.GetProductByIdHandler(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
