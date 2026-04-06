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