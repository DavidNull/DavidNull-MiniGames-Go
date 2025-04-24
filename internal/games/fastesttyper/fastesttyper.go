package fastesttyper

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"davgames/internal/users"
)

var wordLists = map[string][]string{
	"animals":    {"dog", "cat", "lion", "tiger", "elephant", "zebra", "giraffe", "monkey", "bear", "wolf", "rabbit", "snake", "panda", "koala", "penguin"},
	"colors":     {"red", "blue", "green", "yellow", "purple", "orange", "black", "white", "pink", "brown", "gray", "gold", "silver", "bronze", "maroon"},
	"fruits":     {"apple", "banana", "orange", "grape", "strawberry", "pear", "peach", "mango", "kiwi", "lemon", "lime", "cherry", "plum", "fig", "papaya"},
	"countries":  {"france", "spain", "japan", "brazil", "canada", "india", "china", "russia", "italy", "mexico", "egypt", "korea", "peru", "chile", "cuba"},
	"friends":    {"joel", "jorge", "nestor", "adrian", "lucas", "tomas", "hector", "carlos", "jaime"},
	"sports":     {"soccer", "tennis", "golf", "rugby", "boxing", "skiing", "surfing", "running", "cycling", "swimming", "diving", "skating", "hockey", "judo", "polo"},
	"simple":     {"the", "and", "is", "in", "it", "you", "that", "he", "was", "for", "on", "are", "with", "as", "his"},
	"verbs":      {"run", "jump", "swim", "eat", "sleep", "read", "write", "sing", "dance", "play", "work", "talk", "walk", "think", "laugh"},
	"adjectives": {"big", "small", "fast", "slow", "hot", "cold", "new", "old", "good", "bad", "high", "low", "rich", "poor", "wise"},
}

type Game struct {
	User1          *users.User
	Player2        string
	Scanner        *bufio.Scanner
	WordList       []string
	ChallengeWords []string
}

