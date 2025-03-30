package users

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MinaSamirSaad/ecommerce/config"
	"github.com/MinaSamirSaad/ecommerce/services/auth"
	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store shared.UserStore
}

func NewHandler(db *sql.DB) *Handler {
	UserStore := NewStore(db)
	return &Handler{store: UserStore}
}

func (h *Handler) handleLogin(res http.ResponseWriter, req *http.Request) {
	// get payload
	var payload shared.LoginUserPayload
	if err := shared.ParseJsonBody(req, &payload); err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
		return
	}
	// check if user exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		shared.RespondError(res, http.StatusBadRequest, fmt.Errorf("not found invalid email or password"))
		return
	}

	// compare password
	if err := auth.ComparePassword(u.Password, payload.Password); err != nil {
		shared.RespondError(res, http.StatusBadRequest, fmt.Errorf("not found invalid email or password"))
		return
	}
	// generate token
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		shared.RespondError(res, http.StatusInternalServerError, err)
		return
	}
	// return response
	response := shared.LoginUserResponse{
		User:  *u,
		Token: token,
	}
	shared.WriteJson(res, http.StatusOK, response)
}

func (h *Handler) handleRegister(res http.ResponseWriter, req *http.Request) {
	// get payload
	var payload shared.RegisterUserPayload
	if err := shared.ParseJsonBody(req, &payload); err != nil {
		shared.RespondError(res, http.StatusBadRequest, err)
		return
	}
	// check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		shared.RespondError(res, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	// validate payload
	if err := shared.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		shared.RespondError(res, http.StatusBadRequest, error)
		return
	}
	// hash password
	payload.Password, err = auth.HashPassword(payload.Password)
	if err != nil {
		shared.RespondError(res, http.StatusInternalServerError, err)
		return
	}
	// create user
	newUser := &shared.User{
		Email:     payload.Email,
		Password:  payload.Password,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}
	newUser, err = h.store.CreateUser(newUser)
	if err != nil {
		shared.RespondError(res, http.StatusInternalServerError, fmt.Errorf("failed to create user: %v", err))
		return
	}
	// return response
	shared.WriteJson(res, http.StatusCreated, newUser)
}
