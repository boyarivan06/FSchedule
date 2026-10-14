package api

import (
	"FSchedule/database"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in Get User", "err", err)
		return
	}
	u, err := h.store.Users.GetById(r.Context(), id)
	if err != nil {
		h.logger.Error("user not found", "err", err)
		return
	}
	writeJSON(w, 200, u)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	q := r.URL.Query()
	if q.Has("username") {
		user.Username = q.Get("username")
		if _, err0 := h.store.Users.GetByUsername(r.Context(), user.Username); err0 == nil {
			writeError(w, http.StatusConflict, "such user is already created")
			return
		}
	} else {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	err := h.store.Users.Create(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating user")
		return
	}
	writeResponse(w, http.StatusCreated)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var params database.UserParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	var id int
	if _, err = fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in Update User", "err", err)
		return
	}
	var u database.User
	u, err = h.store.Users.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "no such user")
		return
	}
	if params.UserName != nil {
		if _, err0 := h.store.Users.GetByUsername(r.Context(), *params.UserName); err0 == nil {
			writeError(w, http.StatusBadRequest, "user exists")
		}
	}

	err = h.store.Users.Update(r.Context(), u.ID, params)
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
	writeResponse(w, http.StatusOK)
}
