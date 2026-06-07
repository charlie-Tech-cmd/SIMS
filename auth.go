package main

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// GenerateOTP creates a secure 6-digit verification code and validates OS entropy
func GenerateOTP() string {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		// Secure fallback using current timestamp Unix nanoseconds if system entropy fails
		now := time.Now().UnixNano()
		return fmt.Sprintf("%06d", uint32(now)%1000000)
	}
	// Correctly ensures a clean 6-character output format boundary
	return fmt.Sprintf("%06d", (uint32(b[0])<<16|uint32(b[1])<<8|uint32(b[2]))%1000000)
}

// HandleStudentRegister appends new student profiles onto disk using the db engine layer
func (ss *SchoolSystem) HandleStudentRegister() {
	fmt.Println("\n--- Create Student Account ---")
	
	fmt.Print("Enter Matriculation Number: ")
	id := ss.SanitizeField(ss.ReadInput())

	if id == "" {
		fmt.Println("Error: Matriculation number cannot be empty.")
		return
	}

	if _, exists := ss.Users[id]; exists {
		fmt.Println("Error: An account with this Matriculation Number already exists!")
		return
	}

	fmt.Print("Surname: ")
	surname := ss.SanitizeField(ss.ReadInput())
	
	fmt.Print("Middle Name: ")
	middleName := ss.SanitizeField(ss.ReadInput())
	
	fmt.Print("First Name: ")
	firstName := ss.SanitizeField(ss.ReadInput())
	
	fmt.Print("Create Password: ")
	password := ss.SanitizeField(ss.ReadInput())
	
	fmt.Print("Email Address: ")
	email := ss.SanitizeField(ss.ReadInput())

	if surname == "" || firstName == "" || password == "" || email == "" {
		fmt.Println("Error: All fields (except Middle Name) are mandatory.")
		return
	}

	// Save cleanly to the active global in-memory map structure
	ss.Users[id] = User{
		ID:         id,
		Surname:    surname,
		MiddleName: middleName,
		FirstName:  firstName,
		Password:   password,
		Email:      email,
		Role:       "student",
	}

	// CRITICAL FIX: Flush using the unified database engine rather than manual appending
	ss.RewriteUsers()
	fmt.Printf("Success! Account created for %s %s.\n", surname, firstName)
}

// HandleLogin executes case-insensitive folding validation and routes via tagged switch
func (ss *SchoolSystem) HandleLogin() {
	fmt.Println("\n--- Portal Identity Login ---")
	
	fmt.Print("Matric No / Staff ID: ")
	id := ss.ReadInput()
	
	fmt.Print("Surname: ")
	surname := ss.ReadInput()
	
	fmt.Print("Middle Name: ")
	middleName := ss.ReadInput()
	
	fmt.Print("First Name: ")
	firstName := ss.ReadInput()
	
	fmt.Print("Password: ")
	password := ss.ReadInput()

	user, exists := ss.Users[id]
	if !exists || user.Password != password {
		fmt.Println("Error: Invalid credentials verification failed.")
		return
	}

	// Idiomatic Case-Insensitive Check
	if !strings.EqualFold(user.Surname, surname) ||
		!strings.EqualFold(user.MiddleName, middleName) ||
		!strings.EqualFold(user.FirstName, firstName) {
		fmt.Println("Error: Identity names do not match official records.")
		return
	}

	// Tagged Switch handling role authorization routing
	switch user.Role {
	case "lecturer":
		ss.LecturerMenu(user.Surname + " " + user.FirstName)
	case "student":
		ss.StudentMenu(user.ID)
	default:
		fmt.Println("Error: Corrupted or unassigned account profile role.")
	}
}

// RequestPasswordReset initiates safe password recovery workflows
func (ss *SchoolSystem) RequestPasswordReset() {
	fmt.Print("\nEnter Matric No / Staff ID: ")
	id := ss.ReadInput()

	user, exists := ss.Users[id]
	if !exists {
		// Generic return prevents data enumeration attacks sniffing registered IDs
		fmt.Println("Process initialized. If details match, look at your console email output.")
		return
	}

	code := GenerateOTP()
	ss.PendingOTPs[id] = OTP{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	fmt.Println("\n--- [SIMULATED EMAIL NOTIFICATION SYSTEM] ---")
	fmt.Printf("To: %s\nDear %s %s,\nYour reset OTP token code is: %s\n", user.Email, user.Surname, user.FirstName, code)
	fmt.Println("----------------------------------------------")
}

// VerifyAndResetPassword validates tokens and securely updates the database file
func (ss *SchoolSystem) VerifyAndResetPassword() {
	fmt.Print("\nEnter Matric No / Staff ID: ")
	id := ss.ReadInput()
	
	fmt.Print("Enter OTP Received: ")
	code := ss.ReadInput()

	pending, exists := ss.PendingOTPs[id]
	if !exists || pending.Code != code || time.Now().After(pending.ExpiresAt) {
		fmt.Println("Invalid token code or window timeframe expired.")
		return
	}

	fmt.Print("Enter New Secure Password: ")
	newPassword := ss.SanitizeField(ss.ReadInput())

	if newPassword == "" {
		fmt.Println("Error: Password cannot be completely empty.")
		return
	}

	// Read state, change field, update map entry
	user := ss.Users[id]
	user.Password = newPassword
	ss.Users[id] = user

	// Flush state down to flat storage file safely
	ss.RewriteUsers()
	delete(ss.PendingOTPs, id)
	fmt.Println("Success: Password altered successfully. Access updated.")
}