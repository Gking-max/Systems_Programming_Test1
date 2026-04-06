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