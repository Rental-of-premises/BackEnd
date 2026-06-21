package api_scripts

import (
    "fmt"
    "log"
    "net/http"

    "rent/internal/models"
    "rent/internal/storage/repository"
    "rent/internal/api/utils"
)

type AvatarHelper struct {
    AvatarRepo *repository.AvatarRepository
}

func NewAvatarHelper(AvatarRepo *repository.AvatarRepository) *AvatarHelper {
    return &AvatarHelper{
        AvatarRepo: AvatarRepo,
    }
}

func (h *AvatarHelper) ImageRepoHandleImageAvatar(req *http.Request, userID int64) (string, error) {
    if h.AvatarRepo == nil {
        return "", fmt.Errorf("не инициализирован AvatarRepo")
    }

    if req.MultipartForm == nil {
        return "", fmt.Errorf("форма не содержит данных")
    }

    files := req.MultipartForm.File["avatar"]
    if len(files) == 0 {
        return "", fmt.Errorf("файл не найден")
    }

    fileHeader := files[0]

    prefix := fmt.Sprintf("avatar_%d", userID)
    uploadedImage, err := utils.SaveUploadedFileAvatar(fileHeader, prefix)
    if err != nil {
        return "", fmt.Errorf("ошибка сохранения файла: %w", err)
    }

    oldAvatar, err := h.AvatarRepo.GetByUserID(userID)
    if err != nil {
        return "", fmt.Errorf("ошибка получения старой аватарки: %w", err)
    }

    if oldAvatar != nil {
        _ = utils.DeleteFile(oldAvatar.ImageURL)
        err = h.AvatarRepo.Delete(oldAvatar.ID)
        if err != nil {
            return "", fmt.Errorf("ошибка удаления старой аватарки: %w", err)
        }
    }

    avatar := &models.Avatar{
        UserID:   userID,
        ImageURL: uploadedImage.ImageURL,
    }
    err = h.AvatarRepo.Create(avatar)
    if err != nil {
        _ = utils.DeleteFile(uploadedImage.ImageURL)
        return "", fmt.Errorf("ошибка сохранения аватарки в БД: %w", err)
    }

    return uploadedImage.ImageURL, nil
}

func (h *AvatarHelper) GetImageByUser(user *models.User) (*models.Avatar, error) {
    if user == nil {
        return nil, nil
    }

    image, err := h.AvatarRepo.GetByUserID(user.ID)
    if err != nil {
        return nil, err
    }

    return image, nil
}

func (h *AvatarHelper) DeleteAvatar(userID int64) error {
    image, err := h.AvatarRepo.GetByUserID(userID)
    if err != nil {
        return err
    }

    if image == nil {
    log.Printf("ℹ️ Аватарка для пользователя %d не найдена", userID)
    return nil
    }
    
    err = utils.DeleteFile(image.ImageURL)
    if err != nil {
        return err
    }

    return h.AvatarRepo.DeleteByUserID(userID)
}