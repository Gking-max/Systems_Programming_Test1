CREATE TABLE IF NOT EXISTS event_courses (
    event_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    PRIMARY KEY (event_id, course_id),
    FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);