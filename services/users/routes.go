package users

import (
	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(shared.LOGIN, h.handleLogin).Methods(shared.POST_METHOD)
	router.HandleFunc(shared.REGISTER, h.handleRegister).Methods(shared.POST_METHOD)

}
