package battleship

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"davgames/internal/users"
)

const (
	water      = "🌊"
	hit        = "💥"
	miss       = "❌"
	hidden     = "⬜"
	clearLine  = "\033[2K"
	moveCursor = "\033[1A"
)

var ships = []string{"🛶", "🚤", "🛳️", "🛥️", "🚢"}
var shipSizes = map[string]int{ // same like slots game
	"🛶":  2,
	"🚤":  3,
	"⛵":  3,
	"🛥️": 4,
	"🚢":  5,
}

type Player struct {
	Name      string
	Board     [][]string
	EnemyView [][]string
	Ships     map[string]bool
	ShipsLeft int
}

type Game struct {
	User      *users.User
	Player1   Player
	Player2   Player
	BoardSize int
	Scanner   *bufio.Scanner
}

func New(user *users.User) *Game {
	return &Game{
		User:      user,
		BoardSize: 10,
		Scanner:   bufio.NewScanner(os.Stdin),
	}
}

func (g *Game) Play() {
	rand.Seed(time.Now().UnixNano())

	clearScreen()
	fmt.Println("⚓ \033[1;36mBATTLESHIP\033[0m ⚓")
	fmt.Println("\033[1;33m================================\033[0m")

	g.setup()
	g.gameLoop()
}

func (g *Game) setup() {
	fmt.Printf("\033[1;32m🎮 Player 1: %s\033[0m\n", g.User.Name)

	fmt.Print("\033[1;32m🎮 Player 2 name: \033[0m")
	g.Scanner.Scan()
	player2Name := strings.TrimSpace(g.Scanner.Text())
	if player2Name == "" {
		player2Name = "Player 2"
	}

	g.BoardSize = 10

	g.Player1 = Player{
		Name:      g.User.Name,
		Board:     initializeBoard(g.BoardSize),
		EnemyView: initializeBoard(g.BoardSize),
		Ships:     make(map[string]bool),
		ShipsLeft: len(ships),
	}

	g.Player2 = Player{
		Name:      player2Name,
		Board:     initializeBoard(g.BoardSize),
		EnemyView: initializeBoard(g.BoardSize),
		Ships:     make(map[string]bool),
		ShipsLeft: len(ships),
	}

	fmt.Printf("\n\033[1;36m%s, place your ships:\033[0m\n", g.Player1.Name)
	g.placeShips(&g.Player1)

	clearScreen()
	fmt.Printf("\n\033[1;36m%s, place your ships:\033[0m\n", g.Player2.Name)
	g.placeShips(&g.Player2)

	clearScreen()
	fmt.Println("\033[1;33mLet the naval battle begin!\033[0m")
	time.Sleep(2 * time.Second)
}

func (g *Game) gameLoop() {
	currentTurn := 1

	for {
		var currentPlayer, opponent *Player
		if currentTurn%2 == 1 {
			currentPlayer = &g.Player1
			opponent = &g.Player2
		} else {
			currentPlayer = &g.Player2
			opponent = &g.Player1
		}

		clearScreen()
		fmt.Printf("\033[1;36m%s's turn\033[0m\n", currentPlayer.Name)
		fmt.Println("\033[1;33m================================\033[0m")

		fmt.Printf("\n\033[1;32mYour fleet:\033[0m\n")
		printBoard(currentPlayer.Board)

		fmt.Printf("\n\033[1;31mAttacks on %s:\033[0m\n", opponent.Name)
		printBoard(currentPlayer.EnemyView)

		row, col := g.getAttackCoordinates(currentPlayer)

		hit := g.processAttack(currentPlayer, opponent, row, col)

		if opponent.ShipsLeft == 0 {
			clearScreen()
			fmt.Printf("\033[1;33m%s wins!\033[0m\n", currentPlayer.Name)
			fmt.Println("\033[1;32mAll enemy ships have been sunk!\033[0m")

			fmt.Printf("\n\033[1;36m%s's fleet:\033[0m\n", currentPlayer.Name)
			printBoard(currentPlayer.Board)

			fmt.Printf("\n\033[1;36m%s's fleet:\033[0m\n", opponent.Name)
			printBoard(opponent.Board)

			fmt.Print("\n\033[1;37mPress Enter to continue...\033[0m")
			g.Scanner.Scan()
			return
		}

		if hit {
			fmt.Println("\033[1;32mHit! 💥\033[0m")
		} else {
			fmt.Println("\033[1;34mMiss! 🌊\033[0m")
		}

		fmt.Printf("\n\033[1;33mChanging to %s's turn in ", opponent.Name)
		for i := 3; i > 0; i-- {
			fmt.Printf("%d... ", i)
			time.Sleep(1 * time.Second)
		}

		currentTurn++
	}
}

