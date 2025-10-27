package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/umbralcalc/dexetera/pkg/games"
)

func main() {
	var (
		gameType  = flag.String("game", "minimal_example", "Game type to generate (minimal_example, builder_example, visualization_example)")
		outputDir = flag.String("output", "./generated_game", "Output directory for generated files")
		listGames = flag.Bool("list", false, "List available games")
	)
	flag.Parse()

	if *listGames {
		listAvailableGames()
		return
	}

	// Create the game based on type
	var game games.Game
	switch *gameType {
	case "minimal_example":
		game = games.NewMinimalExampleGame()
	case "builder_example":
		game = games.NewBuilderExampleGame()
	case "visualization_example":
		game = games.NewVisualizationExampleGame()
	default:
		fmt.Printf("❌ Unknown game type: %s\n", *gameType)
		fmt.Println("Available games:")
		listAvailableGames()
		os.Exit(1)
	}

	// Generate the game package
	fmt.Printf("🎮 Generating %s game package...\n", game.GetName())
	fmt.Printf("📁 Output directory: %s\n", *outputDir)

	if err := games.GenerateGamePackage(game, *outputDir); err != nil {
		log.Fatalf("❌ Failed to generate game package: %v", err)
	}

	fmt.Println("✅ Game package generated successfully!")
	fmt.Printf("📝 To build and run:\n")
	fmt.Printf("   1. cd %s\n", *outputDir)
	fmt.Printf("   2. ./build.sh\n")
	fmt.Printf("   3. Start your Python websocket server\n")
	fmt.Printf("   4. Open index.html in a browser\n")
}

func listAvailableGames() {
	fmt.Println("Available games:")
	fmt.Println("  minimal_example      - Simple counter game")
	fmt.Println("  builder_example      - Complex game with multiple partitions")
	fmt.Println("  visualization_example - Advanced visualization demo")
}
