package api_scripts

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    "rent/internal/models"
    "rent/internal/storage/repository"
    "rent/internal/api/utils"
)

type ImageHelper struct {
    ImageRepo *repository.ApartmentImageRepository
}

func NewImageHelper(imageRepo *repository.ApartmentImageRepository, avatarRepo *repository.AvatarRepository) *ImageHelper {
    return &ImageHelper{
        ImageRepo:  imageRepo,
        AvatarRepo: avatarRepo,
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
        return newImages, deletedImageIDs, fmt.Errorf("Не инициализирован ImageRepo")
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
    images, err := h.ImageRepo.GetByApartmentID(apartmentID)
    if err != nil {
        return err
    }

    for _, img := range images {
        err = utils.DeleteFile(img.ImageURL)
        if err != nil {
            return err
        }
    }

    return h.ImageRepo.DeleteByApartmentID(apartmentID)
}


func (h *ImageHelper) ImageRepoHandleImageAvatar(req *http.Request, userID int64) ([]string, []int64, error) {
    var newImage string
    var deletedImage int64

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
        return newImages, deletedImageIDs, fmt.Errorf("Не инициализирован ImageRepo")
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

func (h *ImageHelper) GetImageByUser(user *models.User) ([]*models.Avatar, error) {
    if user == nil {
        return []*models.Avatar{}, nil
    }

    image, err := h.ImageRepo.GetByUserID(user.ID)
    if err != nil {
        return []*models.Avatar{}, err
    }

    return images, nil
}

func (h *ImageHelper) DeleteAvatar(userID int64) error {
    image, err := h.ImageRepo.GetByUserID(userID)
    if err != nil {
        return err
    }

    err = utils.DeleteFile(image.ImageURL)
    if err != nil {
        return err
    }

    return h.ImageRepo.DeleteByUserID(userID)
}