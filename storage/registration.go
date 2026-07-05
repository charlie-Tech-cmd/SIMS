package storage

// RegisterStudentForCourse saves a student's course registration.
func (s *SQLiteEngine) RegisterStudentForCourse(studentID, courseCode string) error {
	query := `
	INSERT INTO registrations (
		student_id,
		course_code
	)
	VALUES (?, ?)
	`

	_, err := s.db.Exec(query, studentID, courseCode)
	return err
}

// GetStudentRegisteredCourses returns all courses
// a student has registered for.
func (s *SQLiteEngine) GetStudentRegisteredCourses(studentID string) ([]Course, error) {
	query := `
	SELECT
		c.code,
		c.title,
		c.units,
		c.prerequisite_code,
		c.lecturer_id,
		c.lecturer_name
	FROM registrations r
	JOIN courses c
		ON r.course_code = c.code
	WHERE r.student_id = ?
	`

	rows, err := s.db.Query(query, studentID)
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

// GetStudentRegisteredUnits calculates the total
// registered units for a student.
func (s *SQLiteEngine) GetStudentRegisteredUnits(studentID string) (int, error) {
	query := `
	SELECT TOTAL(c.units)
	FROM registrations r
	JOIN courses c
		ON r.course_code = c.code
	WHERE r.student_id = ?
	`

	var totalUnits int

	err := s.db.QueryRow(query, studentID).Scan(&totalUnits)
	return totalUnits, err
}

// SaveResult stores a student's result.
// Failed courses are marked as carry-over automatically.
func (s *SQLiteEngine) SaveResult(result *StudentResult) error {
	status := "PASSED"

	if result.Grade == "F" {
		status = "CARRY-OVER"
	}

	query := `
	INSERT INTO results (
		student_id,
		course_code,
		score,
		grade,
		status
	)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		result.StudentID,
		result.CourseCode,
		result.Score,
		result.Grade,
		status,
	)

	return err
}

// GetStudentResults returns all results for a student.
func (s *SQLiteEngine) GetStudentResults(studentID string) ([]StudentResult, error) {
	query := `
	SELECT
		student_id,
		course_code,
		score,
		grade,
		status
	FROM results
	WHERE student_id = ?
	`

	rows, err := s.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []StudentResult

	for rows.Next() {
		var result StudentResult

		err := rows.Scan(
			&result.StudentID,
			&result.CourseCode,
			&result.Score,
			&result.Grade,
			&result.Status,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

// CheckPrerequisiteStatus checks whether
// a student has passed a prerequisite course.
func (s *SQLiteEngine) CheckPrerequisiteStatus(studentID, prereqCode string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM results
		WHERE student_id = ?
		  AND course_code = ?
		  AND status = 'PASSED'
	)
	`

	var passed bool

	err := s.db.QueryRow(query, studentID, prereqCode).Scan(&passed)
	return passed, err
}

