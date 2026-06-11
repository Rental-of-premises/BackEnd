package models

import (
	"time"
)

type Apartment struct {
	ID           int64     `json:"id"`
	SellerID     int64    `json:"seller_id"`
	Name         string    `json:"name"`
	Description  string   `json:"description"`
	Capacity     int16    `json:"capacity"`
	PricePerHour int64    `json:"price_per_hour"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type ApartmentFilter struct {
    IsActive *bool `json:"is_active"`
    SellerID *int64  `json:"seller_id"`
    MinPrice *int64  `json:"min_price"`
    MaxPrice *int64  `json:"max_price"`
    Limit *int  `json:"limit"`
	Offset *int  `json:"offset"` 
}