package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

type EventHandler struct {
	DB *sql.DB
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	rows, err := h.DB.Query(`
		SELECT event_id, title, description, event_type, location, 
		       start_time, end_time, color, created_at, updated_at
		FROM events WHERE user_id = $1 ORDER BY start_time ASC`, userID)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.EventID, &e.Title, &e.Description, &e.EventType, &e.Location,
			&e.StartTime, &e.EndTime, &e.Color, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
			return
		}
		e.UserID = userID
		events = append(events, e)
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Data: events})
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	var event Event
	if err := ParseJSONBody(r, &event); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	event.UserID = userID
	if err := ValidateEvent(&event); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: err.Error()})
		return
	}

	err := h.DB.QueryRow(`
		INSERT INTO events (user_id, title, description, event_type, location, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING event_id, created_at, updated_at`,
		event.UserID, event.Title, event.Description, event.EventType,
		event.Location, event.StartTime, event.EndTime).Scan(&event.EventID, &event.CreatedAt, &event.UpdatedAt)
	
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	SendJSONResponse(w, http.StatusCreated, Response{Success: true, Data: event, Message: "Event created successfully"})
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	eventID, _ := strconv.Atoi(idStr)

	var event Event
	err := h.DB.QueryRow(`
		SELECT event_id, title, description, event_type, location, start_time, end_time, color, created_at, updated_at
		FROM events WHERE event_id = $1 AND user_id = $2`,
		eventID, userID).Scan(&event.EventID, &event.Title, &event.Description, &event.EventType,
		&event.Location, &event.StartTime, &event.EndTime, &event.Color, &event.CreatedAt, &event.UpdatedAt)
	
	if err == sql.ErrNoRows {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Event not found"})
		return
	}
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}
	event.UserID = userID

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Data: event})
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	eventID, _ := strconv.Atoi(idStr)

	var event Event
	if err := ParseJSONBody(r, &event); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	result, err := h.DB.Exec(`
		UPDATE events SET title=$1, description=$2, event_type=$3, location=$4, 
		                  start_time=$5, end_time=$6, updated_at=CURRENT_TIMESTAMP
		WHERE event_id=$7 AND user_id=$8`,
		event.Title, event.Description, event.EventType, event.Location,
		event.StartTime, event.EndTime, eventID, userID)
	
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Event not found"})
		return
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Message: "Event updated successfully"})
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	eventID, _ := strconv.Atoi(idStr)

	result, err := h.DB.Exec("DELETE FROM events WHERE event_id=$1 AND user_id=$2", eventID, userID)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Event not found"})
		return
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Message: "Event deleted successfully"})
}

func (h *EventHandler) GetUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	rows, err := h.DB.Query(`
		SELECT event_id, title, event_type, location, start_time, end_time
		FROM events 
		WHERE user_id = $1 AND start_time > NOW() 
		ORDER BY start_time ASC LIMIT 10`, userID)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		rows.Scan(&e.EventID, &e.Title, &e.EventType, &e.Location, &e.StartTime, &e.EndTime)
		events = append(events, e)
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Data: events})
}