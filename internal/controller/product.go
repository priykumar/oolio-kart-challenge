package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	myerror "github.com/priykumar/oolio-kart-challenge/internal/error"
	"github.com/priykumar/oolio-kart-challenge/internal/model"
	"github.com/priykumar/oolio-kart-challenge/internal/service"
)

type ProductController struct {
	svc service.ProductService
}

func NewProductController(svc service.ProductService) *ProductController {
	return &ProductController{svc}
}

// Get all the available products
func (p *ProductController) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	products, err := p.svc.GetAllAvailableProducts()
	if err != nil {
		generateResponse(w, err.(myerror.KartError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(products)
}

// Get product by productId
func (p *ProductController) GetProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if productId exists in path parameter
	if _, exist := vars["productId"]; !exist {
		fmt.Println("No product Id provided")
		generateResponse(w, myerror.KartError{Code: 400, Msg: "No product Id provided"})
		return
	}

	productId := vars["productId"]

	// Validate product ID
	pId, err := strconv.Atoi(productId)
	if err != nil || pId < 0 {
		fmt.Println("Invalid product Id provided")
		generateResponse(w, myerror.KartError{Code: 400, Msg: "Invalid ID supplied"})
		return
	}

	// Get product from service
	products, err := p.svc.GetProductById(int64(pId))
	if err != nil {
		generateResponse(w, err.(myerror.KartError))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func generateResponse(w http.ResponseWriter, kErr myerror.KartError) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(kErr.Code)
	json.NewEncoder(w).Encode(
		model.Response{
			Code:    int32(kErr.Code),
			Type:    myerror.Code2Err[kErr.Code].Error(),
			Message: kErr.Msg,
		},
	)
}
