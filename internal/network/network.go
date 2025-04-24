package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"davgames/internal/users"
)

var (
	activeGames      = make(map[string]GameSession)
	activeGamesMutex sync.Mutex
)

type GameSession struct {
	ID         string   `json:"id"`
	GameType   string   `json:"game_type"`
	Players    []string `json:"players"`
	InProgress bool     `json:"in_progress"`
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DavidNull Games Server - Running")
}

func HandleUsers(usersData *users.Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			userList := make([]string, 0)
			for _, user := range usersData.Users {
				userList = append(userList, user.Name)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(userList)

		case http.MethodPost:

			var authRequest struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			user, err := usersData.Authenticate(authRequest.Username, authRequest.Password)
			if err != nil {
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"name":    user.Name,
				"balance": user.Balance,
			})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func HandleGames(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		activeGamesMutex.Lock()
		games := make([]GameSession, 0, len(activeGames))
		for _, game := range activeGames {
			games = append(games, game)
		}
		activeGamesMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)

	case http.MethodPost:
		var newGame GameSession
		if err := json.NewDecoder(r.Body).Decode(&newGame); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		activeGamesMutex.Lock()
		activeGames[newGame.ID] = newGame
		activeGamesMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newGame)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func StartClient(serverIP, port string) {
	fmt.Println("\033[1;32m🚀 Starting client...\033[0m")
	fmt.Println("\033[1;33m🔌 Connected to server:", serverIP, ":", port, "\033[0m")

	resp, err := http.Get(fmt.Sprintf("http://%s:%s/games", serverIP, port))
	if err != nil {
		fmt.Printf("\033[1;31m❌ Error getting games: %v\033[0m\n", err)
		return
	}
	defer resp.Body.Close()

	var games []GameSession
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		fmt.Printf("\033[1;31m❌ Error decoding games: %v\033[0m\n", err)
		return
	}

	fmt.Println("\033[1;36m🎮 Available Games:\033[0m")
	if len(games) == 0 {
		fmt.Println("\033[1;33m📭 No active games currently\033[0m")
	} else {
		for i, game := range games {
			status := "Waiting for players"
			if game.InProgress {
				status = "In Progress"
			}
			fmt.Printf("\033[1;32m%d. %s - %s - Players: %d/%d\033[0m\n",
				i+1, game.GameType, status, len(game.Players), 2)
		}
	}

	fmt.Println("\n\033[1;33m🚧 Functionality under development...\033[0m")
	fmt.Println("\033[1;33m📚 Please check the documentation for more information.\033[0m")
}
