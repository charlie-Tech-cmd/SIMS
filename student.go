package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Helper function to read a clean menu option line
func readStudentInput(scanner *bufio.Scanner) string {
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func (ss *SchoolSystem) StudentMenu(studentID string) {
	// Unified scanner to prevent console input buffer pollution
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("\n--- Student Portal (Logged in: %s) ---\n", studentID)
		fmt.Println("1. View My Results")
		fmt.Println("2. View Lecturer Updates & Notices")
		fmt.Println("3. View Exam & Assessment Timetables")
		fmt.Println("4. Logout")
		fmt.Print("Choose action: ")

		choice := readStudentInput(scanner)

		// Fail-safe check for terminal dropouts
		if err := scanner.Err(); err != nil {
			fmt.Println("Terminal interface reading error:", err)
			return
		}

		switch choice {
		case "1":
			fmt.Println("\n=================== OFFICIAL ACADEMIC TRANSCRIPT ===================")
			fmt.Printf("Student Matriculation Number: %s\n", studentID)
			fmt.Println("--------------------------------------------------------------------")
			fmt.Printf("%-12s | %-12s | %-12s | %-12s\n", "Course Code", "Credit Units", "Letter Grade", "Grade Points")
			fmt.Println("--------------------------------------------------------------------")

			var totalQualityPoints float64 // Sum of (GP * Units)
			var totalCreditUnits float64   // Sum of Units
			found := false

			for _, res := range ss.Results {
				if res.StudentID == studentID {
					// 1. Convert Letter Grade to standard 5-point University scale
					var gradePoint float64
					switch strings.ToUpper(res.Grade) {
					case "A":
						gradePoint = 5.0
					case "B":
						gradePoint = 4.0
					case "C":
						gradePoint = 3.0
					case "D":
						gradePoint = 2.0
					case "E":
						gradePoint = 1.0
					default: // "F" or invalid entries
						gradePoint = 0.0
					}

					// 2. Display this specific course line item row
					fmt.Printf("%-12s | %-12d | %-12s | %-12.1f\n", 
						res.CourseCode, res.CreditUnits, res.Grade, gradePoint)

					// 3. Accumulate figures for the final CGPA mathematical logic
					totalQualityPoints += (gradePoint * float64(res.CreditUnits))
					totalCreditUnits += float64(res.CreditUnits)
					found = true
				}
			}

			if !found {
				fmt.Println("\nNo academic results or score records published for your profile yet.")
			} else {
				// 4. Calculate the Final GPA score
				var cgpa float64
				if totalCreditUnits > 0 {
					cgpa = totalQualityPoints / totalCreditUnits
				}

				// 5. Print out the formal summary graduation metrics dashboard
				fmt.Println("--------------------------------------------------------------------")
				fmt.Printf("Total Registered Units Attempted: %.1f\n", totalCreditUnits)
				fmt.Printf("Total Earned Quality Points:      %.1f\n", totalQualityPoints)
				
				// Standard University Classification text
				var class string
				switch {
				case cgpa >= 4.50:
					class = "First Class Honours"
				case cgpa >= 3.50:
					class = "Second Class Honours (Upper Division)"
				case cgpa >= 2.40:
					class = "Second Class Honours (Lower Division)"
				case cgpa >= 1.50:
					class = "Third Class Honours"
				default:
					class = "Pass / Academic Probation Status"
				}
				
				fmt.Printf("CURRENT CGPA:                     \033[1;32m%.2f\033[0m (%s)\n", cgpa, class)
				fmt.Println("====================================================================")
			}

		case "2":
			fmt.Println("\n--- LECTURER NOTICES & ANNOUNCEMENTS ---")
			if len(ss.Notices) == 0 {
				fmt.Println("No active announcements or notices posted.")
			} else {
				for _, n := range ss.Notices {
					// Safeguard: omit empty or broken notices
					if strings.TrimSpace(n.Content) == "" {
						continue
					}
					fmt.Printf("[%s] Prof. %s posted:\n   \"%s\"\n\n", 
						n.Timestamp.Format("2006-01-02 15:04"), n.Author, n.Content)
				}
			}

		case "3":
			fmt.Println("\n--- ASSESSMENT & EXAMINATION TIMETABLE ---")
			if len(ss.Timetables) == 0 {
				fmt.Println("No test or examination schedules assigned yet.")
			} else {
				for _, t := range ss.Timetables {
					fmt.Printf("[%s] Course: %s | When: %s | Venue: %s (ID: %s)\n", 
						t.Type, t.CourseCode, t.DateTime, t.Venue, t.ID)
				}
			}

		case "4":
			fmt.Println("Logging out of Student Portal... Return safe.")
			return

		default:
			fmt.Println("Invalid option. Please input a number from 1 to 4.")
		}
	}
}