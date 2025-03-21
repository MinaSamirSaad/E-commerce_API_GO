package users

import (
	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(LOGIN, h.handleLogin).Methods(shared.POST_METHOD)
	router.HandleFunc(REGISTER, h.handleRegister).Methods(shared.POST_METHOD)

}
