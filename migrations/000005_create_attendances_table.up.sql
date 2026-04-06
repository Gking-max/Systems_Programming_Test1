CREATE TABLE IF NOT EXISTS attendances (
    attendance_id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    status VARCHAR(20) CHECK (status IN ('present', 'absent', 'late', 'excused')) DEFAULT 'present',
    notes TEXT,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE,
    UNIQUE(event_id)
);