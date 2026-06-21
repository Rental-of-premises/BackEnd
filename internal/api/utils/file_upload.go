package utils

import (
    "encoding/base64"
    "fmt"
    "io"
    "mime/multipart"
    "path/filepath"
    "strings"
)

const (
    MaxFileSize = 20 << 20
)

var AllowedExtensions = map[string]bool{
    ".jpg":  true,
    ".jpeg": true,
    ".png":  true,
    ".webp": true,
    ".gif":  true,
}

func ValidateImage(fileHeader *multipart.FileHeader) error {
    if fileHeader.Size > MaxFileSize {
        return fmt.Errorf("размер файла не должен превышать 20MB")
    }

    ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
    if !AllowedExtensions[ext] {
        return fmt.Errorf("допустимые форматы: jpg, jpeg, png, webp, gif")
    }

    return nil
}

func ConvertToBase64(fileHeader *multipart.FileHeader) (string, error) {
    file, err := fileHeader.Open()
    if err != nil {
        return "", fmt.Errorf("не удалось открыть файл: %w", err)
    }
    defer file.Close()

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        return "", fmt.Errorf("не удалось прочитать файл: %w", err)
    }

    mimeType := detectMimeType(fileBytes)

    base64Str := base64.StdEncoding.EncodeToString(fileBytes)
    
    dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
    
    return dataURL, nil
}

func detectMimeType(data []byte) string {
    mimeTypes := map[string]string{
        "\xFF\xD8\xFF": "image/jpeg",
        "\x89PNG":      "image/png",
        "GIF":          "image/gif",
        "RIFF":         "image/webp",
    }
    
    for signature, mimeType := range mimeTypes {
        if len(data) >= len(signature) && string(data[:len(signature)]) == signature {
            return mimeType
        }
    }
    
    return "image/jpeg" 
}

func SaveUploadedFiles(files []*multipart.FileHeader, prefix string) ([]string, error) {
    var base64Images []string
    
    for _, fileHeader := range files {
        if err := ValidateImage(fileHeader); err != nil {
            return nil, err
        }

        base64Str, err := ConvertToBase64(fileHeader)
        if err != nil {
            return nil, err
        }

        base64Images = append(base64Images, base64Str)
    }

    return base64Images, nil
}

func SaveUploadedFileAvatar(fileHeader *multipart.FileHeader, prefix string) (string, error) {
    // Валидация
    if err := ValidateImage(fileHeader); err != nil {
        return "", err
    }

    base64Str, err := ConvertToBase64(fileHeader)
    if err != nil {
        return "", err
    }

    return base64Str, nil
}