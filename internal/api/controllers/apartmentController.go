package api_controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"rent/internal/api/middleware"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/models"
	"rent/internal/storage/repository"
)

type ApartmentController struct {
	Rep *repository.ApartmentRepository
	IH  *api_scripts.ImageHelper
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

	images, err := ac.IH.GetImagesByApartment(apartment)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка получения изображений: "+err.Error())
		return
	}

	response := struct {
		Apartment *models.Apartment        `json:"apartment"`
		Images    []*models.ApartmentImage `json:"images"`
	}{
		Apartment: apartment,
		Images:    images,
	}

	api_scripts.RespondJSON(res, http.StatusOK, response)
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

	var allImages [][]*models.ApartmentImage
	for _, apartment := range apartments {
		images, err := ac.IH.GetImagesByApartment(apartment)
		if err != nil {
			api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка получения изображений: "+err.Error())
			return
		}
		allImages = append(allImages, images)
	}

	response := struct {
		Apartments []*models.Apartment        `json:"apartments"`
		Images     [][]*models.ApartmentImage `json:"images"`
	}{
		Apartments: apartments,
		Images:     allImages,
	}

	api_scripts.RespondJSON(res, http.StatusOK, response)
}

func (ac *ApartmentController) GetMyApartments(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}

	filter, err := api_scripts.ParseApartmentFilter(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}
	filter.SellerID = &userID

	apartments, err := ac.Rep.GetAll(filter)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	var allImages [][]*models.ApartmentImage
	for _, apartment := range apartments {
		images, err := ac.IH.GetImagesByApartment(apartment)
		if err != nil {
			api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка получения изображений: "+err.Error())
			return
		}
		allImages = append(allImages, images)
	}

	response := struct {
		Apartments []*models.Apartment        `json:"apartments"`
		Images     [][]*models.ApartmentImage `json:"images"`
	}{
		Apartments: apartments,
		Images:     allImages,
	}

	api_scripts.RespondJSON(res, http.StatusOK, response)
}

func (ac *ApartmentController) CreateApartment(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный тип id в контексте")
		return
	}

	var requestBody struct {
		SellerID     int64    `json:"seller_id"`
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		Capacity     int16    `json:"capacity"`
		PricePerHour int64    `json:"price_per_hour"`
		IsActive     bool     `json:"is_active"`
		ImageURL     string   `json:"image_url"`
		Metro        string   `json:"metro"`
		Address      string   `json:"address"`
		Amenities    []string `json:"amenities"`
	}

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Неверный JSON")
		return
	}

	// Валидация
	if requestBody.Name == "" {
		api_scripts.RespondError(res, http.StatusBadRequest, "Название обязательно")
		return
	}

	if requestBody.PricePerHour <= 0 {
		api_scripts.RespondError(res, http.StatusBadRequest, "Цена должна быть больше 0")
		return
	}

	if requestBody.Capacity <= 0 {
		api_scripts.RespondError(res, http.StatusBadRequest, "Вместимость должна быть больше 0")
		return
	}

	requestBody.SellerID = userID

	apartment := &models.Apartment{
		SellerID:     requestBody.SellerID,
		Name:         requestBody.Name,
		Description:  requestBody.Description,
		Capacity:     requestBody.Capacity,
		PricePerHour: requestBody.PricePerHour,
		IsActive:     requestBody.IsActive,
		Metro:        &requestBody.Metro,
		Address:      &requestBody.Address,
	}

	// Сохраняем помещение в БД
	err = ac.Rep.Create(apartment)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Не удалось создать объявление: "+err.Error())
		return
	}

	// Обработка изображения
	if requestBody.ImageURL != "" {
		imageRepo := repository.GetApartmentImageRepository()

		if strings.HasPrefix(requestBody.ImageURL, "data:image") {
			if len(requestBody.ImageURL) > 7*1024*1024 {
				api_scripts.RespondError(res, http.StatusBadRequest,
					"Изображение слишком большое для base64. Используйте /upload-images")
				return
			}

			image := &models.ApartmentImage{
				ApartmentID: apartment.ID,
				ImageURL:    requestBody.ImageURL,
				Position:    0,
			}

			if err := imageRepo.Create(image); err != nil {
				log.Printf("⚠️ Ошибка сохранения изображения: %v", err)
			}
		} else if strings.HasPrefix(requestBody.ImageURL, "/uploads/") {
			image := &models.ApartmentImage{
				ApartmentID: apartment.ID,
				ImageURL:    requestBody.ImageURL,
				Position:    0,
			}

			if err := imageRepo.Create(image); err != nil {
				log.Printf("⚠️ Ошибка сохранения изображения: %v", err)
			}
		}
	}

	api_scripts.RespondJSON(res, http.StatusCreated, map[string]interface{}{
		"id":      apartment.ID,
		"message": "Объявление успешно создано",
	})
}

