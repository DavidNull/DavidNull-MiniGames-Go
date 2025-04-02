package dice

import (
	"fmt"
	"math/rand" //random
	"strconv"
	"strings"
	"time"

	"davgames/internal/persistence"
	"davgames/internal/users"
)

type Game struct {
	User *users.User
}

func New(user *users.User) *Game {
	return &Game{
		User: user,
	}
}

func (g *Game) Play() {
	rand.Seed(time.Now().UnixNano())

	for {
		fmt.Print("\033[1;36m\n🎲 DICE GAME 🎲\033[0m\n")
		fmt.Println("================")
		fmt.Printf("\033[1;33mYour balance: $%d 💰\033[0m\n\n", g.User.Balance)

		if g.User.Balance <= 0 {
			fmt.Print("\033[1;31mYou're out of money! Game over.\033[0m\n")
			fmt.Println("womp womp 🥀")
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			return
		}

		bet := g.getBet()
		if bet == 0 {
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			return
		}

		targetNumber := g.getTargetNumber()
		if targetNumber == 0 {
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			return
		}

		diceValue := rand.Intn(6) + 1
		fmt.Print("\n\033[1;35mRolling the dice... 🎲 🎲 🎲\033[0m\n")
		time.Sleep(1 * time.Second)
		fmt.Printf("\033[1;36mThe dice shows: %d 🎯\033[0m\n", diceValue)

		if diceValue == targetNumber {
			winnings := bet * 5
			g.User.Balance += winnings
			fmt.Printf("\033[1;32mYou win $%d! 🎉 🎊 💫\033[0m\n", winnings)
		} else {
			g.User.Balance -= bet
			fmt.Printf("\033[1;31mYou lose $%d. Better luck next time! 😢\033[0m\n", bet)
		}

		if err := persistence.SaveUserBalance(g.User); err != nil {
			fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
		}

		var input string
		fmt.Print("\n\033[1;33mPress Enter to continue or type 'exit' to return to the main menu... 🔄\033[0m\n")
		fmt.Scanln(&input)
		if strings.ToLower(input) == "exit" {
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			return
		}
	}
}

func (g *Game) getBet() int {
	for {
		fmt.Printf("\033[1;33mEnter your bet (1-%d) or 0 to exit: $\033[0m", g.User.Balance)
		var input string
		fmt.Scanln(&input)

		if input == "0" {
			return 0
		}

		bet, err := strconv.Atoi(input)
		if err != nil || bet < 1 || bet > g.User.Balance {
			fmt.Print("\033[1;31m❌ Invalid bet amount. Please try again. 😕\033[0m\n")
			continue
		}

		return bet
	}
}

func (g *Game) getTargetNumber() int {
	for {
		fmt.Print("\033[1;33mChoose a number to bet on (1-6) or 0 to exit: 🎯\033[0m ")
		var input string
		fmt.Scanln(&input)

		if input == "0" {
			return 0
		}

		number, err := strconv.Atoi(input)
		if err != nil || number < 1 || number > 6 {
			fmt.Print("\033[1;31mInvalid number. Please choose between 1 and 6. ❌\033[0m\n")
			continue
		}

		return number
	}
}
