package users

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MinaSamirSaad/ecommerce/services/shared"
)

type Handler struct {
	store UserStore
}

func NewHandler(db *sql.DB) *Handler {
	UserStore := NewStore(db)
	return &Handler{store: UserStore}
}

func (h *Handler) handleLogin(res http.ResponseWriter, req *http.Request) {

}

func (h *Handler) handleRegister(res http.ResponseWriter, req *http.Request) {
	// get payload
	var payload shared.RegisterUserPayload
	if err := shared.ParseJsonBody(req, &payload); err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
	}
	// check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		shared.RespondError(res, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	// create user
	// return response
}
