package api_controllers

import (
    "encoding/json"
    "net/http"
    
    "rent/internal/models"
    api_scripts "rent/internal/api/scripts"
    "rent/internal/storage/repository"
    "golang.org/x/crypto/bcrypt"
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


func (uc *UserController) SignUp(res http.ResponseWriter, req *http.Request) {
    var requestBody struct {
        Email    string `json:"email"`
        Password string `json:"password"`
        Name     string `json:"name"`
    }

    err := json.NewDecoder(req.Body).Decode(&requestBody)
    if err != nil {
        api_scripts.RespondError(res, http.StatusBadRequest, "Неверный JSON")
        return
    }

    if requestBody.Email == "" || requestBody.Password == "" {
        api_scripts.RespondError(res, http.StatusBadRequest, "Email и пароль обязательны")
        return
    }

    existingUser, _ := uc.Rep.GetByEmail(requestBody.Email)
    if existingUser != nil {
        api_scripts.RespondError(res, http.StatusConflict, "Пользователь уже существует")
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка хеширования пароля")
        return
    }

    user := &models.User{
		Name:     requestBody.Name,
		Password: string(hashed),
        Email:    requestBody.Email,
    }

    err = uc.Rep.Create(user)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка сохранения пользователя")
        return
    }

    api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
        "id":    user.ID,
        "email": user.Email,
        "name":  user.Name,
    })
}
