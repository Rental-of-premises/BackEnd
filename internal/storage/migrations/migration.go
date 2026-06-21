package migrations

import (
    "embed"
    "log"
    "rent/internal/storage/db"
)

//go:embed *.sql
var migrationFiles embed.FS 

func CreateTablesIfNotExist() error {
    db, err := db.GetSingletonDB()
    if err != nil {
        return err
    }
    files, err := migrationFiles.ReadDir(".")
    if err != nil {
        return err
    }

    for _, file := range files {
        content, err := migrationFiles.ReadFile(file.Name())
        if err != nil {
            return err
        }

        if _, err := db.Exec(string(content)); err != nil {
            log.Printf("❌ Ошибка в %s: %v", file.Name(), err)
            return err
        }
        log.Printf("✅ Выполнен: %s", file.Name())
    }

    return nil
}