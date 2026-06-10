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
		return 0, fmt.Errorf("invalid ID")
	}
	return id, nil
}

func ParseApartmentFilter(r *http.Request) (*models.ApartmentFilter, error) {
	var filter models.ApartmentFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		return nil, fmt.Errorf("invalid filter")
	}
	return &filter, nil
}

func ParseBookingFilter(r *http.Request) (*models.BookingFilter, error) {
	var filter models.BookingFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		return nil, fmt.Errorf("invalid filter")
	}

	return &filter, nil
}