func (ac *ApartmentController) UploadApartmentImages(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	apartmentID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	apartment, err := ac.Rep.GetByID(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске помещения")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Помещение не найдено")
		return
	}
	if apartment.SellerID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "У вас нет прав на изменение этого помещения")
		return
	}

	err = req.ParseMultipartForm(32 << 20)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Ошибка парсинга формы: "+err.Error())
		return
	}

	if req.MultipartForm == nil || len(req.MultipartForm.File["images"]) == 0 {
		api_scripts.RespondError(res, http.StatusBadRequest, "Файлы не найдены")
		return
	}

	// ✅ УДАЛЯЕМ ВСЕ СТАРЫЕ ИЗОБРАЖЕНИЯ ПЕРЕД ЗАГРУЗКОЙ НОВЫХ
	log.Printf("🗑️ Удаляем старые изображения для помещения %d", apartmentID)
	err = ac.IH.DeleteAllImages(apartmentID)
	if err != nil {
		log.Printf("⚠️ Ошибка удаления старых изображений: %v", err)
		// Продолжаем выполнение, даже если не удалось удалить
	}

	// Загружаем новые изображения
	imageURLs, _, err := ac.IH.ImageRepoHandleImages(req, apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Не удалось сохранить фотографии: "+err.Error())
		return
	}

	log.Printf("✅ Загружено %d новых изображений для помещения %d", len(imageURLs), apartmentID)

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message":      "Изображения успешно загружены",
		"images":       imageURLs,
		"images_count": len(imageURLs),
	})
}

func (ac *ApartmentController) UpdateApartment(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	apartmentID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	apartment, err := ac.Rep.GetByID(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске помещения")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Помещение не найдено")
		return
	}

	if apartment.SellerID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "У вас нет прав на редактирование этого помещения")
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Неверный JSON: "+err.Error())
		return
	}

	log.Printf("📝 Обновление помещения %d: %+v", apartmentID, updates)

	if len(updates) == 0 {
		api_scripts.RespondError(res, http.StatusBadRequest, "Нет данных для обновления")
		return
	}

	err = ac.Rep.UpdatePartial(apartmentID, updates)
	if err != nil {
		log.Printf("❌ Ошибка обновления: %v", err)
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при обновлении помещения: "+err.Error())
		return
	}

	updatedApartment, err := ac.Rep.GetByID(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при получении обновленного помещения")
		return
	}

	images, err := ac.IH.GetImagesByApartment(updatedApartment)
	if err != nil {
		images = []*models.ApartmentImage{}
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"apartment": updatedApartment,
		"images":    images,
		"message":   "Объявление успешно обновлено",
	})
}

func (ac *ApartmentController) UpdateApartmentImages(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	apartmentID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	apartment, err := ac.Rep.GetByID(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске помещения")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Помещение не найдено")
		return
	}

	if apartment.SellerID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "У вас нет прав на редактирование этого помещения")
		return
	}

	err = req.ParseMultipartForm(32 << 20)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Ошибка парсинга формы: "+err.Error())
		return
	}

	imageURLs, deletedIDs, err := ac.IH.ImageRepoHandleImages(req, apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, "Не удалось изменить фотографии: "+err.Error())
		return
	}

	updatedImages, err := ac.IH.GetImagesByApartment(apartment)
	if err != nil {
		updatedImages = []*models.ApartmentImage{}
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message":        "Изображения успешно обновлены",
		"images":         updatedImages,
		"new_images":     imageURLs,
		"deleted_images": deletedIDs,
		"images_count":   len(updatedImages),
	})
}

func (ac *ApartmentController) DeleteApartment(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	apartmentID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	apartment, err := ac.Rep.GetByID(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске помещения")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Помещение не найдено")
		return
	}

	if apartment.SellerID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "У вас нет прав на удаление этого помещения")
		return
	}

	err = ac.Rep.Delete(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при удалении помещения")
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Объявление успешно удалено",
	})
}

func (ac *ApartmentController) DeleteAllApartmentImages(res http.ResponseWriter, req *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(req)
	if !ok {
		api_scripts.RespondError(res, http.StatusUnauthorized, "Не авторизован")
		return
	}

	apartmentID, err := api_scripts.ParseID(req)
	if err != nil {
		api_scripts.RespondError(res, http.StatusBadRequest, err.Error())
		return
	}

	apartment, err := ac.Rep.GetByID(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при поиске помещения")
		return
	}
	if apartment == nil {
		api_scripts.RespondError(res, http.StatusNotFound, "Помещение не найдено")
		return
	}

	if apartment.SellerID != userID {
		api_scripts.RespondError(res, http.StatusForbidden, "У вас нет прав на редактирование этого помещения")
		return
	}

	err = ac.IH.DeleteAllImages(apartmentID)
	if err != nil {
		api_scripts.RespondError(res, http.StatusInternalServerError, "Ошибка при удалении изображений: "+err.Error())
		return
	}

	api_scripts.RespondJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Все изображения успешно удалены",
	})
}
