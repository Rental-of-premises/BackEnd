package api

import (
	"fmt"
	"log"
	"net/http"
	api_controllers "rent/internal/api/controllers"
	"rent/internal/config"
	"rent/internal/storage/repository"

	"github.com/gorilla/mux"
	"rent/internal/api/middleware"
)

func CreateAndRunRoutes() {
    r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            next.ServeHTTP(w, req)
        })
    })
  
	userController := &api_controllers.UserController{repository.GetUserRepository()}
    apartmentController := &api_controllers.ApartmentController{repository.GetApartmentRepository()}
    bookingController := &api_controllers.BookingController{repository.GetBookingRepository()}
    // reviewController := &api_controllers.UserController{repository.GetUserRepository()}
	
    // ========== ПУБЛИЧНЫЕ МАРШРУТЫ (без токена) ==========
    r.HandleFunc("/auth/sign-up", userController.SignUp).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth/sign-in", userController.SignIn).Methods("POST", "OPTIONS")
	r.HandleFunc("/apartments/{id}", apartmentController.GetApartment).Methods("GET", "OPTIONS")
	r.HandleFunc("/apartments", apartmentController.GetAllApartments).Methods("POST", "OPTIONS")
	
	// reviewController := &api_controllers.UserController{repository.GetUserRepository()}
	// r.HandleFunc("/review/{id}", userController.GetUser).Methods("GET")
	
	r.HandleFunc("/bookings/{id}", bookingController.GetBooking).Methods("GET", "OPTIONS")
    r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET", "OPTIONS")
    r.HandleFunc("/apartments/{id}", apartmentController.GetApartment).Methods("GET", "OPTIONS")
    r.HandleFunc("/apartments", apartmentController.GetAllApartments).Methods("POST", "OTIONS")

	// r.HandleFunc("/review/{id}", userController.GetUser).Methods("GET")
    // ========== ЗАЩИЩЕННЫЙ МАРШРУТ (с проверкой токена) ==========
    // Оборачиваем GetUser в AuthMiddleware
	protected := r.PathPrefix("/api").Subrouter()
    protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/logout", userController.LogOut).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/delete", userController.DeleteAccount).Methods("DELETE", "OPTIONS")
    //r.HandleFunc("/account/apartments", apartmentController.GetAllApartments).Methods("POST")
    //r.HandleFunc("/booking", bookingController.GetAllBookings).Methods("POST")

	port := config.GetSingletonConfig().ServerPort
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
