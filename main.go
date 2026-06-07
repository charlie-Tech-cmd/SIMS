package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// 1. Initialize the system and immediately assign a shared global scanner
	app := &SchoolSystem{
		UserFile:      "users.txt",
		ResultFile:    "results.txt",
		NoticeFile:    "notices.txt",
		TimetableFile: "timetables.txt",
		ConsoleScanner: bufio.NewScanner(os.Stdin), // Keeps input clean across the entire app
	}

	// 2. Synchronize memory state with disk files
	if err := app.LoadAllData(); err != nil {
		log.Fatal("Fatal Initialization Error: Could not synchronize database files: ", err)
	}

	for {
		fmt.Println("\n=== CAMPUS CENTRAL PORTAL ===")
		fmt.Println("1. Secure Identity Portal Login")
		fmt.Println("2. Register New Student Account")
		fmt.Println("3. Forgotten Password Request (Generate OTP)")
		fmt.Println("4. Execute Account Password Change (Verify OTP)")
		fmt.Println("5. Disconnect System")
		fmt.Print("Select option: ")
		
		// 3. Read input cleanly using our global structural scanner
		var choice string
		if app.ConsoleScanner.Scan() {
			choice = strings.TrimSpace(app.ConsoleScanner.Text())
		}

		// Check if standard input experienced an OS-level terminal disconnect
		if err := app.ConsoleScanner.Err(); err != nil {
			log.Println("Terminal stream error encountered:", err)
			return
		}

		switch choice {
		case "1":
			app.HandleLogin()
		case "2":
			app.HandleStudentRegister() 
		case "3":
			app.RequestPasswordReset()
		case "4":
			app.VerifyAndResetPassword()
		case "5":
			fmt.Println("System Offline. Safe shutdown completed.")
			return
		default:
			fmt.Println("Invalid input. Please choose a valid numeric command option (1-5).")
		}
	}
}