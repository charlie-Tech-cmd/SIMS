package main

import (
	"fmt"
	"strings"

	"sims/storage"
)

// EnsureSeedDataExist adds default users and courses
// if the database is empty.
func (ss *SchoolSystem) EnsureSeedDataExist() {
	exists, err := ss.DB.CheckLecturerExists()
	if err != nil || exists {
		return
	}

	fmt.Println("Setting up default system data...")

	seedUsers := []storage.User{
		{
			ID:               "STF001",
			Surname:          "Eze",
			FirstName:        "Chidi",
			Password:         "password123",
			Email:            "c.eze@university.edu",
			Role:             "lecturer",
			SecurityQuestion: "What is your favorite color?",
			SecurityAnswer:   "blue",
		},
		{
			ID:               "STF002",
			Surname:          "Okonkwo",
			FirstName:        "Bisi",
			Password:         "password123",
			Email:            "b.okonkwo@university.edu",
			Role:             "lecturer",
			SecurityQuestion: "What town were you born in?",
			SecurityAnswer:   "lagos",
		},
		{
			ID:               "STU001",
			Surname:          "Adeleke",
			FirstName:        "Tunde",
			Password:         "password123",
			Email:            "t.adeleke@university.edu",
			Role:             "student",
			SecurityQuestion: "What is your pet's name?",
			SecurityAnswer:   "rex",
		},
	}

	for _, user := range seedUsers {
		_ = ss.DB.SaveUser(&user)
	}

	seedCourses := []storage.Course{
		{
			Code:             "CMP101",
			Title:            "Introduction to Computer Science",
			Units:            3,
			PrerequisiteCode: "",
			LecturerID:       "STF001",
			LecturerName:     "Eze Chidi",
		},
		{
			Code:             "MTH101",
			Title:            "General Mathematics I",
			Units:            4,
			PrerequisiteCode: "",
			LecturerID:       "STF002",
			LecturerName:     "Okonkwo Bisi",
		},
		{
			Code:             "CMP201",
			Title:            "Data Structures and Algorithms",
			Units:            4,
			PrerequisiteCode: "CMP101",
			LecturerID:       "STF001",
			LecturerName:     "Eze Chidi",
		},
		{
			Code:             "CMP301",
			Title:            "Advanced Software Architecture",
			Units:            6,
			PrerequisiteCode: "CMP201",
			LecturerID:       "STF001",
			LecturerName:     "Eze Chidi",
		},
	}

	for _, course := range seedCourses {
		_ = ss.DB.SaveCourse(&course)
	}

	mockResult := &storage.StudentResult{
		StudentID:  "STU001",
		CourseCode: "MTH101",
		Score:      35,
		Grade:      "F",
		Status:     "CARRY-OVER",
	}

	_ = ss.DB.SaveResult(mockResult)

	fmt.Println("Default users and courses added successfully.")
}

// HandleUserLogin authenticates a user account.
func (ss *SchoolSystem) HandleUserLogin() {
	fmt.Println("\n--- Login ---")

	fmt.Print("Enter ID: ")
	id := ss.ReadInput()

	if id == "" {
		return
	}

	fmt.Print("Enter Password: ")
	password := ss.ReadInput()

	user, err := ss.DB.GetUserByID(id)
	if err != nil {
		fmt.Printf("Could not complete login: %v\n", err)
		return
	}

	if user == nil || user.Password != password {
		fmt.Println("Invalid ID or password.")
		return
	}

	fmt.Printf("\nWelcome back, %s %s.\n", user.Surname, user.FirstName)

	if user.Role == "lecturer" {
		ss.RunLecturerDashboard(user)
	} else {
		ss.RunStudentDashboard(user)
	}
}

// HandleDirectPasswordRecovery helps users reset forgotten passwords.
func (ss *SchoolSystem) HandleDirectPasswordRecovery() {
	fmt.Println("\n--- Password Recovery ---")

	fmt.Print("Enter your User ID: ")
	id := ss.ReadInput()

	if id == "" {
		return
	}

	user, err := ss.DB.GetUserByID(id)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return
	}

	if user == nil {
		fmt.Println("No account found with that ID.")
		return
	}

	if user.SecurityQuestion == "" || user.SecurityAnswer == "" {
		fmt.Println("This account does not have recovery settings enabled.")
		return
	}

	fmt.Printf("\nSecurity Question for %s:\n", user.ID)
	fmt.Printf("%s\n", user.SecurityQuestion)

	fmt.Print("Answer: ")
	answer := strings.ToLower(ss.ReadInput())

	if answer != strings.ToLower(user.SecurityAnswer) {
		fmt.Println("Incorrect answer.")
		return
	}

	fmt.Print("\nEnter your new password: ")
	newPassword := ss.ReadInput()

	if len(newPassword) < 6 {
		fmt.Println("Password must be at least 6 characters long.")
		return
	}

	if err := ss.DB.UpdateUserPassword(user.ID, newPassword); err != nil {
		fmt.Printf("Could not update password: %v\n", err)
		return
	}

	fmt.Println("Password updated successfully.")
}

// HandleNewUserRegistration creates a new account.
func (ss *SchoolSystem) HandleNewUserRegistration() {
	fmt.Println("\n--- Create Account ---")

	fmt.Print("Choose Role (1 = Student, 2 = Lecturer): ")
	roleChoice := ss.ReadInput()

	var role string

	if roleChoice == "1" {
		role = "student"
	} else if roleChoice == "2" {
		role = "lecturer"
	} else {
		fmt.Println("Invalid role selection.")
		return
	}

	fmt.Print("Enter Unique ID: ")
	id := ss.ReadInput()

	if existingUser, _ := ss.DB.GetUserByID(id); existingUser != nil {
		fmt.Println("This ID is already registered.")
		return
	}

	fmt.Print("Enter Surname: ")
	surname := ss.ReadInput()

	fmt.Print("Enter First Name: ")
	firstName := ss.ReadInput()

	fmt.Print("Enter Official Email: ")
	email := strings.ToLower(ss.ReadInput())

	if !strings.HasSuffix(email, "@university.edu") {
		fmt.Println("Please use your official school email.")
		return
	}

	fmt.Print("Set Password: ")
	password := ss.ReadInput()

	fmt.Print("Set Recovery Question: ")
	securityQuestion := ss.ReadInput()

	fmt.Print("Set Recovery Answer: ")
	securityAnswer := strings.ToLower(ss.ReadInput())

	newUser := &storage.User{
		ID:               id,
		Surname:          surname,
		FirstName:        firstName,
		Password:         password,
		Email:            email,
		Role:             role,
		SecurityQuestion: securityQuestion,
		SecurityAnswer:   securityAnswer,
	}

	if err := ss.DB.SaveUser(newUser); err != nil {
		fmt.Printf("Could not create account: %v\n", err)
		return
	}

	fmt.Println("Account created successfully.")
}
