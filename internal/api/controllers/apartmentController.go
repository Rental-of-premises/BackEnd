package api_controllers

import (
	"net/http"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/storage/repository"
)

type ApartmentController struct {
	Rep *repository.ApartmentRepository
}

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
	if apartments == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Apartment not found")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, apartments)
}
