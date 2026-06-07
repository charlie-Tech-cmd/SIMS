package main

import (
	"bufio"
	"time"
)

// User represents the standardized identity profile layout for all campus users.
type User struct {
	ID         string // Unique Primary Key: Matriculation Number (Student) or Staff ID (Lecturer)
	Surname    string // Automatically capitalized on registration
	MiddleName string // Optional field, automatically capitalized
	FirstName  string // Automatically capitalized on registration
	Password   string // Plaintext secure key, sanitized against delimiter hijacking
	Email      string // Target address for simulated OTP notification dispatches
	Role       string // Authorization Level: Exactly "student" or "lecturer"
}

// StudentResult holds score evaluations mapped to individuals with CGPA credit weighting.
type StudentResult struct {
	StudentID   string // Links to User.ID
	CourseCode  string // Course catalog tracking identifier (e.g., "CMP301")
	Grade       string // Performance score letter metric achieved (e.g., "A", "F")
	CreditUnits int    // University Standard weight impact indicator (e.g., 3 or 4)
}

// Notice tracks broadcast announcements published exclusively by administrative staff.
type Notice struct {
	ID        string    // Unique string token generated from current Unix epoch time
	Author    string    // Dynamic string displaying the posting professor's full name
	Content   string    // Message text string, sanitized against pipe character conflicts
	Timestamp time.Time // Exact moment of creation, stored in RFC3339 layout standard
}

// TimetableItem holds organizational structures for tests and finals schedules.
type TimetableItem struct {
	ID         string // Custom identification reference code
	Type       string // Categorization token: Exactly "Assessment" or "Examination"
	CourseCode string // Target course identifier
	DateTime   string // Free text timing (e.g., "Monday morning, 10:00 AM")
	Venue      string // Physical campus location descriptor (e.g., "Room 402, Block C")
}

// OTP stores state details for password recovery routines.
type OTP struct {
	Code      string    // Cryptographically generated random 6-digit confirmation key
	ExpiresAt time.Time // Expiry absolute timeline flag (Valid for 5 minutes)
}

// SchoolSystem coordinates all structural global engine parameters and memory map caches.
type SchoolSystem struct {
	// Storage Configuration (File Paths to Flat-File Databases)
	UserFile      string
	ResultFile    string
	NoticeFile    string
	TimetableFile string

	// Live Memory Map Tracking Structures
	Users       map[string]User
	Results     map[string]StudentResult // Compound Key Strategy: "StudentID_CourseCode"
	Notices     map[string]Notice        // Key: Notice.ID string
	Timetables  map[string]TimetableItem // Key: TimetableItem.ID string
	PendingOTPs map[string]OTP           // Key: User.ID string

	// Central Input Management (Prevents OS Standard Input buffer pollution)
	ConsoleScanner *bufio.Scanner
}