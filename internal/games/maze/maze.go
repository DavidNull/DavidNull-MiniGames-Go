// https://github.com/naveensajeendran/PythonEscapeTheMaze?tab=readme-ov-file
// DFS algorithm
package maze

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"davgames/internal/users"
)

const (
	// Maze constants
	mazeWidth  = 21 // Must be odd
	mazeHeight = 11 // Must be odd

	// Cell types
	wall   = '#'
	path   = ' '
	player = 'P'
	exit   = 'E'

	// Movement keys
	moveUp    = 'w'
	moveDown  = 's'
	moveLeft  = 'a'
	moveRight = 'd'
)

// Game represents the maze game
type Game struct {
	User1       *users.User
	Player2     string
	Scanner     *bufio.Scanner
	MazePlayer1 [][]rune
	MazePlayer2 [][]rune
	PlayerPosP1 [2]int // [y, x]
	ExitPosP1   [2]int // [y, x]
	PlayerPosP2 [2]int // [y, x]
	ExitPosP2   [2]int // [y, x]
}

// Position represents a 2D position in the maze
type Position struct {
	Y, X int
}

// Direction represents possible movement directions
type Direction struct {
	Y, X int
}

// Possible movement directions
var directions = []Direction{
	{-2, 0}, // Up
	{2, 0},  // Down
	{0, -2}, // Left
	{0, 2},  // Right
}

