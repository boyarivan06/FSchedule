package api

import (
	"FSchedule/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	var params database.EventParams
	q := r.URL.Query()

	if q.Has("name") {
		params.Name = new(q.Get("name"))
	}

	if q.Has("start_at") {
		t, err := time.Parse(time.RFC3339, q.Get("start_at"))
		if err != nil {
			http.Error(w, "invalid start_at", http.StatusBadRequest)
			return
		}
		params.StartAt = new(t)
	}

	if q.Has("end_at") {
		t, err := time.Parse(time.RFC3339, q.Get("end_at"))
		if err != nil {
			http.Error(w, "invalid end_at", http.StatusBadRequest)
			return
		}
		params.EndAt = new(t)
	}

	if q.Has("repeat") {
		v, err := strconv.ParseInt(q.Get("repeat"), 10, 64)
		if err != nil {
			http.Error(w, "invalid repeat", http.StatusBadRequest)
			return
		}
		params.Repeat = new(v)
	}

	if q.Has("child_id") {
		v, err := strconv.ParseInt(q.Get("child_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid child_id", http.StatusBadRequest)
			return
		}
		params.ChildID = new(v)
	}

	if params.Name == nil && params.StartAt == nil && params.EndAt == nil &&
		params.Repeat == nil && params.ChildID == nil {
		http.Error(w, "nothing to update", http.StatusBadRequest)
		return
	}

	events, err := h.store.Events.GetManyByParams(r.Context(),
		params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error getting events")
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in get event", "err", err)
		return
	}
	event, err := h.store.Events.GetById(r.Context(), id)
	if err != nil {
		h.logger.Error("event not found", "err", err)
		return
	}
	writeJSON(w, 200, event)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e database.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	event, err := h.store.Events.Create(r.Context(), e)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating event")
		return
	}
	writeJSON(w, http.StatusCreated, event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in update event", "err", err)
		return
	}
	var params database.EventParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	var c database.Event
	c, err = h.store.Events.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "no such event")
		return
	}
	err = h.store.Events.Update(r.Context(), c.ID, params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error updating event")
		return
	}
	writeResponse(w, http.StatusCreated)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Error("no id given in delete event", "err", err)
		return
	}
	var _, err = h.store.Events.GetById(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "no such event")
		return
	}
	err = h.store.Events.Delete(r.Context(), int64(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error deleting event")
		return
	}
	writeResponse(w, http.StatusOK)
}
