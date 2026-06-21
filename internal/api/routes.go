package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	api_controllers "rent/internal/api/controllers"
	"rent/internal/api/middleware"
	api_scripts "rent/internal/api/scripts"
	"rent/internal/config"
	"rent/internal/email"
	"rent/internal/storage/repository"
)

func CreateAndRunRoutes() {
	mainRouter := mux.NewRouter()

	// ===== СТАТИЧЕСКИЕ ФАЙЛЫ =====
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		log.Printf("📁 Создаем папку %s", uploadsDir)
		if err := os.MkdirAll(uploadsDir+"/apartments", 0755); err != nil {
			log.Printf("⚠️ Ошибка создания папки: %v", err)
		}
	}

	fileServer := http.FileServer(http.Dir("./uploads"))
	mainRouter.PathPrefix("/uploads/").Handler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			http.StripPrefix("/uploads/", fileServer).ServeHTTP(w, r)
		}),
	)

	// ===== API РОУТЫ =====
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.CORSMiddleware)
	apiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})

	cfg := config.Load()
	db := repository.GetUserRepository().Db
	if db == nil {
		log.Fatal("❌ Репозиторий user не инициализирован! БД не подключена.")
	}

	emailService := email.NewEmailService(cfg)
	iRepo := repository.GetApartmentImageRepository()
	aRepo := repository.GetAvatarRepository()

	userController := &api_controllers.UserController{
		Rep:          repository.GetUserRepository(),
		EmailService: emailService,
		AH:           api_scripts.NewAvatarHelper(aRepo),
	}

	apartmentController := &api_controllers.ApartmentController{
		Rep: repository.GetApartmentRepository(),
		IH:  api_scripts.NewImageHelper(iRepo),
	}

	bookingController := &api_controllers.BookingController{
		Rep:           repository.GetBookingRepository(),
		ApartmentRepo: repository.GetApartmentRepository(),
	}

	reviewController := &api_controllers.ReviewController{
		Rep:           repository.GetReviewRepository(),
		ApartmentRepo: repository.GetApartmentRepository(),
		BookingRepo:   repository.GetBookingRepository(),
	}

    // ========== ПУБЛИЧНЫЕ МАРШРУТЫ (без токена) ==========
    apiRouter.HandleFunc("/users/{id}", userController.GetUser).Methods("GET", "OPTIONS")
    apiRouter.HandleFunc("/auth/sign-up", userController.SignUp).Methods("POST", "OPTIONS")
    apiRouter.HandleFunc("/auth/sign-in", userController.SignIn).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/auth/confirm-email", userController.ConfirmEmail).Methods("GET", "OPTIONS")
    
    apiRouter.HandleFunc("/apartments/{id}", apartmentController.GetApartment).Methods("GET", "OPTIONS")
    apiRouter.HandleFunc("/apartments", apartmentController.GetAllApartments).Methods("GET", "OPTIONS")
    apiRouter.HandleFunc("/amenities", apartmentController.GetAllAmenities).Methods("GET", "OPTIONS")
    apiRouter.HandleFunc("/amenity/{id}", apartmentController.GetAmenity).Methods("GET", "OPTIONS")
    
    apiRouter.HandleFunc("/apartments/{id}/reviews", reviewController.GetAllReviews).Methods("GET", "OPTIONS")

    apiRouter.HandleFunc("/bookings/{id}", bookingController.GetBooking).Methods("GET", "OPTIONS")
    apiRouter.HandleFunc("/apartments/{id}/calendar", bookingController.GetBookingsByApartment).Methods("GET", "OPTIONS")

    // ========== ЗАЩИЩЕННЫЕ МАРШРУТЫ (с проверкой токена) ==========
    protected := mainRouter.PathPrefix("/api").Subrouter()
    protected.Use(middleware.AuthMiddleware)

    protected.HandleFunc("/auth/logout", userController.LogOut).Methods("POST", "OPTIONS")
    protected.HandleFunc("/auth/delete", userController.DeleteAccount).Methods("DELETE", "OPTIONS")
    // protected.HandleFunc("/api/account/settings/profile/avatar", userController.ChangeAvatar).Methods("PATCH", "OPTIONS")
    // protected.HandleFunc("/api/account/settings/profile/name", userController.ChangeName).Methods("PATCH", "OPTIONS")
    // protected.HandleFunc("/api/account/settings/security/password", userController.ChangePassword).Methods("PATCH", "OPTIONS")

    protected.HandleFunc("/account/my-apartments", apartmentController.GetMyApartments).Methods("GET", "OPTIONS")
    protected.HandleFunc("/account/new-apartment", apartmentController.CreateApartment).Methods("POST", "OPTIONS")
    protected.HandleFunc("/account/apartments/{id}/edit", apartmentController.UpdateApartment).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/apartments/{id}/delete", apartmentController.DeleteApartment).Methods("DELETE", "OPTIONS")
    protected.HandleFunc("/account/apartments/{id}/upload-images", apartmentController.UploadApartmentImages).Methods("POST", "OPTIONS")
    protected.HandleFunc("/account/apartments/{id}/update-images", apartmentController.UpdateApartmentImages).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/apartments/{id}/delete-images", apartmentController.DeleteAllApartmentImages).Methods("DELETE", "OPTIONS")

    protected.HandleFunc("/account/my-bookings", bookingController.GetMyBookings).Methods("GET", "OPTIONS")
    protected.HandleFunc("/account/new-booking", bookingController.CreateBooking).Methods("POST", "OPTIONS")
    protected.HandleFunc("/account/my-bookings/{id}/cancel", bookingController.CancelBooking).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/bookings/{id}/confirm", bookingController.ConfirmBookingBySeller).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/bookings/{id}/reject", bookingController.RejectBookingBySeller).Methods("PATCH", "OPTIONS")
    protected.HandleFunc("/account/bookings", bookingController.GetBookings).Methods("GET", "OPTIONS")

    protected.HandleFunc("/apartments/{id}/new-review", reviewController.CreateReview).Methods("POST", "OPTIONS")
    protected.HandleFunc("/account/delete-review/{id}", reviewController.DeleteReview).Methods("DELETE", "OPTIONS")

    port := config.GetSingletonConfig().ServerPort
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mainRouter))
}