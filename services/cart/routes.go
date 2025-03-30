package cart

import (
	"github.com/MinaSamirSaad/ecommerce/services/auth"
	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(shared.CHECKOUT, auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(shared.POST_METHOD)
}
