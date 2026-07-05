package storage

import (
	"database/sql"
	"errors"
)

// SaveUser stores a new user account.
func (s *SQLiteEngine) SaveUser(user *User) error {
	query := `
	INSERT INTO users (
		id,
		surname,
		first_name,
		email,
		password,
		role,
		security_question,
		security_answer
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		user.ID,
		user.Surname,
		user.FirstName,
		user.Email,
		user.Password,
		user.Role,
		user.SecurityQuestion,
		user.SecurityAnswer,
	)

	return err
}

// GetUserByID returns a user by ID.
func (s *SQLiteEngine) GetUserByID(id string) (*User, error) {
	query := `
	SELECT
		id,
		surname,
		first_name,
		email,
		password,
		role,
		security_question,
		security_answer
	FROM users
	WHERE id = ?
	`

	row := s.db.QueryRow(query, id)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Surname,
		&user.FirstName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.SecurityQuestion,
		&user.SecurityAnswer,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserPassword changes a user's password.
func (s *SQLiteEngine) UpdateUserPassword(id string, newPassword string) error {
	query := `
	UPDATE users
	SET password = ?
	WHERE id = ?
	`

	_, err := s.db.Exec(query, newPassword, id)
	return err
}

// CheckLecturerExists checks whether
// at least one lecturer account exists.
func (s *SQLiteEngine) CheckLecturerExists() (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE role = 'lecturer'
		LIMIT 1
	)
	`

	var exists bool

	err := s.db.QueryRow(query).Scan(&exists)
	return exists, err
}