func (g *Game) placeShips(player *Player) {
	for i, shipType := range ships {
		size := shipSizes[shipType]

		for {
			fmt.Printf("\n\033[1;36mPlacing %s (length %d)\033[0m\n", shipType, size)
			fmt.Println("\033[1;33m1. Auto placement\033[0m")
			fmt.Println("\033[1;33m2. Manual placement\033[0m")
			fmt.Print("\033[1;37mChoose an option: \033[0m")

			g.Scanner.Scan()
			option := strings.TrimSpace(g.Scanner.Text())

			if option == "1" {
				if placeShipRandomly(player.Board, shipType, size) {
					player.Ships[shipType] = true
					fmt.Printf("\033[1;32mShip %s placed automatically\033[0m\n", shipType)
					printBoard(player.Board)
					break
				}
				fmt.Println("\033[1;31mCouldn't place ship. Trying again.\033[0m")
			} else if option == "2" {
				printBoard(player.Board)
				fmt.Println("\033[1;36mDirection:\033[0m")
				fmt.Println("\033[1;33m1. Horizontal\033[0m")
				fmt.Println("\033[1;33m2. Vertical\033[0m")
				fmt.Print("\033[1;37mChoose an option: \033[0m")

				g.Scanner.Scan()
				dirOption := strings.TrimSpace(g.Scanner.Text())
				isHorizontal := dirOption != "2"

				fmt.Printf("\033[1;36mEnter starting coordinate (e.g. A3): \033[0m")
				g.Scanner.Scan()
				coordStr := strings.ToUpper(strings.TrimSpace(g.Scanner.Text()))

				if len(coordStr) < 2 {
					fmt.Println("\033[1;31mInvalid coordinate\033[0m")
					continue
				}

				row := int(coordStr[0] - 'A')
				col, err := strconv.Atoi(coordStr[1:])
				col--

				if err != nil || row < 0 || row >= g.BoardSize || col < 0 || col >= g.BoardSize {
					fmt.Println("\033[1;31mCoordinate out of range\033[0m")
					continue
				}

				if isHorizontal {
					if col+size > g.BoardSize {
						fmt.Println("\033[1;31mShip won't fit horizontally\033[0m")
						continue
					}

					canPlace := true
					for j := 0; j < size; j++ {
						if player.Board[row][col+j] != water {
							canPlace = false
							break
						}
					}

					if !canPlace {
						fmt.Println("\033[1;31mThere's a ship in the way\033[0m")
						continue
					}

					for j := 0; j < size; j++ {
						player.Board[row][col+j] = shipType
					}
				} else {
					if row+size > g.BoardSize {
						fmt.Println("\033[1;31mShip won't fit vertically\033[0m")
						continue
					}

					canPlace := true
					for j := 0; j < size; j++ {
						if player.Board[row+j][col] != water {
							canPlace = false
							break
						}
					}

					if !canPlace {
						fmt.Println("\033[1;31mThere's a ship in the way\033[0m")
						continue
					}

					for j := 0; j < size; j++ {
						player.Board[row+j][col] = shipType
					}
				}

				player.Ships[shipType] = true
				fmt.Printf("\033[1;32mShip %s placed manually\033[0m\n", shipType)
				printBoard(player.Board)
				break
			} else {
				fmt.Println("\033[1;31mInvalid option\033[0m")
			}
		}

		if i == len(ships)-1 {
			fmt.Print("\n\033[1;37mPress Enter to continue...\033[0m")
			g.Scanner.Scan()
		}
	}
}

