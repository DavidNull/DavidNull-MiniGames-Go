# 🎲 Casino Games Console

A collection of casino-style games that can be played in the console.

## Available Games

1. **Dice Game** - Bet on a number and roll the dice. Win 5x your bet if you guess correctly!
2. **Blackjack** - Coming soon
3. **Slot Machine** - Coming soon

## How to Run

```bash
# Build the game
go build -o casino ./cmd/app

# Run the game
./casino
```

## Project Structure

- `cmd/app/` - Main application entry point
- `internal/games/` - Game implementations
  - `dice/` - Dice game
  - `blackjack/` - Blackjack game (coming soon)
  - `slots/` - Slot machine game (coming soon)

## Adding New Games

To add a new game:

1. Create a new directory under `internal/games/`
2. Implement your game logic
3. Add an entry to the menu in `cmd/app/main.go`

## License

MIT