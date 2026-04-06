package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    // Database connection string
    dbURL := "postgresql://postgres:postgres@localhost:5432/student_event_tracker?sslmode=disable"
    
    // Connect to PostgreSQL
    db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Create database if it doesn't exist
    _, err = db.Exec("CREATE DATABASE student_event_tracker")
    if err != nil {
        fmt.Println("Database may already exist:", err)
    }
    
    // Reconnect to the specific database
    db.Close()
    db, err = sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Create migration driver
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Initialize migrations
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres", driver)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run migrations
    fmt.Println("Running migrations...")
    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatal("Migration failed:", err)
    }
    
    fmt.Println("✅ Migrations completed successfully!")
    
    // Verify tables
    verifyTables(db)
    
    // Insert sample data
    insertSampleData(db)
    
    // Run queries
    runQueries(db)
}

func verifyTables(db *sql.DB) {
    tables := []string{"users", "events", "courses", "event_courses", "attendances", 
                       "reminders", "study_logs", "grades", "todos", "notifications"}
    
    fmt.Println("\n📋 Verifying tables:")
    for _, table := range tables {
        var exists bool
        err := db.QueryRow(`
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_name = $1
            )`, table).Scan(&exists)
        
        if err != nil {
            fmt.Printf("❌ Error checking %s: %v\n", table, err)
            continue
        }
        
        if exists {
            fmt.Printf("✅ Table '%s' created\n", table)
        } else {
            fmt.Printf("❌ Table '%s' not found\n", table)
        }
    }
}

func insertSampleData(db *sql.DB) {
    fmt.Println("\n📝 Inserting sample data...")
    
    // Insert a sample user
    var userID int
    err := db.QueryRow(`
        INSERT INTO users (username, email, password_hash, full_name, student_id)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (username) DO NOTHING
        RETURNING user_id`,
        "john_doe", "john@university.edu", "hashed_password_123", "John Doe", "STU12345",
    ).Scan(&userID)
    
    if err != nil {
        fmt.Println("User insertion error (may already exist):", err)
        // Get existing user ID
        db.QueryRow("SELECT user_id FROM users WHERE username = 'john_doe'").Scan(&userID)
    }
    
    // Insert sample courses
    courses := []struct {
        code, name, instructor, room, schedule string
        credits int
        color string
    }{
        {"CS101", "Introduction to Programming", "Prof. Smith", "Room 201", "Mon/Wed 10:00-11:30", 3, "#3498db"},
        {"MATH201", "Calculus II", "Prof. Johnson", "Room 105", "Tue/Thu 13:00-14:30", 4, "#2ecc71"},
        {"PHY101", "Physics Fundamentals", "Prof. Williams", "Lab 3", "Mon/Wed 15:00-16:30", 3, "#e74c3c"},
    }
    
    var courseIDs []int
    for _, c := range courses {
        var courseID int
        err := db.QueryRow(`
            INSERT INTO courses (user_id, course_code, course_name, instructor, credits, room, schedule, color)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            ON CONFLICT (user_id, course_code) DO UPDATE SET course_name = EXCLUDED.course_name
            RETURNING course_id`,
            userID, c.code, c.name, c.instructor, c.credits, c.room, c.schedule, c.color,
        ).Scan(&courseID)
        if err == nil {
            courseIDs = append(courseIDs, courseID)
        }
    }
    
    // Insert sample events
    events := []struct {
        title, description, eventType, location string
        startTime, endTime time.Time
    }{
        {"CS101 Lecture", "Introduction to variables and data types", "class", "Room 201", 
            time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 1).Add(2 * time.Hour)},
        {"MATH201 Exam", "Midterm Exam - Chapters 1-5", "exam", "Room 105", 
            time.Now().AddDate(0, 0, 7), time.Now().AddDate(0, 0, 7).Add(3 * time.Hour)},
        {"Study Session", "Group study for Physics", "study", "Library Room 5", 
            time.Now().AddDate(0, 0, 3), time.Now().AddDate(0, 0, 3).Add(3 * time.Hour)},
    }
    
    var eventIDs []int
    for _, e := range events {
        var eventID int
        err := db.QueryRow(`
            INSERT INTO events (user_id, title, description, event_type, location, start_time, end_time)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING event_id`,
            userID, e.title, e.description, e.eventType, e.location, e.startTime, e.endTime,
        ).Scan(&eventID)
        if err == nil {
            eventIDs = append(eventIDs, eventID)
        }
    }
    
    // Insert sample todos
    todos := []struct {
        title, description string
        priority string
        dueDate time.Time
    }{
        {"Complete CS101 Assignment", "Finish coding exercises", "high", time.Now().AddDate(0, 0, 3)},
        {"Read Chapter 4 for MATH201", "Review calculus concepts", "medium", time.Now().AddDate(0, 0, 2)},
        {"Submit Physics Lab Report", "Format according to guidelines", "high", time.Now().AddDate(0, 0, 5)},
    }
    
    for _, t := range todos {
        _, err := db.Exec(`
            INSERT INTO todos (user_id, title, description, priority, due_date)
            VALUES ($1, $2, $3, $4, $5)`,
            userID, t.title, t.description, t.priority, t.dueDate,
        )
        if err != nil {
            fmt.Printf("Error inserting todo: %v\n", err)
        }
    }
    
    fmt.Println("✅ Sample data inserted successfully!")
}

