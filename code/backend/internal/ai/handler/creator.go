// Package AIhandler provides a factory for creating AI instances.
package AIhandler

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/ai/fafo"
	"digital-innovation/gostrategy/internal/ai/fato"
	"digital-innovation/gostrategy/internal/ai/heuristic"
	"digital-innovation/gostrategy/internal/ai/mcts"
	"digital-innovation/gostrategy/internal/ai/minimax"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"fmt"
)

const defaultName = "default"

func optFloat64(opts map[string]any, key string, def float64) float64 {
	if v, ok := opts[key].(float64); ok {
		return v
	}
	return def
}

func optString(opts map[string]any, key string, def string) string {
	if v, ok := opts[key].(string); ok {
		return v
	}
	return def
}

func setOptFloat64(opts map[string]any, key string, target map[string]any) {
	if v, ok := opts[key].(float64); ok {
		target[key] = v
	}
}

// CreateAI returns an AI instance based on the given type with default options.
func CreateAI(aiType string, player *game.Player) (ai.AI, error) {
	return CreateAIWithOptions(aiType, player, nil)
}

// CreateAIWithOptions creates an AI instance with customizable configuration options.
func CreateAIWithOptions(aiType string, player *game.Player, options map[string]any) (ai.AI, error) {
	switch aiType {
	case models.Fafo:
		return fafo.NewAI(player, false), nil
	case models.Fato:
		return fato.NewAIWithAggression(player, true, optFloat64(options, "aggression", 0.5)), nil
	case models.Heuristic:
		params, _ := ai.Load(models.Heuristic, optString(options, "name", defaultName))
		params.Aggression = optFloat64(options, "aggression", params.Aggression)
		return heuristic.NewAIWithParams(player, true, params), nil
	case models.Minimax:
		params, _ := ai.Load(models.Minimax, optString(options, "name", defaultName))
		params.Aggression = optFloat64(options, "aggression", params.Aggression)
		setOptFloat64(options, "depth", params.Config)
		return minimax.NewAIWithParams(player, true, params), nil
	case models.Mcts:
		params, _ := ai.Load(models.Mcts, optString(options, "name", defaultName))
		params.Aggression = optFloat64(options, "aggression", params.Aggression)
		setOptFloat64(options, "iterations", params.Config)
		return mcts.NewAIWithParams(player, true, params), nil
	default:
		return nil, fmt.Errorf("i don't know that AI! %s", aiType)
	}
}
