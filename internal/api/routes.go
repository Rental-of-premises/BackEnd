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

    // ========== ПУБЛИЧНЫЕ МАРШРУТЫ (без токена) ==========
    r.HandleFunc("/auth/sign-up", userController.SignUp).Methods("POST")
    r.HandleFunc("/auth/sign-in", userController.SignIn).Methods("POST")
	r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")
    
    // ========== ЗАЩИЩЕННЫЙ МАРШРУТ (с проверкой токена) ==========
    // Оборачиваем GetUser в AuthMiddleware
   

    port := config.GetSingletonConfig().ServerPort
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}