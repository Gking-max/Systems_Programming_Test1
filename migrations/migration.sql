-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    student_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Events table
CREATE TABLE IF NOT EXISTS events (
    event_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    event_type VARCHAR(50) CHECK (event_type IN ('class', 'exam', 'study', 'assignment', 'extracurricular', 'personal')) NOT NULL,
    location VARCHAR(200),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    color VARCHAR(7) DEFAULT '#3498db',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Courses table
CREATE TABLE IF NOT EXISTS courses (
    course_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    course_code VARCHAR(20) NOT NULL,
    course_name VARCHAR(200) NOT NULL,
    instructor VARCHAR(100),
    credits INTEGER DEFAULT 3,
    room VARCHAR(50),
    schedule VARCHAR(200),
    color VARCHAR(7) DEFAULT '#2ecc71',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    UNIQUE(user_id, course_code)
);

-- Event_Courses junction table
CREATE TABLE IF NOT EXISTS event_courses (
    event_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    PRIMARY KEY (event_id, course_id),
    FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

-- Study logs table
CREATE TABLE IF NOT EXISTS study_logs (
    log_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    course_id INTEGER,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_minutes INTEGER,
    topic VARCHAR(200),
    notes TEXT,
    productivity_rating INTEGER CHECK (productivity_rating BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE SET NULL
);

-- Grades table
CREATE TABLE IF NOT EXISTS grades (
    grade_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    assignment_name VARCHAR(200) NOT NULL,
    grade DECIMAL(5,2),
    max_grade DECIMAL(5,2),
    weight DECIMAL(5,2),
    due_date DATE,
    submitted_date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

-- Todos table
CREATE TABLE IF NOT EXISTS todos (
    todo_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    event_id INTEGER,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    is_completed BOOLEAN DEFAULT FALSE,
    priority VARCHAR(20) CHECK (priority IN ('low', 'medium', 'high')) DEFAULT 'medium',
    due_date TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE SET NULL
);

-- Insert sample user
INSERT INTO users (username, email, password_hash, full_name, student_id) 
VALUES ('john_doe', 'john@university.edu', 'hashed_password_123', 'John Doe', 'STU12345')
ON CONFLICT (username) DO NOTHING;

-- Insert sample courses
INSERT INTO courses (user_id, course_code, course_name, instructor, credits, room, schedule, color) 
VALUES 
    ((SELECT user_id FROM users WHERE username = 'john_doe'), 'CS101', 'Introduction to Programming', 'Prof. Smith', 3, 'Room 201', 'Mon/Wed 10:00-11:30', '#3498db'),
    ((SELECT user_id FROM users WHERE username = 'john_doe'), 'MATH201', 'Calculus II', 'Prof. Johnson', 4, 'Room 105', 'Tue/Thu 13:00-14:30', '#2ecc71'),
    ((SELECT user_id FROM users WHERE username = 'john_doe'), 'PHY101', 'Physics Fundamentals', 'Prof. Williams', 3, 'Lab 3', 'Mon/Wed 15:00-16:30', '#e74c3c')
ON CONFLICT (user_id, course_code) DO NOTHING;

-- Insert sample events
INSERT INTO events (user_id, title, description, event_type, location, start_time, end_time) 
VALUES 
    ((SELECT user_id FROM users WHERE username = 'john_doe'), 'CS101 Lecture', 'Introduction to variables', 'class', 'Room 201', NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day' + INTERVAL '2 hours'),
    ((SELECT user_id FROM users WHERE username = 'john_doe'), 'Math Midterm', 'Chapters 1-5', 'exam', 'Room 105', NOW() + INTERVAL '7 days', NOW() + INTERVAL '7 days' + INTERVAL '3 hours')
ON CONFLICT DO NOTHING;

-- Drop junction tables first
DROP TABLE IF EXISTS event_courses;

-- Drop tables with foreign keys to others
DROP TABLE IF EXISTS study_logs;
DROP TABLE IF EXISTS grades;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS courses;

-- Drop the parent table last
DROP TABLE IF EXISTS users;
