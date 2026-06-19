package api_scripts

import (
	"encoding/json"
	"fmt"
	"net/http"
	api_models "rent/internal/api/models"
	"strconv"
	models "rent/internal/models"

	"github.com/gorilla/mux"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, &api_models.ErrorResponse{Error: message})
}

func ParseID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Неверный ID")
	}
	return id, nil
}

func ParseApartmentFilter(r *http.Request) (*models.ApartmentFilter, error) {
	filter := &models.ApartmentFilter{}
	
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			return nil, fmt.Errorf("invalid is_active value, must be true/false")
		}
		filter.IsActive = &isActive
	}
	
	if sellerIDStr := r.URL.Query().Get("seller_id"); sellerIDStr != "" {
		sellerID, err := strconv.ParseInt(sellerIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid seller_id value, must be integer")
		}
		filter.SellerID = &sellerID
	}
	
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseInt(minPriceStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min_price value, must be integer")
		}
		filter.MinPrice = &minPrice
	}
	
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max_price value, must be integer")
		}
		filter.MaxPrice = &maxPrice
	}
	
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid limit value, must be integer")
		}
		if parsed > 0 {
			limit = parsed
		}
	}
	filter.Limit = &limit
	
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsed, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid offset value, must be integer")
		}
		if parsed >= 0 {
			offset = parsed
		}
	}
	filter.Offset = &offset

	amenityIDs := r.URL.Query()["amenities"]
	if len(amenityIDs) > 0 {
		for _, idStr := range amenityIDs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err == nil {
				filter.Amenities = append(filter.Amenities, id)
			}
		}
	}
	
	return filter, nil
}

func ParseBookingFilter(r *http.Request) (*models.BookingFilter, error) {
	filter := &models.BookingFilter{}
	
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		filter.Status = &statusStr
	}
	
	if sellerIDStr := r.URL.Query().Get("seller_id"); sellerIDStr != "" {
		sellerID, err := strconv.ParseInt(sellerIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid seller_id value, must be integer")
		}
		filter.SellerID = &sellerID
	}
	
	if apartmentIDStr := r.URL.Query().Get("apartment_id"); apartmentIDStr != "" {
		apartmentID, err := strconv.ParseInt(apartmentIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid apartment_id value, must be integer")
		}
		filter.ApartmentID = &apartmentID
	}
	
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id value, must be integer")
		}
		filter.UserID = &userID
	}
	
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseInt(minPriceStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min_price value, must be integer")
		}
		filter.MinPrice = &minPrice
	}
	
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max_price value, must be integer")
		}
		filter.MaxPrice = &maxPrice
	}
	
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid limit value, must be integer")
		}
		if parsed > 0 {
			limit = parsed
		}
	}
	filter.Limit = &limit
	
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsed, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid offset value, must be integer")
		}
		if parsed >= 0 {
			offset = parsed
		}
	}
	filter.Offset = &offset
	
	return filter, nil
}
