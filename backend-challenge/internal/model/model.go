package model

type Response struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}

type Product struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Image    Image   `json:"image"`
}

type OrderedProduct struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderDetail struct {
	CouponCode     string           `json:"couponCode"`
	OrderedProduct []OrderedProduct `json:"items"`
}

type OrderResp struct {
	Id             string           `json:"id"`
	Total          float64          `json:"total"`
	Discount       float64          `json:"discounts"`
	OrderedProduct []OrderedProduct `json:"items"`
}
