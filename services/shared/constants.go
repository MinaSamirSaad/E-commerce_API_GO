package shared

const (
	POST_METHOD   = "POST"
	GET_METHOD    = "GET"
	PATCH_METHOD  = "PATCH"
	DELETE_METHOD = "DELETE"
)

const API_PREFIX = "/api"

// End points
const (
	LOGIN    = "/login"
	REGISTER = "/register"
	PRODUCTS = "/products"
	PRODUCT  = "/products/{id}"
	CHECKOUT = "/cart/checkout"
)

const ErrMissingURLParam = "missing URL parameter"
