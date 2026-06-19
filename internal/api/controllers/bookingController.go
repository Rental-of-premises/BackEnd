package api_controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"rent/internal/api/middleware"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/models"
	"rent/internal/storage/repository"
)

type BookingController struct {
	Rep           *repository.BookingRepository
	ApartmentRepo *repository.ApartmentRepository
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

func (bc *BookingController) GetBookingsByApartment(res http.ResponseWriter, req *http.Request) {
	apartmentID, err := api_scripts.ParseID(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	allBookings, err := bc.Rep.GetByApartment(apartmentID, limit, offset)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, err.Error())
		return
	}

	var activeBookings []*models.Booking
	for _, b := range allBookings {
		if b.Status == "waiting" || b.Status == "confirmed" || b.Status == "completed" {
			activeBookings = append(activeBookings, b)
		}
	}

	if activeBookings == nil {
		activeBookings = []*models.Booking{}
	}

	api_scripts.RespondJSON(res, http.StatusOK, activeBookings)
}

func (bc *BookingController) GetMyBookings(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}

	statusFilter := req.URL.Query().Get("statusFilter")
	if statusFilter == "" {
		statusFilter = "active"
	}

	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	filter := &models.BookingFilter{
		UserID: &userID,
		Limit:  &limit,
		Offset: &offset,
	}

	var allBookings []*models.Booking
	var statuses []string

	switch statusFilter {
	case "active":
		statuses = append(statuses, "confirmed", "waiting")
	case "history":
		statuses = append(statuses, "cancelled", "completed", "rejected")
	case "pending":
		statuses = append(statuses, "waiting")
	case "all":
		statuses = append(statuses, "confirmed", "cancelled", "completed", "waiting", "rejected")
	default:
		statuses = append(statuses, "confirmed", "waiting")
	}

	for _, status := range statuses {
		filter.Status = &status
		bookings, err := bc.Rep.GetAll(filter)
		if err != nil {
			api_scripts.RespondError(res, http.StatusInternalServerError, err.Error())
			return
		}
		allBookings = append(allBookings, bookings...)
	}

	if allBookings == nil {
		allBookings = []*models.Booking{}
	}

	api_scripts.RespondJSON(res, http.StatusOK, allBookings)
}

func (bc *BookingController) GetBookings(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}

	statusFilter := req.URL.Query().Get("statusFilter")
	if statusFilter == "" {
		statusFilter = "active"
	}

	limit := 9999999
	offset := 0

	filter := &models.BookingFilter{
		SellerID: &userID,
		Limit:    &limit,
		Offset:   &offset,
	}

	var allBookings []*models.Booking
	var statuses []string

	switch statusFilter {
	case "active":
		statuses = append(statuses, "confirmed", "waiting")
	case "history":
		statuses = append(statuses, "cancelled", "completed", "rejected")
	case "pending":
		statuses = append(statuses, "waiting")
	case "all":
		statuses = append(statuses, "confirmed", "cancelled", "completed", "waiting", "rejected")
	default:
		statuses = append(statuses, "confirmed", "waiting")
	}

	for _, status := range statuses {
		filter.Status = &status
		bookings, err := bc.Rep.GetAll(filter)
		if err != nil {
			api_scripts.RespondError(res, http.StatusInternalServerError, err.Error())
			return
		}
		allBookings = append(allBookings, bookings...)
	}

	if allBookings == nil {
		allBookings = []*models.Booking{}
	}

	api_scripts.RespondJSON(res, http.StatusOK, allBookings)
}

func (bc *BookingController) CreateBooking(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}
	var requestBody struct {
		UserID      int64     `json:"user_id"`
		ApartmentID int64     `json:"apartment_id"`
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
	check, err := bc.Rep.CheckAvailability(requestBody.ApartmentID, requestBody.TimeFrom, requestBody.TimeTo, userID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Ошибка при проверке возможности брони: " + err.Error())
		return
	}
	if !check {
		api_scripts.RespondError(res, http.StatusBadRequest, "Нельзя забронировать")
		return
	}

	booking := &models.Booking{
		UserID:      requestBody.UserID,
		ApartmentID: requestBody.ApartmentID,
		Status:      requestBody.Status,
		TimeFrom:    requestBody.TimeFrom,
		TimeTo:      requestBody.TimeTo,
	}

	err = bc.Rep.Create(booking)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Ошибка при создании брони: " + err.Error())
		return
	}

	api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
		"id": booking.ID,
	})
}

func (bc *BookingController) CancelBooking(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	bookingID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := bc.Rep.GetByID(bookingID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске бронирования")
		return
	}
	if booking == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Бронирование не найдено")
		return
	}

	if booking.UserID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "У вас нет прав на отмену этого бронирования")
		return
	}

	if booking.Status != "waiting" && booking.Status != "confirmed" {
		api_scripts.RespondError(res, http.StatusBadRequest, "Нельзя отменить бронирование в статусе "+booking.Status)
		return
	}

	err = bc.Rep.UpdateStatus(bookingID, "cancelled")
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при отмене бронирования")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Бронирование успешно отменено",
		"status":  "cancelled",
	})
}

func (bc *BookingController) ConfirmBookingBySeller(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	bookingID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := bc.Rep.GetByID(bookingID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске бронирования")
		return
	}
	if booking == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Бронирование не найдено")
		return
	}

	apartment, err := bc.ApartmentRepo.GetByID(booking.ApartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске квартиры")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Квартира не найдена")
		return
	}

	if apartment.SellerID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "Только продавец может подтвердить бронирование")
		return
	}

	if booking.Status != "waiting" {
		api_scripts.RespondError(res, http.StatusBadRequest, "Нельзя подтвердить бронирование в статусе "+booking.Status)
		return
	}

	err = bc.Rep.UpdateStatus(bookingID, "confirmed")
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при подтверждении бронирования")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Бронирование успешно подтверждено",
		"status":  "confirmed",
	})
}

func (bc *BookingController) RejectBookingBySeller(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	bookingID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := bc.Rep.GetByID(bookingID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске бронирования")
		return
	}
	if booking == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Бронирование не найдено")
		return
	}

	apartment, err := bc.ApartmentRepo.GetByID(booking.ApartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске квартиры")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Квартира не найдена")
		return
	}

	if apartment.SellerID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "Только продавец может отклонить бронирование")
		return
	}

	if booking.Status != "waiting" {
		api_scripts.RespondError(res, http.StatusBadRequest, "Нельзя отклонить бронирование в статусе "+booking.Status)
		return
	}

	err = bc.Rep.UpdateStatus(bookingID, "rejected")
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при отклонении бронирования")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Бронирование отклонено",
		"status":  "rejected",
	})
}
