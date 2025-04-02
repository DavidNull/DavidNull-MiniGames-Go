package redorblack

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"davgames/internal/persistence"
	"davgames/internal/users"
)

type Game struct {
	User            *users.User
	consecutiveWins int
	lastCard        int
}

func New(user *users.User) *Game {
	return &Game{
		User:            user,
		consecutiveWins: 0,
	}
}

func (g *Game) Play() {
	rand.Seed(time.Now().UnixNano())
	for {
		fmt.Print("\033[1;36m\n🃏 RED OR BLACK 🃏\033[0m\n")
		fmt.Println("==================")
		fmt.Printf("\033[1;33mYour balance: $%d 💰\033[0m\n\n", g.User.Balance)

		if g.User.Balance <= 0 {
			fmt.Print("\033[1;31mYou're out of money! Game over.\033[0m\n")
			fmt.Println("womp womp 🥀")
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			return
		}

		fmt.Print("\033[1;33mType 'help' to see the rules or press Enter to continue: \033[0m")
		var helpInput string
		fmt.Scanln(&helpInput)
		if strings.ToLower(helpInput) == "help" {
			fmt.Print("\n\033[1;32m=== GAME RULES ===\033[0m\n")
			fmt.Print("\033[1;37m1. Multiple betting rounds:\033[0m\n")
			fmt.Print("\033[1;37m   - First: Choose red or black\033[0m\n")
			fmt.Print("\033[1;37m   - Second: Higher or lower than previous card\033[0m\n")
			fmt.Print("\033[1;37m   - Third: Between or outside previous numbers\033[0m\n")
			fmt.Print("\033[1;37m   - Final: Guess the suit\033[0m\n")
			fmt.Print("\033[1;33m2. Consecutive wins increase multiplier:\033[0m\n")
			fmt.Print("\033[1;33m   - x2, x3, and x20 on final round\033[0m\n")
			fmt.Print("\033[1;32m================\033[0m\n\n")
			fmt.Print("Press Enter to continue...")
			fmt.Scanln()
		}

		bet := g.getBet()
		if bet == 0 {
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			return
		}

		choice := g.getChoice()
		if choice == "" {
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			return
		}

		result := []string{"red", "black"}[rand.Intn(2)]
		g.lastCard = rand.Intn(13) + 1 // 1-13 for card values
		fmt.Print("\n\033[1;35mDealing the card... 🃏\033[0m\n")
		time.Sleep(1 * time.Second)

		if result == "red" {
			fmt.Printf("\033[1;31mIt's RED! ♥️ ♦️ Card value: %d\033[0m\n", g.lastCard)
		} else {
			fmt.Printf("\033[1;30mIt's BLACK! ♠️ ♣️ Card value: %d\033[0m\n", g.lastCard)
		}

		if result != choice {
			g.consecutiveWins = 0
			g.User.Balance -= bet
			fmt.Printf("\033[1;31mYou lose $%d. Game over! 😢\033[0m\n", bet)
			continue
		}

		winnings := bet * 2
		fmt.Printf("\033[1;32mYou won! You can take $%d (x2) now or continue for higher multipliers\033[0m\n", winnings)
		fmt.Print("\033[1;33mDo you want to continue? (y/n): \033[0m")
		var continueChoice string
		fmt.Scanln(&continueChoice)
		if strings.ToLower(continueChoice) != "y" {
			g.User.Balance += winnings
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			continue
		}

		fmt.Print("\n\033[1;33mHigher or Lower than the previous card? (h/l): \033[0m")
		var hlChoice string
		fmt.Scanln(&hlChoice)
		newCard := rand.Intn(13) + 1
		fmt.Println("Dealing the card...")
		time.Sleep(1 * time.Second)
		fmt.Printf("\033[1;35mNew card: %d\033[0m\n", newCard)

		correctHL := (hlChoice == "h" && newCard > g.lastCard) || (hlChoice == "l" && newCard < g.lastCard)
		if !correctHL {
			g.consecutiveWins = 0
			g.User.Balance -= bet
			fmt.Printf("\033[1;31mWrong! You lose $%d 😢\033[0m\n", bet)
			continue
		}

		winnings = bet * 3
		fmt.Printf("\033[1;32mYou won! You can take $%d (x3) now or continue for higher multipliers.\033[0m\n", winnings)
		fmt.Print("\033[1;33mDo you want to continue? (y/n): \033[0m")
		fmt.Scanln(&continueChoice)
		if strings.ToLower(continueChoice) != "y" {
			g.User.Balance += winnings
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			continue
		}

		fmt.Print("\n\033[1;33mBetween or Outside the previous numbers? (b/o): \033[0m")
		var boChoice string
		fmt.Scanln(&boChoice)
		finalCard := rand.Intn(13) + 1
		fmt.Println("Dealing the card...")
		time.Sleep(1 * time.Second)
		fmt.Printf("\033[1;35mFinal card: %d\033[0m\n", finalCard)

		min := g.lastCard
		max := newCard
		if g.lastCard > newCard {
			min = newCard
			max = g.lastCard
		}

		isBetween := finalCard > min && finalCard < max
		correctBO := (boChoice == "b" && isBetween) || (boChoice == "o" && !isBetween)

		if !correctBO {
			g.consecutiveWins = 0
			g.User.Balance -= bet
			fmt.Printf("\033[1;31mWrong! You lose $%d 😢\033[0m\n", bet)
			continue
		}

		winnings = bet * 5
		fmt.Printf("\033[1;32mYou won! You can take $%d (x5) now or risk it all for x20 on the final round!\033[0m\n", winnings)
		fmt.Print("\033[1;33mDo you want to continue? (y/n): \033[0m")
		fmt.Scanln(&continueChoice)
		if strings.ToLower(continueChoice) != "y" {
			g.User.Balance += winnings
			if err := persistence.SaveUserBalance(g.User); err != nil {
				fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
			}
			continue
		}

		fmt.Print("\n\033[1;33mGuess the suit (hearts/diamonds/clubs/spades): \033[0m")
		var suitChoice string
		fmt.Scanln(&suitChoice)
		fmt.Println("Dealing the card...")
		time.Sleep(1 * time.Second)
		suits := []string{"hearts", "diamonds", "clubs", "spades"}
		finalSuit := suits[rand.Intn(4)]
		fmt.Printf("\033[1;35mFinal suit: %s\033[0m\n", finalSuit)

		if suitChoice != finalSuit {
			g.consecutiveWins = 0
			g.User.Balance -= bet
			fmt.Printf("\033[1;31mWrong! You lose $%d 😢\033[0m\n", bet)
			continue
		}

		g.consecutiveWins++
		winnings = bet * 20
		g.User.Balance += winnings
		fmt.Printf("\033[1;32mCongratulations! You win $%d! (x20 multiplier!) 🎉\033[0m\n", winnings)

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

func (g *Game) getChoice() string {
	for {
		fmt.Print("\033[1;33mChoose 'red' or 'black' (or 'exit' to quit): 🎯\033[0m ")
		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(input)

		if input == "exit" {
			return ""
		}

		if input != "red" && input != "black" {
			fmt.Print("\033[1;31mInvalid choice. Please choose 'red' or 'black'. ❌\033[0m\n")
			continue
		}

		return input
	}
}

// long as hell
