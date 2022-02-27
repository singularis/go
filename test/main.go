package main

import (
	"booking-app/helper"
	"fmt"
	"time"
)

var conferenceName = "Go test"

const conferenceTickets uint = 50

var remainingTickets uint = 50

var bookings = make([]UserData, 0)

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOftickets uint
}

func main() {

	gretUsers()

	for remainingTickets > 0 && len(bookings) < 50 {
		var firstName string
		var lastName string
		var email string
		var userTickets uint
		fmt.Println("Enter your first name")
		fmt.Scan(&firstName)
		fmt.Println("Enter your last name")
		fmt.Scan(&lastName)
		fmt.Println("Enter your email")
		fmt.Scan(&email)
		fmt.Println("Enter number of tickets")
		fmt.Scan(&userTickets)

		isValidName, isValideEmail, isValideTickers := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

		if isValidName && isValideEmail && isValideTickers {

			bookTickets(userTickets, firstName, lastName, email)
			go sendTicket(userTickets, firstName, lastName, email)

			fmt.Printf("%v tickets remaining from %v tickets \n", remainingTickets, conferenceTickets)

			fmt.Printf("First name of bookings is  %v \n", getFirstNames())

			noTicketsRemainig := remainingTickets == 0

			if noTicketsRemainig {
				fmt.Println("All tickets are sold")
				break
			}
		} else {
			if !isValidName {
				fmt.Println("Incorrect name")
			}
			if !isValideEmail {
				fmt.Println("Incorrect emale")
			}
			if !isValideTickers {
				fmt.Println("Incorrect amount of tickers")
			}
		}
	}
}

func gretUsers() {
	fmt.Printf("Welcome %v test application \n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v still avaliable \n", conferenceTickets, remainingTickets)
	fmt.Println("I am testind booking application")
}

func getFirstNames() []string {
	firstNames := []string{}
	for _, name := range bookings {
		firstNames = append(firstNames, name.firstName)
	}
	return firstNames
}

func bookTickets(userTickets uint, firstName string, lastName string, email string) uint {
	remainingTickets = remainingTickets - userTickets

	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOftickets: userTickets,
	}
	bookings = append(bookings, userData)
	fmt.Printf("Thank you %v %v for booking %v tickets. You will recieved it by mail %v \n", firstName, lastName, userTickets, email)
	return remainingTickets
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(10 * time.Second)
	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("###########")
	fmt.Printf("Sending ticket: \n %v \n for email address %v \n", ticket, email)
	fmt.Println("###########")
}
