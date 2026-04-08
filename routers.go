package main

import (
	"net/http"
)

func setupRoutes(eventHandler *EventHandler, courseHandler *CourseHandler, studyHandler *StudyHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Event routes
	mux.HandleFunc("GET /api/events", eventHandler.GetEvents)
	mux.HandleFunc("POST /api/events", eventHandler.CreateEvent)
	mux.HandleFunc("GET /api/events/{id}", eventHandler.GetEvent)
	mux.HandleFunc("PUT /api/events/{id}", eventHandler.UpdateEvent)
	mux.HandleFunc("DELETE /api/events/{id}", eventHandler.DeleteEvent)
	mux.HandleFunc("GET /api/events/upcoming", eventHandler.GetUpcomingEvents)

	// Course routes
	mux.HandleFunc("GET /api/courses", courseHandler.GetCourses)
	mux.HandleFunc("POST /api/courses", courseHandler.CreateCourse)
	mux.HandleFunc("GET /api/courses/{id}", courseHandler.GetCourse)
	mux.HandleFunc("PUT /api/courses/{id}", courseHandler.UpdateCourse)
	mux.HandleFunc("DELETE /api/courses/{id}", courseHandler.DeleteCourse)
	mux.HandleFunc("GET /api/courses/{id}/events", courseHandler.GetCourseEvents)

	// Study routes
	mux.HandleFunc("GET /api/study-logs", studyHandler.GetStudyLogs)
	mux.HandleFunc("POST /api/study-logs", studyHandler.CreateStudyLog)
	mux.HandleFunc("PUT /api/study-logs/{id}", studyHandler.UpdateStudyLog)
	mux.HandleFunc("DELETE /api/study-logs/{id}", studyHandler.DeleteStudyLog)
	mux.HandleFunc("GET /api/study-stats", studyHandler.GetStudyStats)

	// Grade routes
	mux.HandleFunc("GET /api/grades", getGrades)
	mux.HandleFunc("POST /api/grades", createGrade)
	mux.HandleFunc("PUT /api/grades/{id}", updateGrade)
	mux.HandleFunc("DELETE /api/grades/{id}", deleteGrade)
	mux.HandleFunc("GET /api/grade-stats", getGradeStats)

	// Todo routes
	mux.HandleFunc("GET /api/todos", getTodos)
	mux.HandleFunc("POST /api/todos", createTodo)
	mux.HandleFunc("PUT /api/todos/{id}", updateTodo)
	mux.HandleFunc("DELETE /api/todos/{id}", deleteTodo)
	mux.HandleFunc("PUT /api/todos/{id}/complete", completeTodo)

	return mux
}