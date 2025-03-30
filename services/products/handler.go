package products

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MinaSamirSaad/ecommerce/services/shared"
)

type Handler struct {
	store ProductStore
}

func NewHandler(db *sql.DB) *Handler {
	ProductStore := NewStore(db)
	return &Handler{store: ProductStore}
}

func (h *Handler) GetProducts(res http.ResponseWriter, req *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		shared.RespondError(res, http.StatusInternalServerError, err)
		return
	}
	shared.WriteJson(res, http.StatusOK, products)
}

func (h *Handler) CreateProduct(res http.ResponseWriter, req *http.Request) {
	// get payload
	var payload shared.CreateProductPayload
	if err := shared.ParseJsonBody(req, &payload); err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
		return
	}
	// create product
	product, err := h.store.CreateProduct(&payload)
	if err != nil {
		shared.RespondError(res, http.StatusInternalServerError, err)
		return
	}
	shared.WriteJson(res, http.StatusCreated, product)
}

func (h *Handler) GetProductByID(res http.ResponseWriter, req *http.Request) {
	id, ok := shared.GetURLParam(req, "id")
	if !ok {
		shared.RespondError(res, http.StatusBadRequest, fmt.Errorf(shared.ErrMissingURLParam))
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
		return
	}
	product, err := h.store.GetProductByID(intID)
	if err != nil {
		shared.RespondError(res, http.StatusInternalServerError, err)
		return
	}
	if product == nil {
		shared.RespondError(res, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	shared.WriteJson(res, http.StatusOK, product)
}
