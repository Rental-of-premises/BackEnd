package api_controllers

import (
	"net/http"
    "encoding/json"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/storage/repository"
	"rent/internal/models"
	"rent/internal/api/middleware"
    "rent/internal/api/utils"
)

type ApartmentController struct {
	Rep *repository.ApartmentRepository
}

const(
	inf = 9999999
)

func (ac *ApartmentController) GetApartment(res http.ResponseWriter, req *http.Request) {
	id, err := api_scripts.ParseID(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	apartment, err := ac.Rep.GetByID(id)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Failed to get apartment")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Apartment not found")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, apartment)
}

func (ac *ApartmentController) GetAllApartments(res http.ResponseWriter, req *http.Request) {
	filter, err := api_scripts.ParseApartmentFilter(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	apartments, err := ac.Rep.GetAll(filter)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, apartments)
}


func (ac *ApartmentController) GetMyApartments(res http.ResponseWriter, req *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}

	limit := inf
	offset := 0
	
	filter := models.ApartmentFilter{
		SellerID : &userID,
		Limit : &limit,
		Offset : &offset,
	}
	apartments, err := ac.Rep.GetAll(&filter)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, apartments)
}

func (ac *ApartmentController) CreateApartment(res http.ResponseWriter, req *http.Request) {

	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}
	var requestBody struct {
		SellerID     int64    `json:"seller_id"`
		Name         string    `json:"name"`
		Description  string   `json:"description"`
		Capacity     int16    `json:"capacity"`
		PricePerHour int64    `json:"price_per_hour"`
		IsActive     bool      `json:"is_active"`
    }
	
    err := json.NewDecoder(req.Body).Decode(&requestBody)
    if err != nil {
        api_scripts.RespondError(res, http.StatusBadRequest, "Неверный JSON")
        return
    }
    mes := utils.ValidateRequired(requestBody.Name)
    if mes != "" {
        api_scripts.RespondError(res, http.StatusBadRequest, mes)
        return
    }
    mes = utils.ValidateRequired(requestBody.Description)
    if mes != "" {
        api_scripts.RespondError(res, http.StatusBadRequest, mes)
        return
    }
	requestBody.SellerID = userID
	
	apartment := &models.Apartment{
		SellerID : requestBody.SellerID,    
		Name : requestBody.Name,         
		Description : requestBody.Description,  
		Capacity : requestBody.Capacity,     
		PricePerHour : requestBody.PricePerHour, 
		IsActive : requestBody.IsActive,
	}

	err = ac.Rep.Create(apartment)    
	if err != nil {
        api_scripts.RespondError(res, http.StatusBadRequest, "Не удалось создать объявление")
        return
    }

    api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
        "id": apartment.ID,
    })
}