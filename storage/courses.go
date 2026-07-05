package storage

import (
	"database/sql"
	"errors"
)

// SaveCourse adds a new course to the database.
func (s *SQLiteEngine) SaveCourse(course *Course) error {
	query := `
	INSERT INTO courses (
		code,
		title,
		units,
		prerequisite_code,
		lecturer_id,
		lecturer_name
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		course.Code,
		course.Title,
		course.Units,
		course.PrerequisiteCode,
		course.LecturerID,
		course.LecturerName,
	)

	return err
}

// GetCourse returns a course using its course code.
func (s *SQLiteEngine) GetCourse(code string) (*Course, error) {
	query := `
	SELECT
		code,
		title,
		units,
		prerequisite_code,
		lecturer_id,
		lecturer_name
	FROM courses
	WHERE code = ?
	`

	row := s.db.QueryRow(query, code)

	var course Course

	err := row.Scan(
		&course.Code,
		&course.Title,
		&course.Units,
		&course.PrerequisiteCode,
		&course.LecturerID,
		&course.LecturerName,
	)

	// Return nil if the course does not exist.
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &course, nil
}

// GetAllCourses loads every course from the database.
func (s *SQLiteEngine) GetAllCourses() ([]Course, error) {
	query := `
	SELECT
		code,
		title,
		units,
		prerequisite_code,
		lecturer_id,
		lecturer_name
	FROM courses
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course

	for rows.Next() {
		var course Course

		err := rows.Scan(
			&course.Code,
			&course.Title,
			&course.Units,
			&course.PrerequisiteCode,
			&course.LecturerID,
			&course.LecturerName,
		)

		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	return courses, nil
}
