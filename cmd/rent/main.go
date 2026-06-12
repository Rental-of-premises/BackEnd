package main

import (
	"database/sql"
	"log"
	"rent/internal/api"
	"rent/internal/models"
	"rent/internal/storage/repository"
	"time"
	//"rent/internal/email"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Протестировано
func printUsers(database *sql.DB) {
	rows, err := database.Query("SELECT * FROM users")
	checkErr(err)
	for rows.Next() {
		user := &models.User{}

		rows.Scan(
			&user.ID,
			&user.Name,
			&user.Password,
			&user.Email,
			&user.CreatedAt,
		)

		log.Println(user)
	}
}

// Протестировано
func createUser(database *sql.DB) {
	userRepo := repository.NewUserRepository(database)

	user := &models.User{
		Name:      "Alis",
		Password:  "1234",
		Email:     "obidjanowa@mail.ru",
		CreatedAt: time.Now(),
	}

	err := userRepo.Create(user)

	checkErr(err)
}

// протестировано
func runController() {
	api.CreateAndRunRoutes()
}

// протестировано
func getUserById() {
	r := repository.GetUserRepository()
	result, _ := r.GetByID(1)
	log.Printf("user by id: %v", result)
}

func main() {
	runController()
	// getUserById()
}
