//go:build integration

package api

import (
	"FSchedule/database"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func setupTestHandler(t *testing.T) *Handler {
	storage := database.NewStorage(database.SetupTestDb(t))
	handler := NewHandler(storage, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	return handler
}

func newTestRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/users/{id}", h.GetUser)
	r.Post("/users", h.CreateUser)
	r.Patch("/users/{id}", h.UpdateUser)
	r.Delete("/users/{id}", h.DeleteUser)
	return r
}

func TestHandler_CreateUser(t *testing.T) {
	h := setupTestHandler(t)
	r := newTestRouter(h)
	// usual create
	req := httptest.NewRequest(http.MethodPost, "/users?username=biba", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, 201, rec.Code)
	var u database.User
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&u))
	require.Equal(t, "biba", u.Username)
	// taken username
	req = httptest.NewRequest(http.MethodPost, "/users?username=biba", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, 409, rec.Code)
	//	require.Error(t, json.NewDecoder(rec.Body).Decode(&u))
	// sql-injection
	req = httptest.NewRequest(http.MethodPost, "/users?username=biba'", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, 500, rec.Code)
}

func TestHandler_GetUser(t *testing.T) {

}
