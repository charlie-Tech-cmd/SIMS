package storage

import (
	"database/sql"
	"errors"
)

// SaveNotice stores a new notice in the database.
func (s *SQLiteEngine) SaveNotice(notice *Notice) error {
	query := `
	INSERT INTO notices (
		id,
		author_id,
		author,
		target_department,
		content
	)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		notice.ID,
		notice.AuthorID,
		notice.Author,
		notice.TargetDepartment,
		notice.Content,
	)

	return err
}

// GetAllNotices returns all notices sorted by latest first.
func (s *SQLiteEngine) GetAllNotices() ([]Notice, error) {
	query := `
	SELECT
		id,
		author_id,
		author,
		target_department,
		content,
		created_at
	FROM notices
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notices []Notice

	for rows.Next() {
		var notice Notice

		err := rows.Scan(
			&notice.ID,
			&notice.AuthorID,
			&notice.Author,
			&notice.TargetDepartment,
			&notice.Content,
			&notice.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		notices = append(notices, notice)
	}

	return notices, nil
}

// GetAllNoticesForDept returns notices for a department,
// including general announcements.
func (s *SQLiteEngine) GetAllNoticesForDept(dept string) ([]Notice, error) {
	query := `
	SELECT
		id,
		author_id,
		author,
		target_department,
		content,
		created_at
	FROM notices
	WHERE target_department = ?
	   OR target_department = 'GENERAL'
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, dept)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notices []Notice

	for rows.Next() {
		var notice Notice

		err := rows.Scan(
			&notice.ID,
			&notice.AuthorID,
			&notice.Author,
			&notice.TargetDepartment,
			&notice.Content,
			&notice.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		notices = append(notices, notice)
	}

	return notices, nil
}

// GetLatestNoticeForDept returns the newest notice
// for a department or general audience.
func (s *SQLiteEngine) GetLatestNoticeForDept(dept string) (*Notice, error) {
	query := `
	SELECT
		id,
		author_id,
		author,
		target_department,
		content,
		created_at
	FROM notices
	WHERE target_department = ?
	   OR target_department = 'GENERAL'
	ORDER BY created_at DESC
	LIMIT 1
	`

	var notice Notice

	err := s.db.QueryRow(query, dept).Scan(
		&notice.ID,
		&notice.AuthorID,
		&notice.Author,
		&notice.TargetDepartment,
		&notice.Content,
		&notice.Timestamp,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &notice, nil
}

// GetLatestNotice returns the most recent notice.
func (s *SQLiteEngine) GetLatestNotice() (*Notice, error) {
	query := `
	SELECT
		id,
		author_id,
		author,
		target_department,
		content,
		created_at
	FROM notices
	ORDER BY created_at DESC
	LIMIT 1
	`

	var notice Notice

	err := s.db.QueryRow(query).Scan(
		&notice.ID,
		&notice.AuthorID,
		&notice.Author,
		&notice.TargetDepartment,
		&notice.Content,
		&notice.Timestamp,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &notice, nil
}
