package api_controllers

import (
	"net/http"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/storage/repository"
)

type BookingController struct {
	Rep *repository.BookingRepository
}

func (bc *BookingController) GetBooking(res http.ResponseWriter, req *http.Request) {
	id, err := api_scripts.ParseID(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := bc.Rep.GetByID(id)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Failed to get booking")
		return
	}
	if booking == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Booking not found")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, booking)
}
	
func (bc *BookingController) GetAllBookings(res http.ResponseWriter, req *http.Request) {
	filter, err := api_scripts.ParseBookingFilter(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	bookings, err := bc.Rep.GetAll(filter)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}
	if bookings == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Apartment not found")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, bookings)
}
