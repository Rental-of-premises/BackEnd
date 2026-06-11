package api_controllers

import (
    "encoding/json"
    "net/http"
    
    "rent/internal/models"
    //"github.com/gorilla/mux"
    api_scripts "rent/internal/api/scripts"
    "rent/internal/storage/repository"
    "rent/internal/api/utils"
    "rent/internal/api/middleware"
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
    
    user.Password = ""
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

    mes := utils.ValidateEmail(requestBody.Email)
    if mes != "" {
        api_scripts.RespondError(res, http.StatusBadRequest, mes)
        return
    }

    mes = utils.ValidatePassword(requestBody.Password)
    if mes != "" {
        api_scripts.RespondError(res, http.StatusBadRequest, mes)
        return
    }

    existingUser, _ := uc.Rep.GetByEmail(requestBody.Email)
    if existingUser != nil {
        api_scripts.RespondError(res, http.StatusConflict, "Пользователь уже существует")
        return
    }

    hashed, err := utils.HashPassword(requestBody.Password)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка хеширования пароля")
        return
    }

    user := &models.User{
		Name:     requestBody.Name,
		Password: hashed,
        Email:    requestBody.Email,
    }

    err = uc.Rep.Create(user)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка сохранения пользователя")
        return
    }

    api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
        "id":    user.ID,
    })
}

func (uc *UserController) SignIn(res http.ResponseWriter, req *http.Request) {
    var requestBody struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(req.Body).Decode(&requestBody)
    if err != nil {
        api_scripts.RespondError(res, http.StatusBadRequest, "Неверный JSON")
        return
    }


    user, err := uc.Rep.GetByEmail(requestBody.Email)
    if err != nil || user == nil {
        api_scripts.RespondError(res, http.StatusUnauthorized, "Пользователь не существует")
        return
    }

    err = utils.CheckPassword(requestBody.Password, user.Password)
    if err != nil {
        api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный пароль")
        return
    }

    token, err := utils.GenerateJWT(user.ID, user.Email)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка генерации токена")
        return
    }

     http.SetCookie(res, &http.Cookie{
        Name:     "token",
        Value:    token,
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
        Path:     "/",
        MaxAge:   86400,
    })

    api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
        //"token": token,
        "user": map[string]interface{}{
            "id":    user.ID,
            "email": user.Email,
            "name":  user.Name,
        },
    })
}

func (uc *UserController) LogOut(res http.ResponseWriter, req *http.Request) {
    http.SetCookie(res, &http.Cookie{
        Name:     "token",
        Value:    "",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
        Path:     "/",
        MaxAge:   -1,
    })

    api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
        "message": "Успешный выход из системы",
    })
}

func (uc * UserController) DeleteAccount(res http.ResponseWriter, req *http.Request) {
    userID, ok := middleware.GetUserIDFromContext(req)
    if !ok {
        api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
        return
    }

    user, err := uc.Rep.GetByID(userID)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске пользователя")
        return
    }
    if user == nil {
        api_scripts.RespondError(res, http.StatusNotFound, "Пользователь не найден")
        return
    }

    err = uc.Rep.Delete(userID)
    if err != nil {
        api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при удалении аккаунта")
        return
    }

    api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
        "message": "Аккаунт успешно удален",
    })
}