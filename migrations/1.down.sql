-- Drop junction tables first (due to foreign key constraints)
DROP TABLE IF EXISTS event_courses;

-- Drop tables with foreign keys to others
DROP TABLE IF EXISTS study_logs;
DROP TABLE IF EXISTS grades;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS courses;

-- Drop the parent table last
DROP TABLE IF EXISTS users;