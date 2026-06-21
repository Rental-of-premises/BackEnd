package middleware

import (
    "net/http"
    "slices"
)

var allowedOrigins = []string{
    "https://yourdomain.com",
    "https://www.yourdomain.com",
    "http://localhost:3000", 
    "http://localhost:5173", 
    "http://localhost:8080",
    "http://localhost:5173",  
    "http://localhost:5432", 
}

func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        isAllowed := slices.Contains(allowedOrigins, origin)
        
        if isAllowed {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Max-Age", "86400")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}