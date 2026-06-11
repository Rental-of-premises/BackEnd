package api

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    api_controllers "rent/internal/api/controllers"
    "rent/internal/api/middleware"
    "rent/internal/config"
    "rent/internal/storage/repository"
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
    r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET", "OPTIONS")
    r.HandleFunc("/auth/sign-up", userController.SignUp).Methods("POST", "OPTIONS")
    r.HandleFunc("/auth/sign-in", userController.SignIn).Methods("POST", "OPTIONS")

    r.HandleFunc("/apartments/{id}", apartmentController.GetApartment).Methods("GET", "OPTIONS")
    r.HandleFunc("/apartments", apartmentController.GetAllApartments).Methods("POST", "OPTIONS")

    r.HandleFunc("/bookings/{id}", bookingController.GetBooking).Methods("GET", "OPTIONS")

    // ========== ЗАЩИЩЕННЫЕ МАРШРУТЫ (с проверкой токена) ==========
    protected := r.PathPrefix("/api").Subrouter()
    protected.Use(middleware.AuthMiddleware)

    protected.HandleFunc("/auth/logout", userController.LogOut).Methods("POST", "OPTIONS")
    protected.HandleFunc("/auth/delete", userController.DeleteAccount).Methods("DELETE", "OPTIONS")

    protected.HandleFunc("/account/my-apartments", apartmentController.GetMyApartments).Methods("GET", "OPTIONS")
    protected.HandleFunc("/account/new-apartment", apartmentController.CreateApartment).Methods("POST", "OPTIONS")
    protected.HandleFunc("/account/apartments/{id}/edit", apartmentController.UpdateApartment).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/account/apartments/{id}/delete", apartmentController.DeleteApartment).Methods("DELETE", "OPTIONS")

    protected.HandleFunc("/account/my-bookings", bookingController.GetMyBookings).Methods("GET", "OPTIONS")
    protected.HandleFunc("/account/new-booking", bookingController.CreateBooking).Methods("POST", "OPTIONS")

    port := config.GetSingletonConfig().ServerPort
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}