func runQueries(db *sql.DB) {
    fmt.Println("\n🔍 Running sample queries:")
    
    // Query 1: Get all upcoming events
    fmt.Println("\n1. Upcoming Events:")
    rows, err := db.Query(`
        SELECT title, event_type, start_time, location 
        FROM events 
        WHERE start_time > NOW() 
        ORDER BY start_time ASC 
        LIMIT 5
    `)
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var title, eventType, location string
            var startTime time.Time
            rows.Scan(&title, &eventType, &startTime, &location)
            fmt.Printf("   📌 %s - %s at %s (%s)\n", title, eventType, location, startTime.Format("Jan 2, 15:04"))
        }
    }
    
    // Query 2: Get courses with upcoming events
    fmt.Println("\n2. Courses with upcoming events:")
    rows, err = db.Query(`
        SELECT DISTINCT c.course_code, c.course_name, COUNT(e.event_id) as event_count
        FROM courses c
        LEFT JOIN event_courses ec ON c.course_id = ec.course_id
        LEFT JOIN events e ON ec.event_id = e.event_id AND e.start_time > NOW()
        WHERE c.user_id = (SELECT user_id FROM users WHERE username = 'john_doe' LIMIT 1)
        GROUP BY c.course_id, c.course_code, c.course_name
        HAVING COUNT(e.event_id) > 0
    `)
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var courseCode, courseName string
            var eventCount int
            rows.Scan(&courseCode, &courseName, &eventCount)
            fmt.Printf("   📚 %s - %s (%d upcoming events)\n", courseCode, courseName, eventCount)
        }
    }
    
    // Query 3: Get pending high priority todos
    fmt.Println("\n3. High Priority Tasks:")
    rows, err = db.Query(`
        SELECT title, due_date, priority
        FROM todos
        WHERE is_completed = false AND priority = 'high'
        ORDER BY due_date ASC
    `)
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var title, priority string
            var dueDate time.Time
            rows.Scan(&title, &dueDate, &priority)
            fmt.Printf("   ⚠️  %s - Due: %s\n", title, dueDate.Format("Jan 2, 15:04"))
        }
    }
    
    // Query 4: Get study logs summary
    fmt.Println("\n4. Study Logs Summary:")
    var totalMinutes, avgRating int
    err = db.QueryRow(`
        SELECT COALESCE(SUM(duration_minutes), 0), COALESCE(AVG(productivity_rating)::int, 0)
        FROM study_logs
        WHERE start_time > NOW() - INTERVAL '7 days'
    `).Scan(&totalMinutes, &avgRating)
    if err == nil {
        fmt.Printf("   📖 Total study time (last 7 days): %d hours %d minutes\n", totalMinutes/60, totalMinutes%60)
        fmt.Printf("   ⭐ Average productivity rating: %d/5\n", avgRating)
    }
    
    // Query 5: Get grade summary
    fmt.Println("\n5. Grade Summary:")
    rows, err = db.Query(`
        SELECT c.course_code, c.course_name, 
               ROUND(AVG(g.grade/g.max_grade * 100), 1) as average_score
        FROM grades g
        JOIN courses c ON g.course_id = c.course_id
        WHERE g.grade IS NOT NULL
        GROUP BY c.course_id, c.course_code, c.course_name
    `)
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var courseCode, courseName string
            var avgScore float64
            rows.Scan(&courseCode, &courseName, &avgScore)
            fmt.Printf("   📊 %s - %s: %.1f%%\n", courseCode, courseName, avgScore)
        }
    }
    
    fmt.Println("\n✅ All queries completed successfully!")
}