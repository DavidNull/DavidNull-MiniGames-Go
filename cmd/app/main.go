package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	redorblack "davgames/internal/games/RedorBlack"
	"davgames/internal/games/dice"
	"davgames/internal/users"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	usersData, err := users.LoadUsers()
	if err != nil {
		fmt.Printf("\033[1;31mError loading users: %v\033[0m\n", err)
		return
	}

	clearScreen()
	fmt.Print("\033[1;33m🎲 WELCOME TO CASINO GAMES 🎲\033[0m\n")
	fmt.Print("\033[1;33m=============================\033[0m\n")
	usersData.ShowPlayers()

	fmt.Print("\n\033[1;36mEnter your username🚀: \033[0m")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Print("\033[1;36mEnter your password🔑: \033[0m")
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	currentUser, err := usersData.Authenticate(username, password)
	if err != nil {
		fmt.Printf("\033[1;31mAuthentication failed: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\n\033[1;32m✅ Welcome, %s! Your balance 💰 is $%d\033[0m\n", currentUser.Name, currentUser.Balance)
	fmt.Print("\033[1;37mPress Enter to continue...\033[0m\n")
	scanner.Scan()

	for {
		displayMenu()
		choice := getUserChoice(scanner)

		switch choice {
		case "1":
			playDiceGame(currentUser)
		case "2":
			playRedOrBlack(currentUser)
		case "3":
			playSlotMachine()
		case "4":
			fmt.Print("\033[1;33mLogging out...\033[0m\n")
			main()
		case "q", "Q":
			fmt.Print("\033[1;32mThanks for playing! Goodbye.\033[0m\n")
			return
		default:
			fmt.Print("\033[1;31mInvalid choice. Please try again.\033[0m\n")
		}

		fmt.Print("\n\033[1;37mPress Enter to continue...\033[0m\n")
		scanner.Scan()
	}
}

func displayMenu() {
	clearScreen()
	fmt.Print("\033[1;33m🎲 CASINO GAMES 🎲\033[0m\n")
	fmt.Print("\033[1;33m==================\033[0m\n")
	fmt.Print("\033[1;36m1. Dice Game 🎲\033[0m\n")
	fmt.Print("\033[1;36m2. Red or Black 🃏 \033[1;33m(House Favorite! 🌟)\033[0m\n")
	fmt.Print("\033[1;36m3. Slot Machine 🎰\033[0m\n")
	fmt.Print("\033[1;36m4. Logout 🔑\033[0m\n")
	fmt.Print("\033[1;31mQ. Quit 🚫\033[0m\n")
	fmt.Print("\033[1;33m==================\033[0m\n")
	fmt.Print("\033[1;37mEnter your choice: \033[0m")
}

func getUserChoice(scanner *bufio.Scanner) string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func playDiceGame(user *users.User) {
	game := dice.New(user)
	game.Play()
}

func playRedOrBlack(user *users.User) {
	game := redorblack.New(user)
	game.Play()
}

func playSlotMachine() {
	fmt.Print("\n\033[1;35m🎰 SLOT MACHINE 🎰\033[0m\n")
	fmt.Print("\033[1;33mComing soon...\033[0m\n")
}
