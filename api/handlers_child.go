package api

import (
	"FSchedule/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetChildren(w http.ResponseWriter, r *http.Request) {
	var params database.ChildParams
	q := r.URL.Query()

	if q.Has("name") {
		params.Name = new(q.Get("name"))
	}
	if q.Has("user_id") {
		userID, err := strconv.Atoi(q.Get("user_id"))
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		params.UserID = new(int64(userID))
	}

	children, err := h.store.Children.GetManyByParams(r.Context(),
		params)
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusInternalServerError, "error getting children")
		return
	}
	writeJSON(w, http.StatusOK, children)
}

func (h *Handler) GetChild(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in get child", "err", err)
		return
	}
	child, err := h.store.Children.GetById(r.Context(), id)
	if err != nil {
		h.logger.Error("child not found", "err", err)
		return
	}
	writeJSON(w, 200, child)
}

func (h *Handler) CreateChild(w http.ResponseWriter, r *http.Request) {
	var c database.Child
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	err := h.store.Children.Create(r.Context(), c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating child")
		return
	}
	writeResponse(w, http.StatusCreated)
}

func (h *Handler) UpdateChild(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in update child", "err", err)
		return
	}
	var params database.ChildParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	var c database.Child
	c, err = h.store.Children.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "no such child")
		return
	}
	err = h.store.Children.Update(r.Context(), c.ID, params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error updating child")
		return
	}
	writeResponse(w, http.StatusCreated)
}

func (h *Handler) DeleteChild(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in delete child", "err", err)
		return
	}
	var _, err = h.store.Children.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "no such child")
		return
	}
	err = h.store.Children.Delete(r.Context(), int64(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error deleting child")
		return
	}
	writeResponse(w, http.StatusOK)
}