func New(user1 *users.User) *Game {
	return &Game{
		User1:   user1,
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

func (g *Game) Play() {
	fmt.Println("\n\033[1;36m🤠 FASTEST TYPER IN THE WEST ⌨️\033[0m")
	fmt.Println("\033[1;33m===============================\033[0m")

	fmt.Print("\n\033[1;33mEnter second player's name: \033[0m")
	g.Scanner.Scan()
	g.Player2 = strings.TrimSpace(g.Scanner.Text())
	if g.Player2 == "" {
		g.Player2 = "Player 2"
	}

	fmt.Println("\n\033[1;33mSelect a word list for the duel:\033[0m")
	fmt.Println("\033[1;36m1. Animals\033[0m")
	fmt.Println("\033[1;36m2. Colors\033[0m")
	fmt.Println("\033[1;36m3. Fruits\033[0m")
	fmt.Println("\033[1;36m4. Countries\033[0m")
	fmt.Println("\033[1;36m5. Food\033[0m")
	fmt.Println("\033[1;36m6. Sports\033[0m")
	fmt.Println("\033[1;36m7. Simple words\033[0m")
	fmt.Println("\033[1;36m8. Verbs\033[0m")
	fmt.Println("\033[1;36m9. Adjectives\033[0m")
	fmt.Println("\033[1;36m10. Custom words\033[0m")
	fmt.Print("\033[1;37mEnter your choice (1-10): \033[0m")

	g.Scanner.Scan()
	listChoice := strings.TrimSpace(g.Scanner.Text())

	switch listChoice {
	case "1":
		g.WordList = wordLists["animals"]
	case "2":
		g.WordList = wordLists["colors"]
	case "3":
		g.WordList = wordLists["fruits"]
	case "4":
		g.WordList = wordLists["countries"]
	case "5":
		g.WordList = wordLists["food"]
	case "6":
		g.WordList = wordLists["sports"]
	case "7":
		g.WordList = wordLists["simple"]
	case "8":
		g.WordList = wordLists["friends"]
	case "9":
		g.WordList = wordLists["adjectives"]
	case "10":
		fmt.Print("\n\033[1;33mEnter custom words separated by spaces: \033[0m")
		g.Scanner.Scan()
		customWords := strings.TrimSpace(g.Scanner.Text())
		g.WordList = strings.Fields(customWords)
	default:
		g.WordList = wordLists["simple"]
	}

	g.selectChallengeWords()

	g.displayRules()

	fmt.Print("\n\033[1;36mPress Enter when both players are ready...\033[0m")
	g.Scanner.Scan()

	fmt.Printf("\n\033[1;32m%s's turn! Get ready...\033[0m\n", g.User1.Name)
	time.Sleep(2 * time.Second)
	player1Time, player1Correct, player1Text := g.runTypingRound()

	clearScreen()
	g.selectChallengeWords()
	g.displayRules()

	fmt.Printf("\n\033[1;32m%s's turn! Get ready...\033[0m\n", g.Player2)
	time.Sleep(2 * time.Second)
	player2Time, player2Correct, player2Text := g.runTypingRound()

	clearScreen()
	fmt.Println("\n\033[1;33m🔫 DUEL RESULTS 🔫\033[0m")

	fmt.Printf("\033[1;32m%s typed: \"%s\"\033[0m\n", g.User1.Name, player1Text)
	if player1Correct {
		fmt.Printf("\033[1;32mAll words correct! Time: %.2f seconds\033[0m\n", player1Time.Seconds())
	} else {
		fmt.Printf("\033[1;31mIncorrect words. Time: %.2f seconds\033[0m\n", player1Time.Seconds())
	}

	fmt.Printf("\033[1;32m%s typed: \"%s\"\033[0m\n", g.Player2, player2Text)
	if player2Correct {
		fmt.Printf("\033[1;32mAll words correct! Time: %.2f seconds\033[0m\n", player2Time.Seconds())
	} else {
		fmt.Printf("\033[1;31mIncorrect words. Time: %.2f seconds\033[0m\n", player2Time.Seconds())
	}

	if player1Correct && !player2Correct {
		fmt.Printf("\033[1;33mThe winner is: %s!\033[0m\n", g.User1.Name)
	} else if !player1Correct && player2Correct {
		fmt.Printf("\033[1;33mThe winner is: %s!\033[0m\n", g.Player2)
	} else if player1Correct && player2Correct {
		if player1Time < player2Time {
			fmt.Printf("\033[1;33mThe winner is: %s with the fastest time!\033[0m\n", g.User1.Name)
		} else if player2Time < player1Time {
			fmt.Printf("\033[1;33mThe winner is: %s with the fastest time!\033[0m\n", g.Player2)
		} else {
			fmt.Println("\033[1;33mIt's a tie! Both players were equally fast!\033[0m")
		}
	} else {
		fmt.Println("\033[1;33mNo winner - both players had incorrect words.\033[0m")
	}

	fmt.Print("\n\033[1;37mPress Enter to continue...\033[0m")
	g.Scanner.Scan()
}

func (g *Game) selectChallengeWords() {
	rand.Seed(time.Now().UnixNano())

	wordsCopy := make([]string, len(g.WordList))
	copy(wordsCopy, g.WordList)

	rand.Shuffle(len(wordsCopy), func(i, j int) {
		wordsCopy[i], wordsCopy[j] = wordsCopy[j], wordsCopy[i]
	})

	numWords := 5
	if len(wordsCopy) < 5 {
		numWords = len(wordsCopy)
	}

	g.ChallengeWords = wordsCopy[:numWords]
}

func (g *Game) displayRules() {
	fmt.Println("\n\033[1;32mWords to type:\033[0m")
	fmt.Printf("\033[1;36m%s\033[0m\n", strings.Join(g.ChallengeWords, " "))

	fmt.Println("\n\033[1;33mRules of the duel:\033[0m")
	fmt.Println("1. Type the exact words shown above, separated by spaces")
	fmt.Println("2. Press Enter when you finish typing")
	fmt.Println("3. You will be timed - the fastest correct typist wins!")
	fmt.Println("4. All words must be typed correctly to qualify")
}

func (g *Game) runTypingRound() (time.Duration, bool, string) {
	fmt.Println("\033[1;33m3...\033[0m")
	time.Sleep(1 * time.Second)
	fmt.Println("\033[1;33m2...\033[0m")
	time.Sleep(1 * time.Second)
	fmt.Println("\033[1;33m1...\033[0m")
	time.Sleep(1 * time.Second)
	fmt.Println("\033[1;32mGO! Type the words and press Enter when done:\033[0m")

	startTime := time.Now()

	reader := bufio.NewReader(os.Stdin)
	typedText, _ := reader.ReadString('\n')
	typedText = strings.TrimSpace(typedText)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	typedWords := strings.Fields(typedText)
	isCorrect := compareWords(typedWords, g.ChallengeWords)

	return duration, isCorrect, typedText
}

func compareWords(typed, challenge []string) bool {
	if len(typed) != len(challenge) {
		return false
	}

	for i := range typed {
		if typed[i] != challenge[i] {
			return false
		}
	}

	return true
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
