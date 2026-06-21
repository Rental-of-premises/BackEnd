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
				_ = utils.DeleteFile(image.ImageURL)
				_ = h.ImageRepo.Delete(imageID)
			}
		}
	}

	if h.ImageRepo == nil {
		return newImages, deletedImageIDs, fmt.Errorf("не инициализирован ImageRepo")
	}

	if req.MultipartForm == nil {
		return newImages, deletedImageIDs, nil
	}

	files := req.MultipartForm.File["images"]
	if len(files) > 0 {
		existingImages, _ := h.ImageRepo.GetByApartmentID(apartmentID)
		nextPosition := len(existingImages)

		prefix := fmt.Sprintf("apartment_%d", apartmentID)
		uploadedImages, err := utils.SaveUploadedFiles(files, prefix)
		if err != nil {
			return newImages, deletedImageIDs, err
		}

		for _, img := range uploadedImages {
			image := &models.ApartmentImage{
				ApartmentID: apartmentID,
				ImageURL:    img.ImageURL,
				Position:    nextPosition,
			}
			if err := h.ImageRepo.Create(image); err != nil {
				log.Printf("Ошибка сохранения изображения в БД: %v", err)
				_ = utils.DeleteFile(img.ImageURL)
				continue
			}
			nextPosition++
			newImages = append(newImages, img.ImageURL)
		}
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

	images, err := h.ImageRepo.GetByApartmentID(apartmentID)
	if err != nil {
		log.Printf("⚠️ Ошибка получения изображений: %v", err)
		return err
	}

	log.Printf("📸 Найдено %d изображений для удаления", len(images))

	for _, img := range images {
		if err := utils.DeleteFile(img.ImageURL); err != nil {
			log.Printf("⚠️ Ошибка удаления файла %s: %v", img.ImageURL, err)
		}
	}

	if err := h.ImageRepo.DeleteByApartmentID(apartmentID); err != nil {
		log.Printf("⚠️ Ошибка удаления записей из БД: %v", err)
		return err
	}

	log.Printf("✅ Все изображения для помещения %d удалены", apartmentID)
	return nil
}
