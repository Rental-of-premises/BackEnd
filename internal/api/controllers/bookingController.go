package api_controllers

import (
	"net/http"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/storage/repository"
	"rent/internal/api/middleware"
	"rent/internal/models"
    "encoding/json"
	"time"
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
		api_scripts.RespondError(res, http.StatusInternalServerError, err.Error())
		return
	}
	if booking == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Бронь не найдена")
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
		api_scripts.RespondError(res, http.StatusNotFound, "Брони не найдены")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, bookings)
}

func (bc *BookingController) GetMyBookings(res http.ResponseWriter, req *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}
	
	bookings, err := bc.Rep.GetByUser(userID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}
	if bookings == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Список пуст")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, bookings)
}

func (bc *BookingController) CreateBooking(res http.ResponseWriter, req *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}
	var requestBody struct {
		UserID      int64    `json:"user_id"`
		ApartmentID int64    `json:"apartment_id"`
		Status      string    `json:"status"`
		TimeFrom    time.Time `json:"time_from"`
		TimeTo      time.Time `json:"time_to"`
    }
    err := json.NewDecoder(req.Body).Decode(&requestBody)
    if err != nil {
        api_scripts.RespondError(res, http.StatusBadRequest, "Неверный JSON")
        return
    }
	
	requestBody.UserID = userID
	requestBody.Status = "waiting"
    check, err := bc.Rep.CheckAvailability(requestBody.ApartmentID, requestBody.TimeFrom, requestBody.TimeTo)
	if err != nil {
        api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
        return
    }
    if !check {
        api_scripts.RespondError(res, http.StatusBadRequest, "Брони пересекаются")
        return
    }
	
	booking := &models.Booking{
		UserID       : requestBody.UserID,
		ApartmentID  : requestBody.ApartmentID,
		Status       : requestBody.Status,
		TimeFrom     : requestBody.TimeFrom,
		TimeTo       : requestBody.TimeTo,
	}

	err = bc.Rep.Create(booking)    
	if err != nil {
        api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
        return
    }

    api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
        "id": booking.ID,
    })
}