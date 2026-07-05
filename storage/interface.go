package storage

import "time"

// User account details.
type User struct {
	ID               string
	Surname          string
	FirstName        string
	Email            string
	Password         string
	Role             string // student or lecturer
	SecurityQuestion string
	SecurityAnswer   string
}

// Course information.
type Course struct {
	Code             string
	Title            string
	Units            int
	PrerequisiteCode string
	LecturerID       string
	LecturerName     string
}

// Student academic result.
type StudentResult struct {
	StudentID  string
	CourseCode string
	Score      int
	Grade      string
	Status     string // PASSED or CARRY-OVER
}

// Notice or announcement.
type Notice struct {
	ID               string
	AuthorID         string
	Author           string
	TargetDepartment string
	Content          string
	Timestamp        time.Time
}

// StorageEngine handles all database operations.
type StorageEngine interface {
	Close() error

	// Users
	GetUserByID(id string) (*User, error)
	SaveUser(user *User) error
	UpdateUserPassword(id string, newPassword string) error
	CheckLecturerExists() (bool, error)

	// Courses
	GetCourse(code string) (*Course, error)
	GetAllCourses() ([]Course, error)
	SaveCourse(course *Course) error

	// Results
	SaveResult(result *StudentResult) error
	GetStudentResults(studentID string) ([]StudentResult, error)

	// Notices
	SaveNotice(notice *Notice) error
	GetAllNotices() ([]Notice, error)
	GetAllNoticesForDept(dept string) ([]Notice, error)
	GetLatestNoticeForDept(dept string) (*Notice, error)
	GetLatestNotice() (*Notice, error)

	// Registrations
	GetStudentRegisteredUnits(studentID string) (int, error)
	RegisterStudentForCourse(studentID, courseCode string) error
	GetStudentRegisteredCourses(studentID string) ([]Course, error)
	CheckPrerequisiteStatus(studentID, prereqCode string) (bool, error)
}

