package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventID          = errors.New("invalid or empty id")
	ErrTimeEvent        = errors.New("invalid or empty event date")
	ErrEmptyRequestBody = errors.New("empty request body")
	ErrParseEvent       = errors.New("invalid or empty event")
)

type JSONResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Code  int         `json:"code"`
}

type ID struct {
	ID int32 `json:"id"`
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Errorln("empty request body")
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}

	event := storage.Event{}

	if err := parseEvent(&event, r); err != nil {
		h.log.Errorln("cannot parse event", err)
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
	}

	id, err := h.app.CreateEvent(r.Context(), &event)
	if err != nil {
		h.log.Errorln("failed to create event ", err)
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeOkResponse(w, ID{ID: id})
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Errorln("empty request body")
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}

	id, err := parseIDEvent(r)
	if err != nil {
		h.log.Errorln(err)
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
	}

	event := storage.Event{}

	if err := parseEvent(&event, r); err != nil {
		h.log.Errorln("cannot parse event ", err)
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.app.UpdateEvent(r.Context(), id, &event); err != nil {
		h.log.Errorln("failed to edit an event ", id, err.Error())
		if errors.Is(err, storage.ErrNotExistsEvent) {
			writeErrResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, nil)
}

func (h *EventHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.app.DeleteAllEvent(r.Context()); err != nil {
		h.log.Errorln("failed to remove an event ", err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeOkResponse(w, nil)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDEvent(r)
	if err != nil {
		h.log.Errorln(err)
		writeErrResponse(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.app.DeleteEvent(r.Context(), id); err != nil {
		h.log.Errorln("failed to delete an event", id, err.Error())
		if errors.Is(err, storage.ErrNotExistsEvent) {
			writeErrResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeOkResponse(w, id)
}

func (h *EventHandler) ListALl(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Errorln(ErrEmptyRequestBody)
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.app.ListAllEvents(r.Context())
	if err != nil {
		h.log.Errorln("failed to get all events ", err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeOkResponse(w, events)
}

func (h *EventHandler) ListMonth(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Errorln(ErrEmptyRequestBody)
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}

	date, err := parseDateEvent(r)
	if err != nil {
		h.log.Errorln(ErrTimeEvent)
		writeErrResponse(w, ErrTimeEvent.Error(), http.StatusInternalServerError)
	}

	events, err := h.app.ListMonthEvents(r.Context(), date)
	if err != nil {
		h.log.Errorln("failed to get events to month ", err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, events)
}

func (h *EventHandler) ListWeek(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Errorln("empty request body")
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}

	date, err := parseDateEvent(r)
	if err != nil {
		h.log.Errorln(ErrTimeEvent)
		writeErrResponse(w, ErrTimeEvent.Error(), http.StatusInternalServerError)
	}

	events, err := h.app.ListMonthEvents(r.Context(), date)
	if err != nil {
		h.log.Errorln("failed to get events to week ", err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, events)
}

func (h *EventHandler) ListDay(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Errorln("empty request body")
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}

	date, err := parseDateEvent(r)
	if err != nil {
		h.log.Errorln(ErrTimeEvent)
		writeErrResponse(w, ErrTimeEvent.Error(), http.StatusInternalServerError)
	}

	events, err := h.app.ListMonthEvents(r.Context(), date)
	if err != nil {
		h.log.Errorln("failed to get events to day ", err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeOkResponse(w, events)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Errorln("empty request body")
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}

	id, err := parseIDEvent(r)
	if err != nil {
		h.log.Errorln(ErrEventID)
		writeErrResponse(w, ErrEventID.Error(), http.StatusBadRequest)
	}

	event, err := h.app.GetEvent(r.Context(), id)
	if err != nil {
		h.log.Errorln("failed get event", id, err.Error())
		if errors.Is(err, storage.ErrNotExistsEvent) {
			writeErrResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeOkResponse(w, event)
}

func writeOkResponse(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := JSONResponse{Data: data, Code: http.StatusOK}
	_ = json.NewEncoder(w).Encode(response)
}

func writeErrResponse(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	response := JSONResponse{
		Error: err,
		Code:  status,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func parseEvent(e *storage.Event, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return err
	}
	if e == nil {
		return ErrParseEvent
	}
	return nil
}

func parseIDEvent(r *http.Request) (int32, error) {
	strId := chi.URLParam(r, "id")
	if strId == "" {
		return -1, ErrEventID
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		return -1, err
	}
	return int32(id), nil
}

func parseDateEvent(r *http.Request) (time.Time, error) {
	strStart := chi.URLParam(r, "start")
	if strStart == "" {
		return time.Time{}, ErrTimeEvent
	}
	return time.Parse(time.UnixDate, strStart)
}
