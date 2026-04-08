package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func ParseJSONBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func GetUserID(r *http.Request) int {
	// In production, extract from JWT or session
	// For demo, return default user ID
	return 1
}

func ValidateEvent(event *Event) error {
	if event.Title == "" {
		return fmt.Errorf("title is required")
	}
	if event.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if event.StartTime.IsZero() || event.EndTime.IsZero() {
		return fmt.Errorf("start and end times are required")
	}
	if event.EndTime.Before(event.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}
	return nil
}

func ValidateCourse(course *Course) error {
	if course.CourseCode == "" {
		return fmt.Errorf("course code is required")
	}
	if course.CourseName == "" {
		return fmt.Errorf("course name is required")
	}
	return nil
}

func CalculateDuration(start, end time.Time) int {
	return int(end.Sub(start).Minutes())
}