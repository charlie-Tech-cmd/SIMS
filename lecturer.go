package main

import (
	"fmt"
	"sims/storage"
	"strconv"
)

// RunLecturerDashboard displays lecturer options
// and handles lecturer actions.
func (ss *SchoolSystem) RunLecturerDashboard(lecturer *storage.User) {
	for {
		ss.ClearTerminal()

		fmt.Println("==================================================")
		fmt.Printf("  Lecturer Dashboard — %s %s\n", lecturer.Surname, lecturer.FirstName)
		fmt.Println("==================================================")

		fmt.Println("\n1. View Course Enrollments")
		fmt.Println("2. Upload Student Result")
		fmt.Println("3. Post Notice")
		fmt.Println("4. Logout")

		fmt.Print("\nChoose an option: ")

		choice := ss.ReadInput()

		switch choice {

		case "1":
			ss.ClearTerminal()
			ss.HandleViewCourseEnrollments(lecturer)

			fmt.Print("\nPress ENTER to continue...")
			ss.ReadInput()

		case "2":
			ss.ClearTerminal()
			ss.HandleGradeSubmission(lecturer)

			fmt.Print("\nPress ENTER to continue...")
			ss.ReadInput()

		case "3":
			ss.ClearTerminal()
			ss.HandlePublishNotice(lecturer)

			fmt.Print("\nPress ENTER to continue...")
			ss.ReadInput()

		case "4":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}

}

// HandleViewCourseEnrollments shows courses
// assigned to the lecturer.
func (ss *SchoolSystem) HandleViewCourseEnrollments(lecturer *storage.User) {
	fmt.Println("\n--- Course Enrollments ---")

	allCourses, err := ss.DB.GetAllCourses()

	if err != nil || len(allCourses) == 0 {
		fmt.Println("No courses found.")
		return
	}

	found := false

	for _, course := range allCourses {

		if course.LecturerID == lecturer.ID {
			found = true

			fmt.Printf("\nCourse: %s - %s\n", course.Code, course.Title)
			fmt.Println("--------------------------------------------")
			fmt.Println("Students are currently enrolled in this course.")
		}
	}

	if !found {
		fmt.Println("No courses assigned to you yet.")
	}

}

// HandleGradeSubmission allows lecturers
// to upload student grades.
func (ss *SchoolSystem) HandleGradeSubmission(lecturer *storage.User) {
	fmt.Println("\n--- Upload Student Result ---")

	fmt.Print("Enter Student ID: ")
	studentID := ss.ReadInput()

	if studentID == "" {
		return
	}

	fmt.Print("Enter Course Code: ")
	courseCode := ss.ReadInput()

	if courseCode == "" {
		return
	}

	course, err := ss.DB.GetCourse(courseCode)

	if err != nil || course == nil {
		fmt.Println("Course not found.")
		return
	}

	if course.LecturerID != lecturer.ID {
		fmt.Println("You are not assigned to this course.")
		return
	}

	fmt.Print("Enter Score (0 - 100): ")

	scoreInput := ss.ReadInput()

	score, err := strconv.Atoi(scoreInput)

	if err != nil || score < 0 || score > 100 {
		fmt.Println("Enter a valid score between 0 and 100.")
		return
	}

	var grade string

	switch {
	case score >= 70:
		grade = "A"

	case score >= 60:
		grade = "B"

	case score >= 50:
		grade = "C"

	case score >= 45:
		grade = "D"

	default:
		grade = "F"
	}

	result := &storage.StudentResult{
		StudentID:  studentID,
		CourseCode: courseCode,
		Score:      score,
		Grade:      grade,
	}

	if err := ss.DB.SaveResult(result); err != nil {
		fmt.Println("Could not save result.")
		return
	}

	fmt.Printf("\nResult uploaded successfully.\n")
	fmt.Printf("Student: %s\n", studentID)
	fmt.Printf("Course : %s\n", courseCode)
	fmt.Printf("Grade  : %s\n", grade)

}

// HandlePublishNotice allows lecturers
// to publish notices to students.
func (ss *SchoolSystem) HandlePublishNotice(lecturer *storage.User) {
	fmt.Println("\n--- Post Notice ---")

	fmt.Print("Enter Notice ID: ")
	noticeID := ss.ReadInput()

	if noticeID == "" {
		return
	}

	fmt.Print("Enter Department (or GENERAL): ")
	targetDept := ss.ReadInput()

	if targetDept == "" {
		targetDept = "GENERAL"
	}

	fmt.Print("Enter Notice Message:\n> ")

	content := ss.ReadInput()

	if content == "" {
		return
	}

	notice := &storage.Notice{
		ID:               noticeID,
		AuthorID:         lecturer.ID,
		Author:           lecturer.Surname + " " + lecturer.FirstName,
		TargetDepartment: targetDept,
		Content:          content,
	}

	if err := ss.DB.SaveNotice(notice); err != nil {
		fmt.Println("Could not publish notice.")
		return
	}

	fmt.Println("\nNotice posted successfully.")

}
