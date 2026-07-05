package main

import (
	"fmt"
	"sims/storage"
	"strings"
)

const MAX_SEMESTER_CREDIT_LIMIT = 24

// RunStudentDashboard displays the student dashboard
// and handles student actions.
func (ss *SchoolSystem) RunStudentDashboard(student *storage.User) {

	deptCode := "GENERAL"

	idParts := strings.Split(student.ID, "/")

	if len(idParts) > 1 {
		deptCode = strings.ToUpper(idParts[1])
	}

	for {
		ss.ClearTerminal()

		fmt.Println("==================================================")
		fmt.Printf("  Student Dashboard — %s %s (%s)\n",
			student.Surname,
			student.FirstName,
			student.ID,
		)
		fmt.Println("==================================================")

		// Show latest department notice
		latestNotice, err := ss.DB.GetLatestNoticeForDept(deptCode)

		if err == nil && latestNotice != nil {

			fmt.Printf("\nLatest Notice (%s)\n",
				strings.ToUpper(latestNotice.TargetDepartment),
			)

			fmt.Println("--------------------------------------------------")

			fmt.Printf("From : %s\n", latestNotice.Author)

			fmt.Printf("Date : %s\n",
				latestNotice.Timestamp.Format("Jan 02, 15:04"),
			)

			fmt.Printf("Message:\n%s\n", latestNotice.Content)

			fmt.Println("--------------------------------------------------")
		}

		fmt.Println("\n1. View Results")
		fmt.Println("2. View Notices")
		fmt.Println("3. Register Courses")
		fmt.Println("4. Logout")

		fmt.Print("\nChoose an option: ")

		choice := ss.ReadInput()

		switch choice {

		case "1":
			ss.ClearTerminal()
			ss.HandleViewStudentResults(student)

			fmt.Print("\nPress ENTER to continue...")
			ss.ReadInput()

		case "2":
			ss.ClearTerminal()
			ss.HandleViewBulletinBoard(deptCode)

			fmt.Print("\nPress ENTER to continue...")
			ss.ReadInput()

		case "3":
			ss.ClearTerminal()
			ss.HandleCourseRegistration(student)

			fmt.Print("\nPress ENTER to continue...")
			ss.ReadInput()

		case "4":
			return

		default:
			fmt.Println("Invalid option.")
		}
	}

}

