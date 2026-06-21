package main

import (
	"log"
	"rent/internal/api"
	"rent/internal/storage/repository"
	"rent/internal/storage/migrations"
)


func runController() {
	api.CreateAndRunRoutes()
}


func main() {
	if err := migrations.CreateTablesIfNotExist(); err != nil {
		log.Fatalf("❌ Ошибка создания таблиц: %v", err)
	}

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
	
	runController()
}