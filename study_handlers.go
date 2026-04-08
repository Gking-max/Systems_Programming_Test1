package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
)

type StudyHandler struct {
	DB *sql.DB
}

func (h *StudyHandler) GetStudyLogs(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	rows, err := h.DB.Query(`
		SELECT s.log_id, s.course_id, s.start_time, s.end_time, s.duration_minutes, 
		       s.topic, s.notes, s.productivity_rating, s.created_at,
		       c.course_code, c.course_name
		FROM study_logs s
		LEFT JOIN courses c ON s.course_id = c.course_id
		WHERE s.user_id = $1
		ORDER BY s.start_time DESC LIMIT 50`, userID)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var logID int
		var courseID sql.NullInt64
		var startTime, endTime sql.NullTime
		var duration sql.NullInt64
		var topic, notes, courseCode, courseName sql.NullString
		var rating sql.NullInt64
		var createdAt time.Time
		
		rows.Scan(&logID, &courseID, &startTime, &endTime, &duration, &topic, &notes, &rating, &createdAt,
			&courseCode, &courseName)
		
		log := map[string]interface{}{
			"log_id": logID, "created_at": createdAt,
		}
		if courseID.Valid {
			log["course_id"] = courseID.Int64
		}
		if startTime.Valid {
			log["start_time"] = startTime.Time
		}
		if endTime.Valid {
			log["end_time"] = endTime.Time
		}
		if duration.Valid {
			log["duration_minutes"] = duration.Int64
		}
		if topic.Valid {
			log["topic"] = topic.String
		}
		if notes.Valid {
			log["notes"] = notes.String
		}
		if rating.Valid {
			log["productivity_rating"] = rating.Int64
		}
		if courseCode.Valid {
			log["course_code"] = courseCode.String
		}
		if courseName.Valid {
			log["course_name"] = courseName.String
		}
		logs = append(logs, log)
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Data: logs})
}

func (h *StudyHandler) CreateStudyLog(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	var log StudyLog
	if err := ParseJSONBody(r, &log); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	log.UserID = userID
	
	// Calculate duration if end_time is provided
	if log.EndTime != nil {
		duration := int(log.EndTime.Sub(log.StartTime).Minutes())
		log.DurationMinutes = &duration
	}

	err := h.DB.QueryRow(`
		INSERT INTO study_logs (user_id, course_id, start_time, end_time, duration_minutes, topic, notes, productivity_rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING log_id, created_at`,
		log.UserID, log.CourseID, log.StartTime, log.EndTime, log.DurationMinutes,
		log.Topic, log.Notes, log.ProductivityRating).Scan(&log.LogID, &log.CreatedAt)
	
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	SendJSONResponse(w, http.StatusCreated, Response{Success: true, Data: log, Message: "Study log created"})
}

func (h *StudyHandler) UpdateStudyLog(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	logID, _ := strconv.Atoi(idStr)

	var log StudyLog
	if err := ParseJSONBody(r, &log); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	if log.EndTime != nil {
		duration := int(log.EndTime.Sub(log.StartTime).Minutes())
		log.DurationMinutes = &duration
	}

	result, err := h.DB.Exec(`
		UPDATE study_logs SET course_id=$1, start_time=$2, end_time=$3, duration_minutes=$4,
		                      topic=$5, notes=$6, productivity_rating=$7
		WHERE log_id=$8 AND user_id=$9`,
		log.CourseID, log.StartTime, log.EndTime, log.DurationMinutes,
		log.Topic, log.Notes, log.ProductivityRating, logID, userID)
	
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Study log not found"})
		return
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Message: "Study log updated"})
}

func (h *StudyHandler) DeleteStudyLog(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	idStr := r.PathValue("id")
	logID, _ := strconv.Atoi(idStr)

	result, err := h.DB.Exec("DELETE FROM study_logs WHERE log_id=$1 AND user_id=$2", logID, userID)
	if err != nil {
		SendJSONResponse(w, http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		SendJSONResponse(w, http.StatusNotFound, Response{Success: false, Error: "Study log not found"})
		return
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Message: "Study log deleted"})
}

func (h *StudyHandler) GetStudyStats(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	
	var totalHours float64
	var avgRating float64
	var totalSessions int
	
	h.DB.QueryRow(`
		SELECT COALESCE(SUM(duration_minutes), 0)::float/60, 
		       COALESCE(AVG(productivity_rating), 0),
		       COUNT(*)
		FROM study_logs WHERE user_id = $1 AND start_time > NOW() - INTERVAL '30 days'`,
		userID).Scan(&totalHours, &avgRating, &totalSessions)

	stats := map[string]interface{}{
		"total_hours":    totalHours,
		"avg_rating":     avgRating,
		"total_sessions": totalSessions,
	}

	SendJSONResponse(w, http.StatusOK, Response{Success: true, Data: stats})
}