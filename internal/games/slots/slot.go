package slots

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"davgames/internal/persistence"
	"davgames/internal/users"
)

var symbols = []string{"🍒", "🍊", "🍋", "🍇", "🍉", "🍎", "7️⃣", "🐰"}

var symbolWeights = map[string]int{ // + weights + probabilities of appearing
	"🍒":   25,
	"🍊":   20,
	"🍋":   15,
	"🍇":   12,
	"🍉":   10,
	"🍎":   8,
	"7️⃣": 5,
	"🐰":   2,
}

var payouts = map[string]int{
	"🍒🍒🍒":       3,
	"🍊🍊🍊":       4,
	"🍋🍋🍋":       5,
	"🍇🍇🍇":       6,
	"🍉🍉🍉":       8,
	"🍎🍎🍎":       10,
	"7️⃣7️⃣7️⃣": 15,
	"🐰🐰🐰":       25, //signature symbol ;D
}

type Game struct {
	User    *users.User
	Scanner *bufio.Scanner
}

func New(user *users.User) *Game {
	return &Game{
		User:    user,
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

func (g *Game) Play() {
	rand.Seed(time.Now().UnixNano())

	for {
		fmt.Println("\n\033[1;35m🎰 SLOT MACHINE 🎰\033[0m")
		fmt.Println("\033[1;33m==================\033[0m")
		fmt.Printf("\033[1;32mYour balance: $%d\033[0m\n\n", g.User.Balance)

		fmt.Print("\033[1;33mType 'help' to see the paytable or press Enter to continue: \033[0m")
		g.Scanner.Scan()
		if strings.ToLower(g.Scanner.Text()) == "help" {
			g.displayPaytable()
			fmt.Print("\nPress Enter to continue...")
			g.Scanner.Scan()
		}

		if g.User.Balance <= 0 {
			fmt.Println("\033[1;31mYou're out of money! Game over.\033[0m")
			fmt.Println("\033[1;31mwomp womp 🥀\033[0m")
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

		result := g.spinReels()
		g.displaySpinAnimation(result)

		combination := result[0] + result[1] + result[2]
		multiplier, isWinning := payouts[combination]

		if isWinning {
			winnings := bet * multiplier
			g.User.Balance += winnings
			fmt.Printf("\033[1;32m🎉 You won $%d! (Bet $%d × %dx multiplier)\033[0m\n", winnings, bet, multiplier)
		} else if result[0] == result[1] || result[1] == result[2] || result[0] == result[2] {
			fmt.Printf("\033[1;33m🎯 Two matching symbols! You get your $%d bet back.\033[0m\n", bet)
		} else {
			g.User.Balance -= bet
			fmt.Printf("\033[1;31m❌ Sorry, you lost $%d.\033[0m\n", bet)
		}

		fmt.Printf("\033[1;32mNew balance: $%d\033[0m\n", g.User.Balance)

		if err := persistence.SaveUserBalance(g.User); err != nil {
			fmt.Printf("\033[1;31mError saving user data: %v\033[0m\n", err)
		}

		fmt.Println("\n\033[1;37mPress Enter to play again or type 'exit' to return to the main menu...\033[0m")
		g.Scanner.Scan()
		if strings.ToLower(g.Scanner.Text()) == "exit" {
			return
		}
	}
}

func (g *Game) getBet() int {
	for {
		fmt.Printf("\033[1;36mEnter your bet (1-%d) or 0 to exit: $\033[0m", g.User.Balance)
		g.Scanner.Scan()
		input := g.Scanner.Text()

		if input == "0" {
			return 0
		}

		bet, err := strconv.Atoi(input)
		if err != nil || bet < 1 || bet > g.User.Balance {
			fmt.Println("\033[1;31mInvalid bet amount. Please try again.\033[0m")
			continue
		}

		return bet
	}
}

func getWeightedSymbol() string {
	totalWeight := 0
	for _, weight := range symbolWeights {
		totalWeight += weight
	}

	r := rand.Intn(totalWeight)

	currentWeight := 0
	for _, symbol := range symbols {
		currentWeight += symbolWeights[symbol]
		if r < currentWeight {
			return symbol
		}
	}

	return symbols[0]
}

func (g *Game) spinReels() []string {
	reels := make([]string, 3)
	for i := 0; i < 3; i++ {
		reels[i] = getWeightedSymbol()
	}
	return reels
}

func (g *Game) displaySpinAnimation(finalResult []string) {
	fmt.Println("\n\033[1;36mSpinning the reels...\033[0m")

	frames := 10

	for frame := 0; frame < frames; frame++ {
		if frame > 0 {
			fmt.Print("\r")
		}

		if frame < frames-1 {
			tempReels := make([]string, 3)
			for i := 0; i < 3; i++ {
				tempReels[i] = getWeightedSymbol()
			}
			fmt.Printf("\033[1;35m[ %s | %s | %s ]\033[0m", tempReels[0], tempReels[1], tempReels[2])
		} else {
			fmt.Printf("\033[1;33m[ %s | %s | %s ]\033[0m", finalResult[0], finalResult[1], finalResult[2])
		}

		sleepTime := 100 + (frames-frame)*50
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}

	fmt.Println()
}

func (g *Game) displayPaytable() {
	fmt.Println("\n\033[1;33m💰 PAYTABLE 💰\033[0m")
	fmt.Println("\033[1;33m========================\033[0m")
	fmt.Println("\033[1;36m🍒🍒🍒\033[0m = 3× your bet")
	fmt.Println("\033[1;36m🍊🍊🍊\033[0m = 4× your bet")
	fmt.Println("\033[1;36m🍋🍋🍋\033[0m = 5× your bet")
	fmt.Println("\033[1;36m🍇🍇🍇\033[0m = 6× your bet")
	fmt.Println("\033[1;36m🍉🍉🍉\033[0m = 8× your bet")
	fmt.Println("\033[1;36m🍎🍎🍎\033[0m = 10× your bet")
	fmt.Println("\033[1;36m7️⃣7️⃣7️⃣\033[0m = 15× your bet")
	fmt.Println("\033[1;36m🐰🐰🐰\033[0m = 25× your bet")
	fmt.Println("\033[1;33mAny two matching symbols = Get your bet back\033[0m")
	fmt.Println()
	fmt.Println("\033[1;34mPROBABILITIES:\033[0m")
	fmt.Println(" 🍒 > 🍊 > 🍋 > 🍇 > 🍉 > 🍎 > 7️⃣ > 🐰 ")
	fmt.Println("\033[1;33m========================\033[0m")
}

//probabilities:
// 🍒: 25%
// 🍊: 20%
// 🍋: 15%
// 🍇: 12%
// 🍉: 10%
// 🍎: 8%
// 7️⃣: 5%
// 🐰: 2%

// centerText centra el texto en un ancho específico
func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}

	spaces := width - len(text)
	leftSpaces := spaces / 2
	rightSpaces := spaces - leftSpaces

	return strings.Repeat(" ", leftSpaces) + text + strings.Repeat(" ", rightSpaces)
}
