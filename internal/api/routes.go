package api

import (
	"fmt"
	"log"
	"net/http"
	api_controllers "rent/internal/api/controllers"
	"rent/internal/config"
	"rent/internal/storage/repository"

	"github.com/gorilla/mux"
	//"rent/internal/api/middleware"
)

func CreateAndRunRoutes() {
    r := mux.NewRouter()
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            next.ServeHTTP(w, req)
        })
    })
  
	userController := &api_controllers.UserController{repository.GetUserRepository()}
  apartmentController := &api_controllers.ApartmentController{repository.GetApartmentRepository()}
  bookingController := &api_controllers.BookingController{repository.GetBookingRepository()}
	
  // ========== ПУБЛИЧНЫЕ МАРШРУТЫ (без токена) ==========
  r.HandleFunc("/auth/sign-up", userController.SignUp).Methods("POST")
  r.HandleFunc("/auth/sign-in", userController.SignIn).Methods("POST")
  r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")
  r.HandleFunc("/apartments/{id}", apartmentController.GetApartment).Methods("GET")
  r.HandleFunc("/apartments", apartmentController.GetAllApartments).Methods("POST")
	
	// reviewController := &api_controllers.UserController{repository.GetUserRepository()}
	// r.HandleFunc("/review/{id}", userController.GetUser).Methods("GET")
    
  r.HandleFunc("/booking/{id}", bookingController.GetBooking).Methods("GET")
	//r.HandleFunc("/booking", bookingController.GetAllBookings).Methods("POST")
  // ========== ЗАЩИЩЕННЫЙ МАРШРУТ (с проверкой токена) ==========
  // Оборачиваем GetUser в AuthMiddleware

	port := config.GetSingletonConfig().ServerPort
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
