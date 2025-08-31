package repo

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	myerror "github.com/priykumar/oolio-kart-challenge/internal/error"
	"github.com/priykumar/oolio-kart-challenge/internal/model"
)

var mu = &sync.Mutex{}
var repo *kartRepository

type KartRepository interface {
	ListAvailableProducts() ([]model.Product, error)
	GetProductById(int64) (*model.Product, error)
	PlaceOrder(model.OrderDetail) (*model.OrderResp, error)
	PopulateCoupons(string)
}

type kartRepository struct {
	dbClient *sql.DB
}

func getDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "../repo/mydb.db")
	if err != nil || db == nil {
		fmt.Println("Error while opening db driver for sql-lite. Error: ", err)
		panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Error while connecting to sqlite db:", err)
		panic(err)
	}

	fmt.Println("Successfully initialized driver for my-sql")
	return db
}

func InitialiseDatabase() KartRepository {
	// singleton design pattern
	if repo == nil || repo.dbClient == nil {
		mu.Lock()
		defer mu.Unlock()
		if repo == nil || repo.dbClient == nil {
			repo = &kartRepository{
				dbClient: getDatabase(),
			}

			repo.CreateTables()
		}
	}

	return repo
}

// Get list of available products
func (k *kartRepository) ListAvailableProducts() ([]model.Product, error) {
	cmd := `SELECT id, name, price, category, image_thumbnail, image_mobile, image_tablet, image_desktop
	FROM products WHERE is_available=1`

	rows, err := k.dbClient.Query(cmd)
	if err != nil {
		fmt.Println("Failed quering products table for available products. Error:", err)
		return nil, myerror.KartError{Code: 500, Msg: "Failed quering DB"}
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var thumb, mobile, tablet, desktop string
		var id int
		err := rows.Scan(
			&id,
			&p.Name,
			&p.Price,
			&p.Category,
			&thumb,
			&mobile,
			&tablet,
			&desktop,
		)
		if err != nil {
			fmt.Println("Failed scanning rows. Error:", err)
			return nil, myerror.KartError{Code: 500, Msg: "Failed scanning rows in DB"}
		}

		// Image details
		p.Image = model.Image{
			Thumbnail: thumb,
			Mobile:    mobile,
			Tablet:    tablet,
			Desktop:   desktop,
		}

		p.Id = fmt.Sprintf("%d", id)
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Failed scanning rows. Error:", err)
		return nil, myerror.KartError{Code: 500, Msg: "Failed scanning rows in DB"}
	}

	return products, nil
}

// Get single product by ID
func (k *kartRepository) GetProductById(productId int64) (*model.Product, error) {
	cmd := `SELECT id, name, price, category, image_thumbnail, image_mobile, image_tablet, image_desktop
	FROM products WHERE id = ? AND is_available=1`

	var p model.Product
	var thumb, mobile, tablet, desktop string
	var id int
	err := k.dbClient.QueryRow(cmd, productId).Scan(
		&id,
		&p.Name,
		&p.Price,
		&p.Category,
		&thumb,
		&mobile,
		&tablet,
		&desktop,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Product not found or not available: ID", productId)
			return nil, myerror.KartError{Code: 404, Msg: "Product not found or not available"}
		}
		fmt.Printf("Failed querying DB for product ID %d. Error: %v\n", productId, err)
		return nil, myerror.KartError{Code: 500, Msg: "Failed quering DB"}
	}

	// Image details
	p.Image = model.Image{
		Thumbnail: thumb,
		Mobile:    mobile,
		Tablet:    tablet,
		Desktop:   desktop,
	}

	p.Id = fmt.Sprintf("%d", id)
	return &p, nil
}

func (k *kartRepository) validateCode(promo string) (float64, error) {
	var discount float64 = 0
	err := k.dbClient.QueryRow("SELECT discount FROM coupons WHERE promo_code = ?", promo).Scan(&discount)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Invalid coupon code is provided")
			return 0.0, myerror.KartError{Code: 400, Msg: "Invalid coupon code is provided"}
		}
		fmt.Println("Failed to check promo code in coupon table. Error:", err)
		return 0.0, myerror.KartError{Code: 500, Msg: "Failed quering DB"}
	}

	return discount, nil
}

