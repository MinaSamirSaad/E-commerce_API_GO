package cart

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MinaSamirSaad/ecommerce/services/auth"
	"github.com/MinaSamirSaad/ecommerce/services/order"
	"github.com/MinaSamirSaad/ecommerce/services/products"
	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/MinaSamirSaad/ecommerce/services/users"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	orderStore   order.OrderStore
	productStore products.ProductStore
	userStore    shared.UserStore
}

func NewHandler(db *sql.DB) *Handler {
	OrderStore := order.NewStore(db)
	ProductStore := products.NewStore(db)
	UserStore := users.NewStore(db)
	return &Handler{orderStore: OrderStore, productStore: ProductStore, userStore: UserStore}
}

func (h *Handler) handleCheckout(res http.ResponseWriter, req *http.Request) {
	userID := auth.GetUserIDFromContext(req.Context())
	var cart shared.CartCheckoutPayload
	if err := shared.ParseJsonBody(req, &cart); err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
		return
	}

	if err := shared.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		shared.RespondError(res, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	productIds, err := GetCartItemsIDs(cart.Items)
	if err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
		return
	}

	// get products
	products, err := h.productStore.GetProductsByID(productIds)
	if err != nil {
		shared.RespondError(res, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(products, cart.Items, userID)
	if err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
		return
	}

	shared.WriteJson(res, http.StatusOK, map[string]interface{}{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}

func GetCartItemsIDs(items []shared.CartCheckoutItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productIds[i] = item.ProductID
	}

	return productIds, nil
}

func calculateTotalPrice(cartItems []shared.CartCheckoutItem, products map[int]shared.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}

func checkIfCartIsInStock(cartItems []shared.CartCheckoutItem, products map[int]shared.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func (h *Handler) createOrder(products []shared.Product, cartItems []shared.CartCheckoutItem, userID int) (int, float64, error) {
	// create a map of products for easier access
	productsMap := make(map[int]shared.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	// check if all products are available
	if err := checkIfCartIsInStock(cartItems, productsMap); err != nil {
		return 0, 0, err
	}

	// calculate total price
	totalPrice := calculateTotalPrice(cartItems, productsMap)

	// reduce the quantity of products in the store
	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProduct(product)
	}

	// create order record
	orderID, err := h.orderStore.CreateOrder(shared.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address", // could fetch address from a user addresses table
	})
	if err != nil {
		return 0, 0, err
	}

	// create order the items records
	for _, item := range cartItems {
		h.orderStore.CreateOrderItem(shared.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}
