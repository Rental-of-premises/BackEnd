package api

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    api_controllers "rent/internal/api/controllers"
    "rent/internal/api/middleware"
    "rent/internal/config"
    "rent/internal/email"
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

    // Загружаем конфиг и создаём email сервис
    cfg := config.Load()
    db := repository.GetUserRepository().Db
    if db == nil {
        log.Fatal("❌ Репозиторий user не инициализирован! БД не подключена.")
    }
    
    var dbName, dbUser string
    err := db.QueryRow("SELECT current_database(), current_user").Scan(&dbName, &dbUser)
    if err != nil {
        log.Printf("❌ Ошибка получения информации о БД: %v", err)
    } else {
        log.Printf("✅ Подключено к БД: %s, пользователь: %s", dbName, dbUser)
    }
    
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    if err != nil {
        log.Printf("❌ Ошибка подсчета пользователей: %v", err)
    } else {
        log.Printf("📊 В таблице users: %d записей", count)
    }
    emailService := email.NewEmailService(cfg)

    // Контроллеры
    userController := &api_controllers.UserController{
        Rep:          repository.GetUserRepository(),
        EmailService: emailService,
    }
    
    apartmentController := &api_controllers.ApartmentController{
        Rep: repository.GetApartmentRepository(),
    }
    
    bookingController := &api_controllers.BookingController{
        Rep:           repository.GetBookingRepository(),
        ApartmentRepo: repository.GetApartmentRepository(),
    }
    
    // reviewController := &api_controllers.UserController{repository.GetUserRepository()}

    // ========== ПУБЛИЧНЫЕ МАРШРУТЫ (без токена) ==========
    r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET", "OPTIONS")
    r.HandleFunc("/auth/sign-up", userController.SignUp).Methods("POST", "OPTIONS")
    r.HandleFunc("/auth/sign-in", userController.SignIn).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth/confirm-email", userController.ConfirmEmail).Methods("GET", "OPTIONS")

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
    protected.HandleFunc("/account/my-bookings/{id}/cancel", bookingController.CancelBooking).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/bookings/{id}/confirm", bookingController.ConfirmBookingBySeller).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/bookings/{id}/reject", bookingController.RejectBookingBySeller).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/bookings", bookingController.GetBookings).Methods("GET", "OPTIONS")

    port := config.GetSingletonConfig().ServerPort
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}