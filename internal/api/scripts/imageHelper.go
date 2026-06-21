package api_scripts

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"rent/internal/api/utils"
	"rent/internal/models"
	"rent/internal/storage/repository"
)

type ImageHelper struct {
	ImageRepo *repository.ApartmentImageRepository
}

func NewImageHelper(imageRepo *repository.ApartmentImageRepository) *ImageHelper {
    return &ImageHelper{
        ImageRepo:  imageRepo,
    }
}

func (h *ImageHelper) ImageRepoHandleImages(req *http.Request, apartmentID int64) ([]string, []int64, error) {
	var newImages []string
	var deletedImageIDs []int64

	deleteImagesStr := req.FormValue("delete_images")
	if deleteImagesStr != "" {
		ids := strings.Split(deleteImagesStr, ",")
		for _, idStr := range ids {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err == nil {
				deletedImageIDs = append(deletedImageIDs, id)
			}
		}
	}

	if len(deletedImageIDs) > 0 {
		for _, imageID := range deletedImageIDs {
			image, err := h.ImageRepo.GetByID(imageID)
			if err == nil && image.ApartmentID == apartmentID {
				_ = h.ImageRepo.Delete(imageID)
				log.Printf("🗑️ Удалено изображение #%d из БД", imageID)
			}
		}
	}

	if h.ImageRepo == nil {
		return newImages, deletedImageIDs, fmt.Errorf("не инициализирован ImageRepo")
	}

	files := req.MultipartForm.File["images"]
	if len(files) > 0 {
		existingImages, _ := h.ImageRepo.GetByApartmentID(apartmentID)
		nextPosition := len(existingImages)

		base64Images, err := utils.SaveUploadedFiles(files, fmt.Sprintf("apartment_%d", apartmentID))
		if err != nil {
			return newImages, deletedImageIDs, err
		}

		for _, base64Data := range base64Images {
			image := &models.ApartmentImage{
				ApartmentID: apartmentID,
				ImageData:   base64Data, 
				Position:    nextPosition,
			}
			if err := h.ImageRepo.Create(image); err != nil {
				log.Printf("❌ Ошибка сохранения изображения в БД: %v", err)
				continue
			}
			nextPosition++
			newImages = append(newImages, base64Data) // 
		}
		
		log.Printf("📸 Загружено %d новых изображений", len(newImages))
	}

	return newImages, deletedImageIDs, nil
}

func (h *ImageHelper) GetImagesByApartment(apartment *models.Apartment) ([]*models.ApartmentImage, error) {
	if apartment == nil {
		return []*models.ApartmentImage{}, nil
	}

	images, err := h.ImageRepo.GetByApartmentID(apartment.ID)
	if err != nil {
		return []*models.ApartmentImage{}, err
	}

	return images, nil
}

func (h *ImageHelper) DeleteAllImages(apartmentID int64) error {
	log.Printf("🗑️ DeleteAllImages: удаление всех изображений для помещения %d", apartmentID)

	// Получаем изображения
	images, err := h.ImageRepo.GetByApartmentID(apartmentID)
	if err != nil {
		log.Printf("⚠️ Ошибка получения изображений: %v", err)
		return err
	}

	log.Printf("📸 Найдено %d изображений для удаления", len(images))

	if err := h.ImageRepo.DeleteByApartmentID(apartmentID); err != nil {
		log.Printf("⚠️ Ошибка удаления записей из БД: %v", err)
		return err
	}

	log.Printf("✅ Все изображения для помещения %d удалены из БД", apartmentID)
	return nil
}

func (h *ImageHelper) DeleteImage(imageID int64) error {
	log.Printf("🗑️ DeleteImage: удаление изображения #%d", imageID)

	image, err := h.ImageRepo.GetByID(imageID)
	if err != nil {
		return fmt.Errorf("ошибка получения изображения: %w", err)
	}

	if image == nil {
		return fmt.Errorf("изображение не найдено")
	}

	if err := h.ImageRepo.Delete(imageID); err != nil {
		return fmt.Errorf("ошибка удаления из БД: %w", err)
	}

	log.Printf("✅ Изображение #%d удалено из БД", imageID)
	return nil
}