func (g *Game) getAttackCoordinates(player *Player) (int, int) {
	for {
		fmt.Print("\n\033[1;36mEnter attack coordinates (e.g. B5): \033[0m")
		g.Scanner.Scan()
		coordStr := strings.ToUpper(strings.TrimSpace(g.Scanner.Text()))

		if len(coordStr) < 2 {
			fmt.Println("\033[1;31mInvalid coordinate\033[0m")
			continue
		}

		row := int(coordStr[0] - 'A')
		col, err := strconv.Atoi(coordStr[1:])
		col--

		if err != nil || row < 0 || row >= g.BoardSize || col < 0 || col >= g.BoardSize {
			fmt.Println("\033[1;31mCoordinate out of range\033[0m")
			continue
		}

		if player.EnemyView[row][col] == hit || player.EnemyView[row][col] == miss {
			fmt.Println("\033[1;31mYou've already attacked that position\033[0m")
			continue
		}

		return row, col
	}
}

func (g *Game) processAttack(attacker *Player, defender *Player, row, col int) bool {
	targetCell := defender.Board[row][col]

	if targetCell != water && targetCell != hit && targetCell != miss {
		shipType := targetCell
		attacker.EnemyView[row][col] = hit
		defender.Board[row][col] = hit

		shipSunk := true
		for i := 0; i < g.BoardSize; i++ {
			for j := 0; j < g.BoardSize; j++ {
				if defender.Board[i][j] == shipType {
					shipSunk = false
					break
				}
			}
			if !shipSunk {
				break
			}
		}

		if shipSunk {
			fmt.Printf("\033[1;33mYou've sunk a %s!\033[0m\n", shipType)
			defender.ShipsLeft--
		}

		return true
	} else {
		attacker.EnemyView[row][col] = miss
		defender.Board[row][col] = miss
		return false
	}
}

func initializeBoard(size int) [][]string {
	board := make([][]string, size)
	for i := range board {
		board[i] = make([]string, size)
		for j := range board[i] {
			board[i][j] = water
		}
	}
	return board
}

func printBoard(board [][]string) {
	size := len(board)

	fmt.Print("   ")
	for i := 0; i < size; i++ {
		fmt.Printf("  %d ", i+1)
	}
	fmt.Println()

	fmt.Print("   ")
	for i := 0; i < size; i++ {
		fmt.Print("────")
	}
	fmt.Println()

	for i := 0; i < size; i++ {
		fmt.Printf(" %c │", 'A'+i)
		for j := 0; j < size; j++ {
			fmt.Printf(" %s ", board[i][j])
		}
		fmt.Println()
	}
}

func placeShipRandomly(board [][]string, shipType string, size int) bool {
	boardSize := len(board)
	maxAttempts := 100

	for attempt := 0; attempt < maxAttempts; attempt++ {
		isHorizontal := rand.Intn(2) == 0

		var row, col int
		if isHorizontal {
			row = rand.Intn(boardSize)
			col = rand.Intn(boardSize - size + 1)
		} else {
			row = rand.Intn(boardSize - size + 1)
			col = rand.Intn(boardSize)
		}

		canPlace := true
		if isHorizontal {
			for j := 0; j < size; j++ {
				if board[row][col+j] != water {
					canPlace = false
					break
				}
			}
		} else {
			for j := 0; j < size; j++ {
				if board[row+j][col] != water {
					canPlace = false
					break
				}
			}
		}

		if canPlace {
			if isHorizontal {
				for j := 0; j < size; j++ {
					board[row][col+j] = shipType
				}
			} else {
				for j := 0; j < size; j++ {
					board[row+j][col] = shipType
				}
			}
			return true
		}
	}

	return false
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
