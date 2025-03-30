package shared

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func ParseJsonBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJson(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func RespondError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}

var Validate = validator.New()

func GetURLParam(r *http.Request, key string) (string, bool) {
	vars := mux.Vars(r)
	value, ok := vars[key]
	return value, ok
}
