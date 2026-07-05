package main

import "fmt"

// RunMainMenuLoop displays the application's main menu
// and keeps the program running until the user exits.
func (ss *SchoolSystem) RunMainMenuLoop() {
	for {
		ss.DisplaySystemHeader("University Portal")

		fmt.Println("1. Login")
		fmt.Println("2. Create Account")
		fmt.Println("3. Recover Password")
		fmt.Println("4. Exit")

		fmt.Print("Choose an option: ")

		choice := ss.ReadInput()

		switch choice {
		case "1":
			ss.HandleUserLogin()

		case "2":
			ss.HandleNewUserRegistration()

		case "3":
			ss.HandleDirectPasswordRecovery()

		case "4":
			fmt.Println("\nThanks for using the School Portal.")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}

}
