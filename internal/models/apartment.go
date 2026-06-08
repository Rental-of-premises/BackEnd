package models

import (
	"time"
)

type Apartment struct {
	ID           int64     `json:"id"`
	SellerID     *int64    `json:"seller_id,omitempty"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	Capacity     *int16    `json:"capacity,omitempty"`
	PricePerHour *int32    `json:"price_per_hour,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

