package repo

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priykumar/oolio-kart-challenge/internal/model"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	return db
}

func TestListAvailableProducts(t *testing.T) {
	// Test success
	db := setupTestDB()
	defer db.Close()

	repo := &kartRepository{dbClient: db}
	repo.CreateTables()

	db.Exec(`DELETE FROM products`)
	// Insert test product
	db.Exec(`INSERT INTO products (name, price, category, image_thumbnail, image_mobile, image_tablet, image_desktop, is_available) 
		VALUES ('Test Product', 100.0, 'Test', 'thumb.jpg', 'mobile.jpg', 'tablet.jpg', 'desktop.jpg', 1)`)

	products, err := repo.ListAvailableProducts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(products))
	}

	// Test with unavailable products
	db.Exec(`UPDATE products SET is_available = 0`)
	products, err = repo.ListAvailableProducts()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(products) != 0 {
		t.Errorf("Expected 0 products, got %d", len(products))
	}
}

func TestGetProductById_Success(t *testing.T) {
	// Test success
	db := setupTestDB()
	defer db.Close()

	repo := &kartRepository{dbClient: db}
	repo.CreateTables()
	db.Exec(`DELETE FROM products`)

	db.Exec(`INSERT INTO products (id, name, price, category, image_thumbnail, image_mobile, image_tablet, image_desktop, is_available) 
		VALUES (1, 'Test Product', 100.0, 'Test', 'thumb.jpg', 'mobile.jpg', 'tablet.jpg', 'desktop.jpg', 1)`)

	product, err := repo.GetProductById(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if product.Name != "Test Product" {
		t.Errorf("Expected 'Test Product', got %s", product.Name)
	}

	// Test product not found
	_, err = repo.GetProductById(999)
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}
}

func TestGetProductById_Failure(t *testing.T) {
	// Test failure
	db := setupTestDB()
	defer db.Close()

	repo := &kartRepository{dbClient: db}
	repo.CreateTables()
	db.Exec(`DELETE FROM products`)

	db.Exec(`INSERT INTO products (id, name, price, category, image_thumbnail, image_mobile, image_tablet, image_desktop, is_available) 
		VALUES (1, 'Test Product', 100.0, 'Test', 'thumb.jpg', 'mobile.jpg', 'tablet.jpg', 'desktop.jpg', 1)`)

	_, err := repo.GetProductById(999)
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}
}

func TestValidateCode_Success(t *testing.T) {
	// Test valid coupon
	db := setupTestDB()
	defer db.Close()

	repo := &kartRepository{dbClient: db}
	repo.CreateTables()
	db.Exec(`DELETE FROM products`)

	db.Exec(`INSERT INTO coupons (promo_code, discount) VALUES ('SAVE10', 10.0)`)

	discount, err := repo.validateCode("SAVE10")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if discount != 10.0 {
		t.Errorf("Expected discount 10.0, got %f", discount)
	}

	// Test invalid coupon
	_, err = repo.validateCode("INVALID")
	if err == nil {
		t.Error("Expected error for invalid coupon, got nil")
	}
}

func TestValidateCode_Failure(t *testing.T) {
	// Test valid coupon
	db := setupTestDB()
	defer db.Close()

	repo := &kartRepository{dbClient: db}
	repo.CreateTables()
	db.Exec(`DELETE FROM products`)

	db.Exec(`INSERT INTO coupons (promo_code, discount) VALUES ('SAVE10', 10.0)`)

	// Test invalid coupon
	_, err := repo.validateCode("INVALID")
	if err == nil {
		t.Error("Expected error for invalid coupon, got nil")
	}
}

func TestKartRepository_PlaceOrder(t *testing.T) {
	// Test successful order
	db := setupTestDB()
	defer db.Close()

	repo := &kartRepository{dbClient: db}
	repo.CreateTables()
	db.Exec(`DELETE FROM products`)

	// Insert test data
	db.Exec(`INSERT INTO products (id, name, price, category, is_available) VALUES (1, 'Test Product', 100.0, 'Test', 1)`)
	db.Exec(`INSERT INTO coupons (promo_code, discount) VALUES ('SAVE10', 10.0)`)

	orderDetail := model.OrderDetail{
		CouponCode: "SAVE10",
		OrderedProduct: []model.OrderedProduct{
			{ProductId: "1", Quantity: 2},
		},
	}

	order, err := repo.PlaceOrder(orderDetail)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if order.Total != 180.0 { // 200 - 10% discount
		t.Errorf("Expected total 180.0, got %f", order.Total)
	}

	// Test invalid product
	orderDetail.OrderedProduct[0].ProductId = "999"
	_, err = repo.PlaceOrder(orderDetail)
	if err == nil {
		t.Error("Expected error for invalid product, got nil")
	}
}
