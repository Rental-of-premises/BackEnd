package main

import (
	"database/sql"
	"log"
	"rent/internal/api"
	"rent/internal/models"
	"rent/internal/storage/repository"
	"time"
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
	// ===== АВТОМАТИЧЕСКОЕ ЗАВЕРШЕНИЕ ПРОШЕДШИХ БРОНИРОВАНИЙ ПРИ СТАРТЕ =====
	log.Println("🔄 Проверка прошедших бронирований...")
	
	bookingRepo := repository.GetBookingRepository()
	if bookingRepo != nil {
		if err := bookingRepo.CompletePastBookings(); err != nil {
			log.Printf("⚠️ Ошибка при завершении прошедших бронирований: %v", err)
		} else {
			log.Println("✅ Проверка прошедших бронирований завершена")
		}
	} else {
		log.Println("⚠️ BookingRepository не инициализирован")
	}
	
	// ===== ЗАПУСК СЕРВЕРА =====
	runController()
	// getUserById()
}
