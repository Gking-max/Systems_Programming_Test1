package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

type CourseHandler struct {
	DB *sql.DB
}

func (h *CourseHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	rows, err := h.DB.Query(`
		SELECT course_id, course_code, course_name, instructor, credits, room, schedule, color, created_at
		FROM courses WHERE user_id = $1 ORDER BY course_code`, userID)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		err := rows.Scan(&c.CourseID, &c.CourseCode, &c.CourseName, &c.Instructor,
			&c.Credits, &c.Room, &c.Schedule, &c.Color, &c.CreatedAt)
		if err != nil {
			SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
			return
		}
		c.UserID = userID
		courses = append(courses, c)
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Data: courses})
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	var course Course
	if err := ParseJSONBody(r, &course); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	course.UserID = userID
	if err := ValidateCourse(&course); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: err.Error()})
		return
	}

	err := h.DB.QueryRow(`
		INSERT INTO courses (user_id, course_code, course_name, instructor, credits, room, schedule, color)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING course_id, created_at`,
		course.UserID, course.CourseCode, course.CourseName, course.Instructor,
		course.Credits, course.Room, course.Schedule, course.Color).Scan(&course.CourseID, &course.CreatedAt)
	
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	SendJSONResponse(w, http.StatusCreated, Response{Success: true, Data: course, Message: "Course created successfully"})
}

func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	courseID, _ := strconv.Atoi(idStr)

	var course Course
	err := h.DB.QueryRow(`
		SELECT course_id, course_code, course_name, instructor, credits, room, schedule, color, created_at
		FROM courses WHERE course_id = $1 AND user_id = $2`,
		courseID, userID).Scan(&course.CourseID, &course.CourseCode, &course.CourseName,
		&course.Instructor, &course.Credits, &course.Room, &course.Schedule, &course.Color, &course.CreatedAt)
	
	if err == sql.ErrNoRows {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Course not found"})
		return
	}
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}
	course.UserID = userID

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Data: course})
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	courseID, _ := strconv.Atoi(idStr)

	var course Course
	if err := ParseJSONBody(r, &course); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	result, err := h.DB.Exec(`
		UPDATE courses SET course_code=$1, course_name=$2, instructor=$3, credits=$4, 
		                   room=$5, schedule=$6, color=$7
		WHERE course_id=$8 AND user_id=$9`,
		course.CourseCode, course.CourseName, course.Instructor, course.Credits,
		course.Room, course.Schedule, course.Color, courseID, userID)
	
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Course not found"})
		return
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Message: "Course updated successfully"})
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	courseID, _ := strconv.Atoi(idStr)

	result, err := h.DB.Exec("DELETE FROM courses WHERE course_id=$1 AND user_id=$2", courseID, userID)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Course not found"})
		return
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Message: "Course deleted successfully"})
}

func (h *CourseHandler) GetCourseEvents(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	courseID, _ := strconv.Atoi(idStr)

	rows, err := h.DB.Query(`
		SELECT e.event_id, e.title, e.event_type, e.location, e.start_time, e.end_time
		FROM events e
		JOIN event_courses ec ON e.event_id = ec.event_id
		WHERE ec.course_id = $1 AND e.user_id = $2
		ORDER BY e.start_time ASC`, courseID, userID)
	
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