// Place order
func (k *kartRepository) PlaceOrder(oDetail model.OrderDetail) (order *model.OrderResp, err error) {
	var discountPercent float64 = 0
	if oDetail.CouponCode != "" {
		discountPercent, err = k.validateCode(oDetail.CouponCode)
		if err != nil {
			fmt.Println("failed validating coupon")
			return nil, err
		}
	}

	// begin the transaction
	tx, err := k.dbClient.Begin()
	if err != nil {
		fmt.Println("Failed to begin transaction. Error:", err)
		return nil, myerror.KartError{Code: 500, Msg: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	orderID := uuid.New().String()

	// Prepare statement for order items
	stmt, err := tx.Prepare(`INSERT INTO order_items (order_id, product_id, quantity) VALUES (?, ?, ?)`)
	if err != nil {
		fmt.Println("failed to prepare statement to be executed")
		return nil, myerror.KartError{Code: 500, Msg: "Failed to prepare statement"}
	}
	defer stmt.Close()

	// Insert all order items
	for _, item := range oDetail.OrderedProduct {
		// Validate product exists and is available
		var exists int
		err = tx.QueryRow("SELECT COUNT(*) FROM products WHERE id = ? AND is_available = 1", item.ProductId).Scan(&exists)
		if err != nil {
			fmt.Println("Failed to validate product. Error:", err)
			return nil, myerror.KartError{Code: 500, Msg: "Failed to validate product"}
		}
		if exists == 0 {
			fmt.Printf("Product %s not found or not available\n", item.ProductId)
			return nil, myerror.KartError{Code: 400, Msg: "Provided product is not valid or is not available"}
		}

		// Insert order item
		_, err = stmt.Exec(orderID, item.ProductId, item.Quantity)
		if err != nil {
			fmt.Println("Failed to execute transaction. Error", err)
			return nil, myerror.KartError{Code: 500, Msg: "Failed to execute transaction"}
		}

		fmt.Printf("Added item: Product %s, Quantity %d\n", item.ProductId, item.Quantity)
	}

	total, err := calculateOrderTotal(tx, oDetail.OrderedProduct)
	if err != nil {
		return nil, myerror.ErrInternalServer
	}
	discount := total * (discountPercent / 100.0)
	finalTotal := total - discount

	discount = float64(int(discount*100)) / 100
	finalTotal = float64(int(finalTotal*100)) / 100

	// Insert main order
	_, err = tx.Exec(`INSERT INTO orders (id, total, discounts, coupon_id) VALUES (?, ?, ?, (SELECT id FROM coupons WHERE promo_code = ?))`,
		orderID, finalTotal, discount, oDetail.CouponCode)
	if err != nil {
		fmt.Println("Failed inserting order detail. Error:", err)
		return nil, myerror.KartError{Code: 500, Msg: "Failed inserting into DB"}
	}

	// Commit transaction - all or nothing
	if err = tx.Commit(); err != nil {
		fmt.Println("Failed to commit transaction. Error:", err)
		return nil, myerror.KartError{Code: 500, Msg: "Failed to commit transaction"}
	}

	// Return created order
	order = &model.OrderResp{
		Id:             orderID,
		Total:          finalTotal,
		Discount:       discount,
		OrderedProduct: oDetail.OrderedProduct,
	}

	fmt.Printf("Order created successfully: %s (Total: %.2f)\n", orderID, finalTotal)
	return order, nil
}

func calculateOrderTotal(tx *sql.Tx, items []model.OrderedProduct) (float64, error) {
	var total float64 = 0.0

	for _, item := range items {
		var price float64
		err := tx.QueryRow("SELECT price FROM products WHERE id = ?", item.ProductId).Scan(&price)
		if err != nil {
			fmt.Println("failed getting price for product", item.ProductId, "Error: ", err)
			return 0, myerror.KartError{Code: 500, Msg: "Failed to get price of product"}
		}
		total += price * float64(item.Quantity)
	}

	return total, nil
}
