// Package main is the entry point for autonomous simulation
package main

import (
	aivsai "digital-innovation/gostrategy/internal/ai/AIvsAI"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/models"
	"flag"
	"fmt"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func run() error {
	aiTypes := flag.String("ai", "fafo:fafo", "Run AI vs AI matches")
	matches := flag.Int("matches", 100, "Number of AI vs AI matches to run")
	format := flag.String("format", "none", "The format used to print the results of an AI vs AI competition, either none or md")
	loggingEnabled := flag.Bool("logging", true, "Show logs in stdout")
	setupPath := flag.String("setup", "", "Path to a board setup file for both AIs (JSON or text)")
	setupPath1 := flag.String("setup1", "", "Path to a board setup file for AI 1 (Alice)")
	setupPath2 := flag.String("setup2", "", "Path to a board setup file for AI 2 (Bob)")

	flag.Parse()

	fmt.Println("=== GoStrategy Autonomous Simulation Running ===")

	var ai1, ai2 string
	switch {
	case aiTypes == nil:
		fmt.Print("No AIs specified, running default FATO vs FATO")
		ai1, ai2 = models.Fato, models.Fato
	case strings.Contains(*aiTypes, ":"):
		aiTypeSplit := strings.Split(*aiTypes, ":")
		ai1, ai2 = aiTypeSplit[0], aiTypeSplit[1]
	default:
		ai1, ai2 = *aiTypes, *aiTypes
	}

	p1SetupFile := *setupPath1
	p2SetupFile := *setupPath2
	if p1SetupFile == "" && *setupPath != "" {
		p1SetupFile = *setupPath
	}
	if p2SetupFile == "" && *setupPath != "" {
		p2SetupFile = *setupPath
	}

	var setup1, setup2 []string
	if p1SetupFile != "" {
		s1, err := game.LoadBoardSetupFromFile(p1SetupFile)
		if err != nil {
			return fmt.Errorf("failed to load AI 1 board setup: %w", err)
		}
		setup1 = s1
	}

	if p2SetupFile != "" {
		s2, err := game.LoadBoardSetupFromFile(p2SetupFile)
		if err != nil {
			return fmt.Errorf("failed to load AI 2 board setup: %w", err)
		}
		setup2 = s2
	}

	start := time.Now()
	aivsai.RunAIvsAIWithSetups(ai1, ai2, *matches, *format, *loggingEnabled, setup1, setup2)
	elapsed := time.Since(start)
	fmt.Printf("\nAI vs AI matches completed in %.2f seconds\n", elapsed.Seconds())
	return nil
}
