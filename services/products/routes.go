package products

import (
	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(shared.PRODUCTS, h.GetProducts).Methods(shared.GET_METHOD)
	router.HandleFunc(shared.PRODUCTS, h.CreateProduct).Methods(shared.POST_METHOD)
	router.HandleFunc(shared.PRODUCT, h.GetProductByID).Methods(shared.GET_METHOD)
}
