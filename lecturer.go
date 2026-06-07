package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LecturerMenu handles administrative dashboard views attached to the SchoolSystem struct context.
func (ss *SchoolSystem) LecturerMenu(lecturerName string) {
	for {
		fmt.Printf("\n--- Lecturer Dashboard (Welcome, Prof. %s) ---\n", lecturerName)
		fmt.Println("1. Post/Update Student Result")
		fmt.Println("2. Post Availability Notice")
		fmt.Println("3. Post/Update Assessment or Exam Timetable")
		fmt.Println("4. Register a New Faculty Lecturer")
		fmt.Println("5. Logout")
		fmt.Print("Choose action: ")

		choice := ss.ReadInput()

		switch choice {
		case "1":
			for {
				fmt.Print("\nEnter Student Matriculation Number: ")
				sID := ss.SanitizeField(ss.ReadInput())

				if _, exists := ss.Users[sID]; !exists {
					fmt.Println(" Error: Student profile ID does not exist in the database.")
					break
				}

				fmt.Print("Enter Course Code (e.g., CMP301): ")
				course := strings.ToUpper(ss.SanitizeField(ss.ReadInput()))

				fmt.Print("Enter Grade/Score Letter (A, B, C, D, E, F): ")
				grade := strings.ToUpper(ss.SanitizeField(ss.ReadInput()))

				if grade != "A" && grade != "B" && grade != "C" && grade != "D" && grade != "E" && grade != "F" {
					fmt.Println(" Error: Invalid grade letter. Must be A, B, C, D, E, or F.")
					continue
				}

				fmt.Print("Enter Course Credit Units (e.g., 3 or 4): ")
				unitsStr := ss.SanitizeField(ss.ReadInput())
				
				units, err := strconv.Atoi(unitsStr)
				if err != nil || units <= 0 || units > 6 {
					fmt.Println(" Error: Credit units must be a number between 1 and 6.")
					continue
				}

				if sID == "" || course == "" {
					fmt.Println(" Error: Missing fields. All parameters are required.")
					continue
				}

				key := sID + "_" + course
				ss.Results[key] = StudentResult{
					StudentID:   sID,
					CourseCode:  course,
					Grade:       grade,
					CreditUnits: units,
				}

				ss.RewriteResults()
				fmt.Println(" Success: Student score updated safely on disk!")
				break
			}

		case "2":
			fmt.Print("\nEnter announcement message: ")
			msg := ss.SanitizeField(ss.ReadInput())

			if msg == "" {
				fmt.Println(" Error: Cannot post an empty announcement.")
				continue
			}

			id := fmt.Sprintf("%d", time.Now().UnixNano())
			ss.Notices[id] = Notice{
				ID:        id,
				Author:    lecturerName,
				Content:   msg,
				Timestamp: time.Now(),
			}

			ss.RewriteNotices()
			fmt.Println(" Notice posted successfully!")

		case "3":
			fmt.Print("\nEnter Schedule ID: ")
			id := ss.SanitizeField(ss.ReadInput())

			fmt.Print("Type (Assessment/Exam): ")
			tType := ss.SanitizeField(ss.ReadInput())

			fmt.Print("Course Code: ")
			course := strings.ToUpper(ss.SanitizeField(ss.ReadInput()))

			fmt.Print("Date & Time: ")
			dt := ss.SanitizeField(ss.ReadInput())

			fmt.Print("Venue: ")
			venue := ss.SanitizeField(ss.ReadInput())

			if id == "" || tType == "" || course == "" || dt == "" || venue == "" {
				fmt.Println(" Error: Missing parameters. Operation aborted.")
				continue
			}

			ss.Timetables[id] = TimetableItem{
				ID:         id,
				Type:       tType,
				CourseCode: course,
				DateTime:   dt,
				Venue:      venue,
			}

			ss.RewriteTimetables()
			fmt.Println(" Timetable updated successfully!")

		case "4":
			fmt.Println("\n--- Create New Lecturer Account ---")
			
			fmt.Print("Enter New Staff ID / Username (e.g., STF002): ")
			staffID := ss.SanitizeField(ss.ReadInput())
			
			if _, exists := ss.Users[staffID]; exists {
				fmt.Println(" Error: A user with this ID already exists in the system.")
				continue
			}

			fmt.Print("Enter Surname: ")
			surname := ss.SanitizeField(ss.ReadInput())
			
			fmt.Print("Enter Middle Name (Optional, press Enter to skip): ")
			middleName := ss.SanitizeField(ss.ReadInput())
			
			fmt.Print("Enter First Name: ")
			firstName := ss.SanitizeField(ss.ReadInput())
			
			fmt.Print("Set Temporary Account Password: ")
			password := ss.SanitizeField(ss.ReadInput())
			
			fmt.Print("Enter Official Email Address: ")
			email := ss.SanitizeField(ss.ReadInput())

			if staffID == "" || surname == "" || firstName == "" || password == "" || email == "" {
				fmt.Println(" Error: All fields except Middle Name are strictly mandatory.")
				continue
			}

			ss.Users[staffID] = User{
				ID:         staffID,
				Surname:    surname,
				MiddleName: middleName,
				FirstName:  firstName,
				Password:   password,
				Email:      email,
				Role:       "lecturer",
			}

			ss.RewriteUsers()
			fmt.Printf(" Success! Account for Prof. %s has been created.\n", surname)

		case "5":
			fmt.Println("Logging out of administrative terminal...")
			return

		default:
			fmt.Println("Invalid choice, please select an option from 1 to 5.")
		}
	}
}