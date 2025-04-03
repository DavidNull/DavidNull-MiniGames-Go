package logging

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type User struct {
	Username string
	Password string
	Balance  int
}

type AuthManager struct {
	Users   []User
	Scanner *bufio.Scanner
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		Users: []User{
			{Username: "player1", Password: "pass1", Balance: 200},
			{Username: "player2", Password: "pass2", Balance: 300},
			{Username: "player3", Password: "pass3", Balance: 500},
			{Username: "admin", Password: "admin123", Balance: 1000},
		},
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

func (am *AuthManager) Authenticate() (*User, bool) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("🎲 WELCOME TO CASINO GAMES 🎲")
	fmt.Println("=============================")

	fmt.Println("\nAvailable Players:")
	for i, user := range am.Users {
		fmt.Printf("%d. %s\n", i+1, user.Username)
	}

	var username string
	for {
		fmt.Print("\nEnter your username (or 'q' to quit): ")
		am.Scanner.Scan()
		username = strings.TrimSpace(am.Scanner.Text())

		if username == "q" {
			return nil, false
		}

		userExists := false
		for _, user := range am.Users {
			if user.Username == username {
				userExists = true
				break
			}
		}

		if userExists {
			break
		} else {
			fmt.Println("⚠️ User not found. Please try again.")
		}
	}

	var password string
	attemptsLeft := 3

	for attemptsLeft > 0 {
		fmt.Print("Enter your password: ")
		am.Scanner.Scan()
		password = strings.TrimSpace(am.Scanner.Text())

		for i, user := range am.Users {
			if user.Username == username && user.Password == password {
				fmt.Printf("\n✅ Welcome, %s! Your balance is $%d\n", user.Username, user.Balance)
				fmt.Println("Press Enter to continue...")
				am.Scanner.Scan()
				return &am.Users[i], true
			}
		}

		attemptsLeft--
		if attemptsLeft > 0 {
			fmt.Printf("❌ Incorrect password. %d attempts left.\n", attemptsLeft)
		} else {
			fmt.Println("❌ Too many failed attempts. Access denied.")
			return nil, false
		}
	}

	return nil, false
}
