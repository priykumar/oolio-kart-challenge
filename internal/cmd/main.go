package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/priykumar/oolio-kart-challenge/internal/controller"
	"github.com/priykumar/oolio-kart-challenge/internal/middleware"
	"github.com/priykumar/oolio-kart-challenge/internal/repo"
	"github.com/priykumar/oolio-kart-challenge/internal/service"
)

type Config struct {
	CouponArtifacts []string `json:"coupon_artifacts"`
	ValidTokenPath  string   `json:"valid_token_path"`
}

func isTokenFileEmpty(filePath string) bool {
	filePath = "../token/" + filePath

	stat, err := os.Stat(filePath)
	if err != nil || stat.Size() == 0 {
		return true
	}

	return false
}

func main() {
	db := repo.InitialiseDatabase()
	psvc := service.NewProductService(db)
	osvc := service.NewOrderService(db)
	p := controller.NewProductController(psvc)
	s := controller.NewOrderController(osvc)

	// read and parse config file
	configFile, err := os.Open("../config/config.json")
	if err != nil {
		panic(fmt.Errorf("failed to open config.json: %w", err))
	}
	defer configFile.Close()

	var cfg Config
	if err := json.NewDecoder(configFile).Decode(&cfg); err != nil {
		panic(fmt.Errorf("failed to decode config.json: %w", err))
	}

	fmt.Println("Coupon artifacts:", cfg.CouponArtifacts)

	if isEmpty := isTokenFileEmpty(cfg.ValidTokenPath); isEmpty {
		fmt.Println("Token files are empty, hence read token artifacts")
		readArtifacts(cfg.CouponArtifacts)
	}
	db.PopulateCoupons(cfg.ValidTokenPath)

	r := mux.NewRouter()
	r.HandleFunc("/product", p.GetProductHandler).Methods("GET")
	r.HandleFunc("/product/{productId}", p.GetProductByIdHandler).Methods("GET")
	r.Handle("/order", middleware.ApiKeyMiddleware(http.HandlerFunc(s.PlaceOrderHandler))).Methods("POST")

	http.ListenAndServe(":8080", r)
}
