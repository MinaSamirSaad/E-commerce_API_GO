package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MinaSamirSaad/ecommerce/services/shared"
	"github.com/gorilla/mux"
)

type mockUserStore struct{}

func (m *mockUserStore) CreateUser(u *shared.User) (*shared.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*shared.User, error) {
	return nil, fmt.Errorf("user with email %s already exists", email)
}

func (m *mockUserStore) GetUserByID(id int) (*shared.User, error) {
	return nil, nil
}

func TestUserServiceHandlers(t *testing.T) {
	mockStore := &mockUserStore{}
	Handler := Handler{store: mockStore}
	t.Run(("should fail it it is not valid email"), func(t *testing.T) {
		payload := &shared.RegisterUserPayload{
			Email:     "notValid",
			Password:  "123456",
			FirstName: "test",
			LastName:  "test",
		}
		marshalPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalPayload))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", Handler.handleRegister)
		router.ServeHTTP(rr, req)
		fmt.Println(rr.Body.String())
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status code 400, got %d", rr.Code)
		}
	})
	t.Run(("should pass if email valid"), func(t *testing.T) {
		payload := &shared.RegisterUserPayload{
			Email:     "valid@test.com",
			Password:  "123456",
			FirstName: "test",
			LastName:  "test",
		}
		marshalPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalPayload))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", Handler.handleRegister)
		router.ServeHTTP(rr, req)
		fmt.Println(rr.Body.String())
		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status code 201, got %d", rr.Code)
		}
	})
}
