package repo

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func (k *kartRepository) CreateTables() error {
	// Create Product table
	productCmd := `
	CREATE TABLE IF NOT EXISTS products 
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		category TEXT NOT NULL,
		image_thumbnail TEXT,
		image_mobile TEXT,
		image_tablet TEXT,
		image_desktop TEXT,
		is_available INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := k.dbClient.Exec(productCmd)
	if err != nil {
		fmt.Println("Failed creating table products. Error: ", err)
		return err
	}

	// Create Order table
	orderCmd := `
	CREATE TABLE IF NOT EXISTS orders 
	(
		id TEXT PRIMARY KEY,
		total REAL NOT NULL,
		discounts REAL DEFAULT 0.0,
		coupon_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (coupon_id) REFERENCES coupons(id)
	);
	`
	_, err = k.dbClient.Exec(orderCmd)
	if err != nil {
		fmt.Println("Failed creating table orders. Error: ", err)
		return err
	}

	// Table to map order to products in order
	orderItemsCmd := `
	CREATE TABLE IF NOT EXISTS order_items 
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id TEXT NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	);
	`
	_, err = k.dbClient.Exec(orderItemsCmd)
	if err != nil {
		fmt.Println("Failed creating table order_items. Error: ", err)
		return err
	}

	// Table to map promocode to discount
	couponCmd := `
	CREATE TABLE IF NOT EXISTS coupons 
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		promo_code TEXT UNIQUE NOT NULL,
		discount REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = k.dbClient.Exec(couponCmd)
	if err != nil {
		fmt.Println("Failed creating table coupons. Error: ", err)
		return err
	}

	fmt.Println("All the tables are successully created")

	// Populate products table if no data found in it
	var count int
	k.dbClient.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		k.PopulateTables()
	}
	return nil
}

func (k *kartRepository) PopulateTables() {
	var products = map[string][]string{
		"Waffle":    {"Chicken Waffle", "Banana Waffle", "Belgian Waffle", "Chocolate Waffle", "Red Velvet Waffle"},
		"Pancakes":  {"Classic Pancakes", "Blueberry Pancakes", "Chocolate Chip Pancakes", "Banana Pancakes"},
		"Burger":    {"Classic Burger", "Cheese Burger", "Mushroom Burger", "Veggie Burger", "Chicken Burger"},
		"Pizza":     {"Margherita Pizza", "Farmhouse Pizza", "Panner Makhani Pizza", "Onion Pizza", "BBQ Chicken Pizza"},
		"Pasta":     {"Arabiatta", "Arabiatta", "Alio Olio", "Mac and Cheese", "Chicken Pasta"},
		"Beverages": {"Fresh Lime Juice", "Iced Coffee", "Hot Chocolate", "Green Tea", "Milkshake"},
	}

	var dish_baseprice_multipler = map[string][]float64{
		"Waffle":    {110, 1.2},
		"Pancakes":  {200, 1.33},
		"Burger":    {180, 1.24},
		"Pizza":     {310, 1.5},
		"Pasta":     {400, 1.08},
		"Beverages": {100, 1.38},
	}

	stmt := `INSERT INTO products (name, price, category, image_thumbnail, image_mobile, image_tablet, image_desktop, is_available) 
		VALUES ('%s', %.2f, '%s', '%s', '%s', '%s', '%s', %d)`

	rand.Seed(time.Now().UnixNano())
	baseurl := "https://orderfoodonline.deno.dev/public/images/"
	url := ""

	// Populate Product Table
	for category, dishes := range products {
		base := dish_baseprice_multipler[category][0]
		mul := dish_baseprice_multipler[category][1]
		for _, dish := range dishes {
			url = baseurl + strings.ReplaceAll(strings.ToLower(dish), " ", "-")
			randomBit := rand.Intn(2)
			s := fmt.Sprintf(stmt, dish, base*mul, category, url+"-thumbnail.jpg", url+"-mobile.jpg", url+"-tablet.jpg", url+"-desktop.jpg", randomBit)
			fmt.Println("Executing", s)
			_, err := k.dbClient.Exec(s)
			if err != nil {
				fmt.Println("Failed inserting into products. Error:", err)
			}
		}
	}

}

func (k *kartRepository) PopulateCoupons(filePath string) {
	fmt.Println("Populating coupons in DB")
	rand.Seed(time.Now().UnixNano())

	filePath = "../token/" + filePath
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening valid token file:", err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code := scanner.Text()
		if code == "" {
			continue
		}

		discount := 10 + rand.Float64()*(50-10)
		discount = float64(int(discount*100)) / 100 // truncate to 2 decimals

		_, err := k.dbClient.Exec("INSERT OR IGNORE INTO coupons (promo_code, discount) VALUES (?, ?)", code, discount)
		if err != nil {
			log.Println("Insert error:", err)
		}
	}
	fmt.Println("Done populating coupons in DB")
}
