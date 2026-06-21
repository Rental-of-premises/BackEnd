package api_scripts

import (
    "fmt"
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
        return "", fmt.Errorf("MultipartForm is nil. Вызовите ParseMultipartForm в контроллере")
    }

    if req.MultipartForm.File == nil {
        return "", fmt.Errorf("req.MultipartForm.File is nil")
    }

    files, ok := req.MultipartForm.File["avatar"]
    if !ok || len(files) == 0 {
        keys := make([]string, 0, len(req.MultipartForm.File))
        for k := range req.MultipartForm.File {
            keys = append(keys, k)
        }
        return "", fmt.Errorf("файл 'avatar' не найден. Доступные ключи: %v", keys)
    }

    fileHeader := files[0]

    if fileHeader.Size > 10<<20 {
        return "", fmt.Errorf("размер файла превышает 10 MB")
    }

    prefix := fmt.Sprintf("avatar_%d", userID)
    base64Data, err := utils.SaveUploadedFileAvatar(fileHeader, prefix)
    if err != nil {
        return "", fmt.Errorf("ошибка сохранения файла: %w", err)
    }

    oldAvatar, err := h.AvatarRepo.GetByUserID(userID)
    if err != nil {
        return "", fmt.Errorf("ошибка получения старой аватарки: %w", err)
    }

    if oldAvatar != nil {
        err = h.AvatarRepo.Delete(oldAvatar.ID)
        if err != nil {
            return "", fmt.Errorf("ошибка удаления старой аватарки: %w", err)
        }
    }

    avatar := &models.Avatar{
        UserID:    userID,
        ImageData: base64Data, 
    }
    err = h.AvatarRepo.Create(avatar)
    if err != nil {
        return "", fmt.Errorf("ошибка сохранения аватарки в БД: %w", err)
    }

    return base64Data, nil
}

func (h *AvatarHelper) GetAvatarByUser(user *models.User) (*models.Avatar, error) {
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
        return nil
    }

    err = h.AvatarRepo.DeleteByUserID(userID)
    if err != nil {
        return fmt.Errorf("ошибка удаления аватарки из БД: %w", err)
    }

    return nil
}