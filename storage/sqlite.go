package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteEngine handles all SQLite database operations.
type SQLiteEngine struct {
	db *sql.DB
}

// NewSQLiteEngine opens the database connection
// and prepares required tables.
func NewSQLiteEngine(dbPath string) (*SQLiteEngine, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	engine := &SQLiteEngine{
		db: db,
	}

	// Create tables if they do not exist
	if err := engine.createTables(); err != nil {
		db.Close()
		return nil, err
	}

	return engine, nil
}

// Close shuts down the database connection.
func (s *SQLiteEngine) Close() error {
	return s.db.Close()
}

// createTables creates all required database tables.
func (s *SQLiteEngine) createTables() error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			surname TEXT NOT NULL,
			first_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL,
			security_question TEXT,
			security_answer TEXT
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS courses (
			code TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			units INTEGER NOT NULL,
			prerequisite_code TEXT,
			lecturer_id TEXT NOT NULL,
			lecturer_name TEXT NOT NULL
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_id TEXT NOT NULL,
			course_code TEXT NOT NULL,
			score INTEGER NOT NULL,
			grade TEXT NOT NULL,
			status TEXT NOT NULL,

			FOREIGN KEY(student_id)
				REFERENCES users(id),

			FOREIGN KEY(course_code)
				REFERENCES courses(code)
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS notices (
			id TEXT PRIMARY KEY,
			author_id TEXT NOT NULL,
			author TEXT NOT NULL,
			target_department TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS registrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_id TEXT NOT NULL,
			course_code TEXT NOT NULL,

			UNIQUE(student_id, course_code),

			FOREIGN KEY(student_id)
				REFERENCES users(id),

			FOREIGN KEY(course_code)
				REFERENCES courses(code)
		);
		`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create tables: %w", err)
		}
	}

	return nil
}

