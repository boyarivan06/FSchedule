package api

import (
	"FSchedule/database"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store  *database.Storage
	logger *slog.Logger
}

func NewHandler(store *database.Storage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) GetUser(w http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in Get User", "err", err)
		return
	}
	u, err := h.store.Users.GetById(request.Context(), id)
	if err != nil {
		h.logger.Error("user not found", "err", err)
		return
	}
	writeJSON(w, 200, u)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u database.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if u.Username == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	} else if _, err0 := h.store.Users.GetByUsername(r.Context(), u.Username); err0 == nil {
		writeError(w, http.StatusConflict, "such user is already created")
		return
	}
	err := h.store.Users.Create(r.Context(), u.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating user")
		return
	}
	writeResponse(w, http.StatusCreated)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in Update User", "err", err)
		return
	}
	var u, err = h.store.Users.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "no such user")
		return
	}
	var newUser database.User
	if err = json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	err = h.store.Users.Update(r.Context(), u.ID, map[string]string{"username": newUser.Username})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error updating user")
		return
	}
	writeResponse(w, http.StatusCreated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in delete user", "err", err)
		return
	}
	var _, err = h.store.Users.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "no such user")
		return
	}
	err = h.store.Users.Delete(r.Context(), int64(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error deleting user")
		return
	}
}
