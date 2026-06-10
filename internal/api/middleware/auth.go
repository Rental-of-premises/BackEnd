package middleware

import (
    "context"
    "net/http"
    
    "rent/internal/api/utils"
    api_scripts "rent/internal/api/scripts"
)

type contextKey string

const (
    UserIDKey    contextKey = "userID"
    UserEmailKey contextKey = "userEmail"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        tokenString := utils.ExtractToken(req)
        if tokenString == "" {
            api_scripts.RespondError(res, http.StatusUnauthorized, "Отсутствует токен авторизации")
            return
        }
        
        claims, err := utils.ParseJWT(tokenString)
        if err != nil {
            api_scripts.RespondError(res, http.StatusUnauthorized, "Неверный или просроченный токен")
            return
        }
        
        ctx := context.WithValue(req.Context(), UserIDKey, claims.ID)
        ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
        
        next.ServeHTTP(res, req.WithContext(ctx))
    })
}

func GetUserIDFromContext(req *http.Request) (int64, bool) {
    userID, ok := req.Context().Value(UserIDKey).(int64)
    return userID, ok
}

func GetUserEmailFromContext(req *http.Request) (string, bool) {
    email, ok := req.Context().Value(UserEmailKey).(string)
    return email, ok
}