// HandleCourseRegistration manages student course registration.
func (ss *SchoolSystem) HandleCourseRegistration(student *storage.User) {

	ss.DisplaySystemHeader("Course Registration")

	results, err := ss.DB.GetStudentResults(student.ID)

	if err != nil {
		fmt.Println("Could not load student results.")
		return
	}

	var carryOverUnits int
	var carryOverCourses []string

	for _, result := range results {

		if result.Status == "CARRY-OVER" {

			course, _ := ss.DB.GetCourse(result.CourseCode)

			if course != nil {
				carryOverUnits += course.Units
				carryOverCourses = append(carryOverCourses, result.CourseCode)
			}
		}
	}

	registeredCourses, err := ss.DB.GetStudentRegisteredCourses(student.ID)

	if err != nil {
		fmt.Println("Could not load registered courses.")
		return
	}

	registeredUnits := 0

	for _, course := range registeredCourses {
		registeredUnits += course.Units
	}

	totalUnits := registeredUnits + carryOverUnits
	remainingUnits := MAX_SEMESTER_CREDIT_LIMIT - totalUnits

	fmt.Printf("Student ID: %s\n", student.ID)

	fmt.Println("--------------------------------------------------")

	fmt.Println("Registered Courses:")

	if len(registeredCourses) == 0 {

		fmt.Println("No registered courses yet.")

	} else {

		for _, course := range registeredCourses {
			fmt.Printf("- %s\n", course.Code)
		}
	}

	fmt.Println("--------------------------------------------------")

	fmt.Printf("Registered Units : %d\n", registeredUnits)

	fmt.Printf("Carry-Over Units : %d %v\n",
		carryOverUnits,
		carryOverCourses,
	)

	fmt.Printf("Total Load       : %d / %d\n",
		totalUnits,
		MAX_SEMESTER_CREDIT_LIMIT,
	)

	fmt.Printf("Remaining Units  : %d\n", remainingUnits)

	fmt.Println("--------------------------------------------------")

	if remainingUnits <= 0 {

		fmt.Println("\nYou have reached the maximum unit limit.")
		return
	}

	allCourses, err := ss.DB.GetAllCourses()

	if err != nil || len(allCourses) == 0 {
		fmt.Println("No courses available.")
		return
	}

	fmt.Println("\nAvailable Courses")

	fmt.Printf("%-10s %-30s %-8s %-15s\n",
		"CODE",
		"TITLE",
		"UNITS",
		"PREREQUISITE",
	)

	fmt.Println("----------------------------------------------------------------")

	for _, course := range allCourses {

		prerequisite := course.PrerequisiteCode

		if prerequisite == "" {
			prerequisite = "NONE"
		}

		fmt.Printf("%-10s %-30s %-8d %-15s\n",
			course.Code,
			course.Title,
			course.Units,
			prerequisite,
		)
	}

	fmt.Print("\nEnter course code (or press ENTER to cancel): ")

	targetCode := ss.ReadInput()

	if targetCode == "" {
		return
	}

	var selectedCourse *storage.Course

	for _, course := range allCourses {

		if course.Code == targetCode {
			selectedCourse = &course
			break
		}
	}

	if selectedCourse == nil {
		fmt.Println("Course not found.")
		return
	}

	for _, course := range registeredCourses {

		if course.Code == selectedCourse.Code {
			fmt.Println("You already registered this course.")
			return
		}
	}

	if selectedCourse.PrerequisiteCode != "" {

		passed, err := ss.DB.CheckPrerequisiteStatus(
			student.ID,
			selectedCourse.PrerequisiteCode,
		)

		if err != nil || !passed {

			fmt.Printf(
				"You must pass %s before registering %s.\n",
				selectedCourse.PrerequisiteCode,
				selectedCourse.Code,
			)

			return
		}
	}

	if selectedCourse.Units > remainingUnits {

		fmt.Printf(
			"Not enough remaining units for %s.\n",
			selectedCourse.Code,
		)

		return
	}

	err = ss.DB.RegisterStudentForCourse(
		student.ID,
		selectedCourse.Code,
	)

	if err != nil {
		fmt.Println("Could not register course.")
		return
	}

	fmt.Printf(
		"\n%s registered successfully.\n",
		selectedCourse.Code,
	)

}

// HandleViewStudentResults displays student grades.
func (ss *SchoolSystem) HandleViewStudentResults(student *storage.User) {

	fmt.Println("\n--- Student Results ---")

	results, err := ss.DB.GetStudentResults(student.ID)

	if err != nil || len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Printf("%-12s %-8s %-8s %-12s\n",
		"COURSE",
		"SCORE",
		"GRADE",
		"STATUS",
	)

	fmt.Println("--------------------------------------------------")

	for _, result := range results {

		fmt.Printf("%-12s %-8d %-8s %-12s\n",
			result.CourseCode,
			result.Score,
			result.Grade,
			result.Status,
		)
	}

}

// HandleViewBulletinBoard displays department notices.
func (ss *SchoolSystem) HandleViewBulletinBoard(deptCode string) {

	fmt.Printf("\n--- Notices (%s) ---\n", deptCode)

	notices, err := ss.DB.GetAllNoticesForDept(deptCode)

	if err != nil || len(notices) == 0 {
		fmt.Println("No notices available.")
		return
	}

	for _, notice := range notices {

		fmt.Println("--------------------------------------------------")

		fmt.Printf("From : %s\n", notice.Author)

		fmt.Printf("Date : %s\n",
			notice.Timestamp.Format("Jan 02, 15:04"),
		)

		fmt.Printf("Message:\n%s\n", notice.Content)
	}

}