// New creates a new maze game instance
func New(user1 *users.User) *Game {
	return &Game{
		User1:   user1,
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

// Play starts the maze game
func (g *Game) Play() {
	fmt.Println("\n\033[1;36m🧭 ESCAPE THE MAZE 🧗‍♂️\033[0m")
	fmt.Println("\033[1;33m===============================\033[0m")

	// Get second player's name
	fmt.Print("\n\033[1;33mEnter second player's name: \033[0m")
	g.Scanner.Scan()
	g.Player2 = strings.TrimSpace(g.Scanner.Text())
	if g.Player2 == "" {
		g.Player2 = "Player 2"
	}

	// Generate mazes for both players
	g.generateMazeForPlayer1()
	g.generateMazeForPlayer2()

	g.displayRules()

	// Wait for players to be ready
	fmt.Print("\n\033[1;36mPress Enter when both players are ready...\033[0m")
	g.Scanner.Scan()

	// Start the game for player 1
	fmt.Printf("\n\033[1;32m%s's turn! Get ready...\033[0m\n", g.User1.Name)
	time.Sleep(2 * time.Second)
	player1Time, player1Success := g.runMazeGame(g.MazePlayer1, &g.PlayerPosP1, g.ExitPosP1)

	// Clear screen and start the game for player 2
	clearScreen()
	g.displayRules()

	fmt.Printf("\n\033[1;32m%s's turn! Get ready...\033[0m\n", g.Player2)
	time.Sleep(2 * time.Second)
	player2Time, player2Success := g.runMazeGame(g.MazePlayer2, &g.PlayerPosP2, g.ExitPosP2)

	// Determine the winner
	clearScreen()
	fmt.Println("\n\033[1;33m🧭 MAZE RACE RESULTS 🧗‍♂️\033[0m")

	fmt.Printf("\033[1;32m%s's result:\033[0m\n", g.User1.Name)
	if player1Success {
		fmt.Printf("\033[1;32mEscaped the maze! Time: %.2f seconds\033[0m\n", player1Time.Seconds())
	} else {
		fmt.Printf("\033[1;31mFailed to escape the maze\033[0m\n")
	}

	fmt.Printf("\033[1;32m%s's result:\033[0m\n", g.Player2)
	if player2Success {
		fmt.Printf("\033[1;32mEscaped the maze! Time: %.2f seconds\033[0m\n", player2Time.Seconds())
	} else {
		fmt.Printf("\033[1;31mFailed to escape the maze\033[0m\n")
	}

	// Determine winner based on success and time
	if player1Success && !player2Success {
		fmt.Printf("\033[1;33mThe winner is: %s!\033[0m\n", g.User1.Name)
	} else if !player1Success && player2Success {
		fmt.Printf("\033[1;33mThe winner is: %s!\033[0m\n", g.Player2)
	} else if player1Success && player2Success {
		if player1Time < player2Time {
			fmt.Printf("\033[1;33mThe winner is: %s with the fastest escape time!\033[0m\n", g.User1.Name)
		} else if player2Time < player1Time {
			fmt.Printf("\033[1;33mThe winner is: %s with the fastest escape time!\033[0m\n", g.Player2)
		} else {
			fmt.Println("\033[1;33mIt's a tie! Both players escaped in the same time!\033[0m")
		}
	} else {
		fmt.Println("\033[1;33mNo winner - both players failed to escape the maze.\033[0m")
	}

	fmt.Print("\n\033[1;37mPress Enter to continue...\033[0m")
	g.Scanner.Scan()
}

// generateMazeForPlayer1 creates a random maze for player 1
func (g *Game) generateMazeForPlayer1() {
	g.MazePlayer1 = make([][]rune, mazeHeight)
	for i := range g.MazePlayer1 {
		g.MazePlayer1[i] = make([]rune, mazeWidth)
		for j := range g.MazePlayer1[i] {
			g.MazePlayer1[i][j] = wall
		}
	}

	rand.Seed(time.Now().UnixNano())
	g.generateMaze(g.MazePlayer1)
	g.placePlayerAndExit(g.MazePlayer1, &g.PlayerPosP1, &g.ExitPosP1)
}

// generateMazeForPlayer2 creates a random maze for player 2
func (g *Game) generateMazeForPlayer2() {
	g.MazePlayer2 = make([][]rune, mazeHeight)
	for i := range g.MazePlayer2 {
		g.MazePlayer2[i] = make([]rune, mazeWidth)
		for j := range g.MazePlayer2[i] {
			g.MazePlayer2[i][j] = wall
		}
	}

	rand.Seed(time.Now().UnixNano() + 1000) // Use a different seed
	g.generateMaze(g.MazePlayer2)
	g.placePlayerAndExit(g.MazePlayer2, &g.PlayerPosP2, &g.ExitPosP2)
}

// generateMaze uses a depth-first search with backtracking to generate a maze
func (g *Game) generateMaze(maze [][]rune) {
	// Start at a random odd position
	startY := 1
	startX := 1
	maze[startY][startX] = path

	// Use a stack for backtracking
	stack := []Position{{startY, startX}}

	// Continue until stack is empty
	for len(stack) > 0 {
		// Get current position
		current := stack[len(stack)-1]

		// Check if there are any valid neighbors
		validDirs := []Direction{}

		for _, dir := range directions {
			// Calculate neighbor position
			newY, newX := current.Y+dir.Y, current.X+dir.X

			// Check if the new position is within bounds and is a wall
			if newY > 0 && newY < mazeHeight-1 && newX > 0 && newX < mazeWidth-1 && maze[newY][newX] == wall {
				// Check if we can tunnel through (we need to check if the cell beyond is a wall)
				// This prevents creating loops
				validDirs = append(validDirs, dir)
			}
		}

		if len(validDirs) > 0 {
			// Choose a random direction
			dir := validDirs[rand.Intn(len(validDirs))]

			// Create a path in that direction
			newY, newX := current.Y+dir.Y, current.X+dir.X
			maze[newY][newX] = path

			// Also create a path one step in the direction (carve through the wall)
			midY, midX := current.Y+dir.Y/2, current.X+dir.X/2
			maze[midY][midX] = path

			// Push the new position to the stack
			stack = append(stack, Position{newY, newX})
		} else {
			// No valid neighbors, backtrack
			stack = stack[:len(stack)-1]
		}
	}
}

// placePlayerAndExit places the player and exit in the maze
func (g *Game) placePlayerAndExit(maze [][]rune, playerPos *[2]int, exitPos *[2]int) {
	// Find all possible positions (path cells)
	var possiblePositions []Position
	for y := 1; y < mazeHeight-1; y++ {
		for x := 1; x < mazeWidth-1; x++ {
			if maze[y][x] == path {
				possiblePositions = append(possiblePositions, Position{y, x})
			}
		}
	}

	// Randomly select start and end positions
	if len(possiblePositions) > 1 {
		startIdx := rand.Intn(len(possiblePositions))
		playerPos[0] = possiblePositions[startIdx].Y
		playerPos[1] = possiblePositions[startIdx].X

		// Make sure exit is far from player
		var farthestPos Position
		maxDist := 0
		for _, pos := range possiblePositions {
			// Calculate Manhattan distance
			dist := abs(pos.Y-playerPos[0]) + abs(pos.X-playerPos[1])
			if dist > maxDist {
				maxDist = dist
				farthestPos = pos
			}
		}

		exitPos[0] = farthestPos.Y
		exitPos[1] = farthestPos.X
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// displayRules displays the game rules
func (g *Game) displayRules() {
	fmt.Println("\n\033[1;33mRules of the maze:\033[0m")
	fmt.Println("1. Each player must escape their own randomly generated maze")
	fmt.Println("2. Use WASD keys to move: W (up), A (left), S (down), D (right)")
	fmt.Println("3. Find the exit (E) to escape the maze")
	fmt.Println("4. The player who escapes in the shortest time wins!")
	fmt.Println("5. Type 'quit' to give up")
}

// runMazeGame runs a single maze game and returns the time taken
func (g *Game) runMazeGame(maze [][]rune, playerPos *[2]int, exitPos [2]int) (time.Duration, bool) {
	fmt.Println("\033[1;33m3...\033[0m")
	time.Sleep(1 * time.Second)
	fmt.Println("\033[1;33m2...\033[0m")
	time.Sleep(1 * time.Second)
	fmt.Println("\033[1;33m1...\033[0m")
	time.Sleep(1 * time.Second)
	fmt.Println("\033[1;32mGO! Navigate the maze to the exit (E):\033[0m")

	// Start the timer
	startTime := time.Now()

	// Set terminal to raw mode to read single keystrokes
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// Disable echo
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// Restore terminal when done
	defer func() {
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	}()

	// Make a copy of the maze to display
	displayMaze := make([][]rune, len(maze))
	for i := range maze {
		displayMaze[i] = make([]rune, len(maze[i]))
		copy(displayMaze[i], maze[i])
	}

	// Set player and exit positions
	displayMaze[playerPos[0]][playerPos[1]] = player
	displayMaze[exitPos[0]][exitPos[1]] = exit

	// Keep track of quit sequence
	var quitSequence []rune
	quitWord := "quit"

	// Game loop
	for {
		// Draw maze
		clearScreen()
		drawMaze(displayMaze)
		fmt.Println("\n\033[1;36mUse W (up), A (left), S (down), D (right) to move. Type 'quit' to give up.\033[0m")

		// Read a single character
		var b = make([]byte, 1)
		os.Stdin.Read(b)
		input := rune(b[0])

		// Check for quit sequence
		quitSequence = append(quitSequence, input)
		if len(quitSequence) > 4 {
			quitSequence = quitSequence[1:] // Keep last 4 characters
		}

		// Check if last 4 characters spell "quit"
		if len(quitSequence) == 4 {
			isQuit := true
			for i, c := range quitWord {
				if i >= len(quitSequence) || quitSequence[i] != c {
					isQuit = false
					break
				}
			}
			if isQuit {
				return time.Since(startTime), false
			}
		}

		// Process move
		newY, newX := playerPos[0], playerPos[1]

		switch input {
		case moveUp:
			newY--
		case moveDown:
			newY++
		case moveLeft:
			newX--
		case moveRight:
			newX++
		default:
			continue // Invalid key, try again
		}

		// Check if new position is valid
		if newY >= 0 && newY < mazeHeight && newX >= 0 && newX < mazeWidth && maze[newY][newX] == path {
			// Update player position
			displayMaze[playerPos[0]][playerPos[1]] = path
			playerPos[0], playerPos[1] = newY, newX
			displayMaze[newY][newX] = player

			// Check if player reached the exit
			if newY == exitPos[0] && newX == exitPos[1] {
				// Draw final state
				clearScreen()
				displayMaze[newY][newX] = player // Show player at exit
				drawMaze(displayMaze)
				fmt.Println("\n\033[1;32mCongratulations! You escaped the maze!\033[0m")
				time.Sleep(2 * time.Second)
				return time.Since(startTime), true
			}
		}
	}
}

// drawMaze draws the maze in the console
func drawMaze(maze [][]rune) {
	for _, row := range maze {
		fmt.Print("\033[1;33m") // Yellow color for maze
		for _, cell := range row {
			switch cell {
			case wall:
				fmt.Print("██")
			case path:
				fmt.Print("  ")
			case player:
				fmt.Print("\033[1;32mP \033[1;33m") // Green player
			case exit:
				fmt.Print("\033[1;31mE \033[1;33m") // Red exit
			}
		}
		fmt.Println("\033[0m") // Reset color
	}
}

// clearScreen clears the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
