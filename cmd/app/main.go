package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	redorblack "davgames/internal/games/RedorBlack"
	"davgames/internal/games/battleship"
	"davgames/internal/games/dice"
	"davgames/internal/games/fastesttyper"
	"davgames/internal/games/maze"
	"davgames/internal/games/slots"
	"davgames/internal/network"
	"davgames/internal/users"
)

const (
	DEFAULT_PORT = "8080"
)

func main() {
	// Definir flags para modos servidor y cliente
	isServer := flag.Bool("server", false, "Run as server")
	connectTo := flag.String("connect", "", "Connect to server IP")
	port := flag.String("port", DEFAULT_PORT, "Port to use for server/client")
	flag.Parse()

	// Ejecutar en modo servidor si se especifica
	if *isServer {
		startServer(*port)
		return
	}

	// Ejecutar en modo cliente si se especifica
	if *connectTo != "" {
		connectToServer(*connectTo, *port)
		return
	}

	// Modo normal de un solo jugador
	startLocalGame()
}

func startServer(port string) {
	fmt.Printf("\033[1;33m🎲 DAVIDNULL GAMES - SERVER MODE 🎲\033[0m\n")
	fmt.Printf("\033[1;33m================================\033[0m\n")
	fmt.Printf("\033[1;32mStarting server on port %s...\033[0m\n", port)

	// Cargar usuarios
	usersData, err := users.LoadUsers()
	if err != nil {
		fmt.Printf("\033[1;31mError loading users: %v\033[0m\n", err)
		return
	}

	// Iniciar servidor HTTP
	http.HandleFunc("/", network.HandleRoot)
	http.HandleFunc("/users", network.HandleUsers(usersData))
	http.HandleFunc("/games", network.HandleGames)

	// Mostrar información del servidor
	localIP := getLocalIP()
	fmt.Printf("\033[1;32mServer running!\033[0m\n")
	fmt.Printf("\033[1;32mLocal IP: %s\033[0m\n", localIP)
	fmt.Printf("\033[1;32mPort: %s\033[0m\n", port)
	fmt.Printf("\033[1;33mPlayers can connect using:\033[0m\n")
	fmt.Printf("\033[1;36m./DavidNullGames --connect %s\033[0m\n", localIP)

	// Iniciar servidor
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectToServer(serverIP, port string) {
	fmt.Printf("\033[1;33m🎲 DAVIDNULL GAMES - CLIENT MODE 🎲\033[0m\n")
	fmt.Printf("\033[1;33m================================\033[0m\n")
	fmt.Printf("\033[1;32mConnecting to server at %s:%s...\033[0m\n", serverIP, port)

	// Intentar conectar al servidor
	_, err := http.Get(fmt.Sprintf("http://%s:%s/", serverIP, port))
	if err != nil {
		fmt.Printf("\033[1;31mError connecting to server: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\033[1;32mConnected to server!\033[0m\n")

	// Iniciar cliente
	network.StartClient(serverIP, port)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

func startLocalGame() {
	scanner := bufio.NewScanner(os.Stdin)

	usersData, err := users.LoadUsers()
	if err != nil {
		fmt.Printf("\033[1;31mError loading users: %v\033[0m\n", err)
		return
	}

	clearScreen()
	fmt.Print("\033[1;33m🎲 WELCOME TO DAVIDNULL GAMES 🎲\033[0m\n")
	fmt.Print("\033[1;33m================================\033[0m\n")
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
			playSlotMachine(currentUser)
		case "4":
			playFastestTyper(currentUser)
		case "5":
			playMazeGame(currentUser)
		case "6":
			playBattleship(currentUser)
		case "7":
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
	fmt.Print("\033[1;33m🎲 DAVIDNULL GAMES 🎲\033[0m\n")
	fmt.Print("\033[1;33m==================\033[0m\n")
	fmt.Print("\033[1;32m🍀 Luck Games 🎲 🎯\033[0m\n")
	fmt.Print("\033[1;36m1. Dice Game 🎲\033[0m\n")
	fmt.Print("\033[1;36m2. Red or Black 🃏 \033[1;33m(House Favorite! 🌟)\033[0m\n")
	fmt.Print("\033[1;36m3. Slot Machine 🎰\033[0m\n")
	fmt.Print("\033[1;32m🎮 Local 2 Players 🎮\033[0m\n")
	fmt.Print("\033[1;36m4. 🤠 Fastest typer in the West ⌨️\033[0m\n")
	fmt.Print("\033[1;36m5. 🧭 Leave the maze! 🧗‍♂️\033[0m\n")
	fmt.Print("\033[1;36m6. 🚢 Battleship 🚢\033[0m\n")
	fmt.Print("\033[1;35m🌐 LAN Games (Coming Soon!) 🌐\033[0m\n")
	fmt.Print("\033[1;36m7. Logout 🔑\033[0m\n")
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

func playSlotMachine(currentUser *users.User) {
	game := slots.New(currentUser)
	game.Play()
}

func playFastestTyper(currentUser *users.User) {
	game := fastesttyper.New(currentUser)
	game.Play()
}

func playMazeGame(currentUser *users.User) {
	game := maze.New(currentUser)
	game.Play()
}

func playBattleship(currentUser *users.User) {
	game := battleship.New(currentUser)
	game.Play()
}
