package api_controllers

import (
	"net/http"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/storage/repository"
)

type UserController struct {
	Rep *repository.UserRepository
}

func (uc *UserController) GetUser(res http.ResponseWriter, req *http.Request) {
	id, err := api_scripts.ParseID(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	user, err := uc.Rep.GetByID(id)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if user == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "User not found")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, user)
}
