package api_controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"rent/internal/api/middleware"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/api/utils"
	"rent/internal/email"
	"rent/internal/models"
	"rent/internal/storage/repository"
)

type UserController struct {
	Rep          *repository.UserRepository
	EmailService *email.EmailService
}

func (uc *UserController) GetUser(res http.ResponseWriter, req *http.Request) {
	id, err := api_scripts.ParseID(req)

	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	user, err := uc.Rep.GetByID(id)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Failed to get user: " + err.Error())
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

	if len(requestBody.Password) < 60 {
		api_scripts.RespondError(res, http.StatusBadRequest, "Неверный формат пароля")
		return
	}

	existingUser, _ := uc.Rep.GetByEmail(requestBody.Email)
	if existingUser != nil {
		api_scripts.RespondError(res, http.StatusConflict, "Пользователь уже существует")
		return
	}

	user := &models.User{
		Name:       requestBody.Name,
		Email:      requestBody.Email,
		Password:   requestBody.Password,
		IsActive:   false,
		EmailToken: nil,
	}

	err = uc.Rep.Create(user)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка сохранения пользователя")
		return
	}

	token, err := utils.GenerateRandomToken()
	if err != nil {
		log.Printf("❌ GenerateRandomToken ошибка: %v", err)
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при регистрации")
		return
	}

	err = uc.Rep.UpdateEmailToken(user.ID, token)
	if err != nil {
		log.Printf("❌ UpdateEmailToken ошибка: %v", err)
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при регистрации")
		return
	}

	go func() {
		confirmURL := fmt.Sprintf("https://team3.verstack.ru/auth/confirm-email?token=%s", token)

		data := struct {
			Name       string
			ConfirmURL string
		}{
			Name:       user.Name,
			ConfirmURL: confirmURL,
		}

		body, err := uc.EmailService.RenderTemplate("welcome.html", data)
		if err != nil {
			return
		}

		uc.EmailService.SendEmail(user.Email, "Подтверждение регистрации", body)
	}()

	api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
		"id": user.ID,
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

	if !user.IsActive {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Подтвердите email перед входом. Проверьте почту.")
		return
	}
	if user.Password != requestBody.Password {
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

func (uc *UserController) DeleteAccount(res http.ResponseWriter, req *http.Request) {
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

func (uc *UserController) ConfirmEmail(res http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	if token == "" {
		api_scripts.RespondError(res, http.StatusBadRequest, "Токен не указан")
		return
	}

	userID, err := uc.Rep.ActivateUser(token)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка активации")
		return
	}
	if userID == 0 {
		api_scripts.RespondError(res, http.StatusNotFound, "Неверный или просроченный токен")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Email успешно подтверждён! Теперь вы можете войти.",
	})
}
