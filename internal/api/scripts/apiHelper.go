package api_scripts

import (
	"encoding/json"
	"fmt"
	"net/http"
	api_models "rent/internal/api/models"
	"strconv"

	"github.com/gorilla/mux"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, &api_models.ErrorResponse{Error: message})
}

func ParseID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID")
	}
	return id, nil
}
