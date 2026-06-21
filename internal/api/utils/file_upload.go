package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rent/internal/models"
)

const (
<<<<<<< HEAD
    MaxFileSize = 20 << 20
    UploadDir   = "./uploads/apartments"
    UploadURL   = "/uploads/apartments"
    UploadAvatarsDir   = "./uploads/avatars"
    UploadAvatarURL   = "/uploads/avatars"
=======
	MaxFileSize = 20 << 20 // 20 MB
	UploadDir   = "./uploads/apartments"
>>>>>>> main
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

func ValidateImage(fileHeader *multipart.FileHeader) error {
<<<<<<< HEAD
    if fileHeader.Size > MaxFileSize {
        return fmt.Errorf("размер файла не должен превышать 20MB")
    }
=======
	if fileHeader.Size > MaxFileSize {
		return fmt.Errorf("размер файла не должен превышать 20MB")
	}
>>>>>>> main

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !AllowedExtensions[ext] {
		return fmt.Errorf("допустимые форматы: jpg, jpeg, png, webp, gif")
	}

	return nil
}

func SaveUploadedFiles(files []*multipart.FileHeader, prefix string) ([]*models.ApartmentImage, error) {
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		return nil, fmt.Errorf("ошибка создания папки для загрузок: %w", err)
	}

	var uploadedImages []*models.ApartmentImage
	timestamp := time.Now().Unix()

	for idx, fileHeader := range files {
		if err := ValidateImage(fileHeader); err != nil {
			return nil, err
		}

		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)

		fileName := fmt.Sprintf("%s_%d_%d%s", prefix, timestamp, idx, ext)
		filePath := filepath.Join(UploadDir, fileName)

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			return nil, err
		}

<<<<<<< HEAD
        imageURL := fmt.Sprintf("%s/%s", UploadURL, fileName)
=======
		imageURL := fmt.Sprintf("/uploads/apartments/%s", fileName)
>>>>>>> main

		uploadedImages = append(uploadedImages, &models.ApartmentImage{
			ImageURL: imageURL,
			Position: idx,
		})
	}

	return uploadedImages, nil
}

func SaveUploadedFileAvatar(fileHeader *multipart.FileHeader, prefix string) (*models.Avatar, error) {

    if err := os.MkdirAll(UploadAvatarsDir, 0755); err != nil {
        return nil, fmt.Errorf("ошибка создания папки для загрузок: %w", err)
    }

    if err := ValidateImage(fileHeader); err != nil {
        return nil, err
    }

    file, err := fileHeader.Open()
    timestamp := time.Now().Unix()
    if err != nil {
        return nil, err
    }
    defer file.Close()

    ext := filepath.Ext(fileHeader.Filename) 
    
    fileName := fmt.Sprintf("%s_%d_%s", prefix, timestamp, ext)
    filePath := filepath.Join(UploadAvatarsDir, fileName)

    dst, err := os.Create(filePath)
    if err != nil {
        return nil, err
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        return nil, err
    }

    imageURL := fmt.Sprintf("/%s/%s", UploadAvatarURL, fileName)

    return &models.Avatar{ ImageURL: imageURL,}, nil
}

func DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("ошибка удаления файла: %w", err)
	}
	return nil
}

func DeleteMultipleFiles(filePaths []string) error {
	var errors []string
	for _, path := range filePaths {
		if err := DeleteFile(path); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("ошибки при удалении: %s", strings.Join(errors, "; "))
	}
	return nil
}
