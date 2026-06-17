package api_controllers

import (
	"encoding/json"
	"net/http"
	"strconv"


	"rent/internal/api/middleware"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/models"
	"rent/internal/storage/repository"
)

type ReviewController struct {
	Rep           *repository.ReviewRepository
	ApartmentRepo *repository.ApartmentRepository
	BookingRepo *repository.BookingRepository
}

func (rc *ReviewController) GetReview(res http.ResponseWriter, req *http.Request) {
	id, err := api_scripts.ParseID(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	review, err := rc.Rep.GetByID(id)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, err.Error())
		return
	}
	if review == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Отзыв не найден")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, review)
}

func (rc *ReviewController) GetAllReviews(res http.ResponseWriter, req *http.Request) {
	id, err := api_scripts.ParseID(req)

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

	allReviews, err := rc.Rep.GetByApartment(id, limit, offset)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, err.Error())
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, allReviews)
}

func (rc *ReviewController) CreateReview(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}	
	apartmentID, err := api_scripts.ParseID(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}
	var requestBody struct {
		UserID      int64     
		ApartmentID int64   
		Stars       int16    `json:"stars"`
		Comment     string    `json:"comment"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Неверный JSON")
		return
	}

	// valid, err := rc.BookingRepo.CheckReviewValidation(requestBody.UserID, requestBody.ApartmentID)
	// if(err != nil) {
	// 	api_scripts.RespondError(res, http.StatusUnauthorized, "Ошибка при проверке валидности написания отзыва: " + err.Error())
	// 	return
	// }
	// if(!valid) {
	// 	api_scripts.RespondError(res, http.StatusUnauthorized, "Нельзя писать отзыв к непроверенным объявлениям")
	// 	return
	// }

	requestBody.UserID = userID
	requestBody.ApartmentID = apartmentID

	review := &models.Review{
		UserID:      &requestBody.UserID,
		ApartmentID: &requestBody.ApartmentID,
		Stars:      &requestBody.Stars,
		Comment:    &requestBody.Comment,
	}

	err = rc.Rep.Create(review)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
		"id": review.ID,
	})
}

func (rc *ReviewController) DeleteReview(res http.ResponseWriter, req *http.Request) {
    userID, ok := middleware.GetUserIDFromContext(req)
    if !ok {
        api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
        return
    }

    err := rc.Rep.Delete(userID)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при удалении отзыва: " + err.Error())
        return
    }

    api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
        "message": "Отзыв успешно удален",
    })
}