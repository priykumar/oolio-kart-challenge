package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	myerror "github.com/priykumar/oolio-kart-challenge/internal/error"
	"github.com/priykumar/oolio-kart-challenge/internal/model"
	"github.com/priykumar/oolio-kart-challenge/internal/service"
)

type OrderController struct {
	svc service.OrderService
}

func NewOrderController(svc service.OrderService) *OrderController {
	return &OrderController{svc}
}

func validateOrder(oDetail model.OrderDetail) error {
	if len(oDetail.OrderedProduct) == 0 {
		return fmt.Errorf("no product provided")
	} else {
		for _, od := range oDetail.OrderedProduct {
			if strings.TrimSpace(od.ProductId) == "" {
				return fmt.Errorf("product id not present in request")
			} else if od.Quantity <= 0 {
				return fmt.Errorf("quantity can't be negative or zero")
			}
		}
	}

	return nil
}

func (o *OrderController) PlaceOrderHandler(w http.ResponseWriter, r *http.Request) {
	var oDetail model.OrderDetail

	if r.Body == nil {
		generateResponse(w, myerror.KartError{Code: 400, Msg: "No request body found"})
		return
	}

	json.NewDecoder(r.Body).Decode(&oDetail)
	if err := validateOrder(oDetail); err != nil {
		generateResponse(w, myerror.KartError{Code: 400, Msg: err.Error()})
		return
	}

	orders, err := o.svc.PlaceOrder(oDetail)
	if err != nil {
		generateResponse(w, err.(myerror.KartError))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
