package main

import "time"

type User struct {
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	StudentID    string    `json:"student_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type Event struct {
	EventID     int       `json:"event_id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Course struct {
	CourseID   int       `json:"course_id"`
	UserID     int       `json:"user_id"`
	CourseCode string    `json:"course_code"`
	CourseName string    `json:"course_name"`
	Instructor string    `json:"instructor"`
	Credits    int       `json:"credits"`
	Room       string    `json:"room"`
	Schedule   string    `json:"schedule"`
	Color      string    `json:"color"`
	CreatedAt  time.Time `json:"created_at"`
}

type StudyLog struct {
	LogID             int        `json:"log_id"`
	UserID            int        `json:"user_id"`
	CourseID          *int       `json:"course_id"`
	StartTime         time.Time  `json:"start_time"`
	EndTime           *time.Time `json:"end_time"`
	DurationMinutes   *int       `json:"duration_minutes"`
	Topic             string     `json:"topic"`
	Notes             string     `json:"notes"`
	ProductivityRating *int      `json:"productivity_rating"`
	CreatedAt         time.Time  `json:"created_at"`
}

type Grade struct {
	GradeID        int        `json:"grade_id"`
	UserID         int        `json:"user_id"`
	CourseID       int        `json:"course_id"`
	AssignmentName string     `json:"assignment_name"`
	Grade          *float64   `json:"grade"`
	MaxGrade       *float64   `json:"max_grade"`
	Weight         *float64   `json:"weight"`
	DueDate        *time.Time `json:"due_date"`
	SubmittedDate  *time.Time `json:"submitted_date"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
}

type Todo struct {
	TodoID      int        `json:"todo_id"`
	UserID      int        `json:"user_id"`
	EventID     *int       `json:"event_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}