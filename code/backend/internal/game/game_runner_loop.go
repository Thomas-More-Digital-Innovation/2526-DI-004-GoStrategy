package game

import (
	"digital-innovation/gostrategy/internal/logging"
	// nosemgrep: go.lang.security.audit.crypto.math_random.math-random-used
	"math/rand"
	"time"
)

// RunToCompletion runs the game until it's over (for AI vs AI)
// Winner can be nil when max turns are reached and both AIs have a similar piece count
func (gr *Runner) RunToCompletion() *Player {
	turnCount := 0
	logging.Debug(logging.TagGame, "Runner: Starting RunToCompletion loop")

	for !gr.isGameOverSafe() && turnCount < gr.maxTurns && !gr.IsStopped() {
		ok, executed := gr.runOneTurn()
		if !ok {
			logging.Debug(logging.TagGame, "Runner: Stop signal received, ending game")
			return nil
		}
		if executed {
			turnCount++
		}
	}

	if turnCount >= gr.maxTurns {
		logging.Debug(logging.TagGame, "Runner: Game ended: Maximum turns reached")
		return gr.calculateWinnerOnMaxTurnsExceeded()
	}

	gr.lock()
	defer gr.unlock()
	return gr.game.GetWinner()
}

// runOneTurn attempts one game turn. Returns (continue, executed).
// continue=false means a stop signal was received.
func (gr *Runner) runOneTurn() (bool, bool) {
	select {
	case <-gr.stopChan:
		return false, false
	default:
	}

	if gr.IsPaused() {
		select {
		case <-gr.stopChan:
			return false, false
		case <-time.After(100 * time.Millisecond):
			return true, false
		}
	}

	if gr.ExecuteTurn() {
		return true, true
	}

	if gr.isGameOverSafe() {
		logging.Debug(logging.TagGame, "Runner: Game ended during ExecuteTurn")
		return true, false
	}

	select {
	case <-gr.stopChan:
		return false, false
	case <-time.After(100 * time.Millisecond):
		return true, false
	}
}

func (gr *Runner) calculateWinnerOnMaxTurnsExceeded() *Player {
	gr.lock()
	defer gr.unlock()
	if float64(gr.game.Players[0].GetPieceScore())/float64(gr.game.Players[1].GetPieceScore()) > 1.15 {
		gr.game.SetWinner(gr.game.Players[0], WinCauseMaxTurns)
	} else if float64(gr.game.Players[1].GetPieceScore())/float64(gr.game.Players[0].GetPieceScore()) > 1.15 {
		gr.game.SetWinner(gr.game.Players[1], WinCauseMaxTurns)
	}
	return gr.game.GetWinner()
}

func (gr *Runner) isGameOverSafe() bool {
	gr.lock()
	defer gr.unlock()
	return gr.game.IsGameOver()
}

// ExecuteTurn executes a single turn. Returns false if waiting for human input.
func (gr *Runner) ExecuteTurn() bool {
	return gr.executeTurn(false)
}

func (gr *Runner) executeTurn(ignorePause bool) bool {
	gr.lock()

	if gr.game.IsGameOver() {
		gr.unlock()
		return false
	}
	if !ignorePause && gr.paused {
		gr.unlock()
		return false
	}

	controller := gr.game.GetCurrentController()
	if controller.GetControllerType() == HumanController {
		return gr.executeHumanTurn(controller)
	}
	return gr.executeAITurn(controller, ignorePause)
}

// executeHumanTurn handles a human player's turn. Caller must hold the lock.
func (gr *Runner) executeHumanTurn(controller PlayerController) bool {
	humanController, ok := controller.(*HumanPlayerController)
	if !ok || !humanController.HasPendingMove() {
		gr.stateMu.Lock()
		if !gr.waitingForHumanInput {
			gr.waitingForHumanInput = true
		}
		gr.stateMu.Unlock()
		gr.unlock()
		return false
	}

	move := humanController.GetPendingMove()
	if move == nil {
		gr.unlock()
		return false
	}

	piece := gr.game.Board.GetPieceAt(move.GetFrom())
	if piece == nil {
		gr.unlock()
		return false
	}

	gr.game.MakeMove(move, piece)
	gr.stateMu.Lock()
	gr.waitingForHumanInput = false
	gr.stateMu.Unlock()
	gr.unlock()

	if gr.onMoveExecuted != nil {
		gr.onMoveExecuted()
	}
	return true
}

// executeAITurn handles an AI player's turn. Caller must hold the lock.
func (gr *Runner) executeAITurn(controller PlayerController, ignorePause bool) bool {
	boardCopy := gr.game.Board.Clone()
	currentPlayer := gr.game.CurrentPlayer
	gr.unlock()

	start := time.Now()
	move := controller.MakeMove(boardCopy)

	if !ignorePause && !gr.applyTurnPacing(time.Since(start)) {
		return false
	}

	gr.lock()

	if gr.game.IsGameOver() || gr.game.CurrentPlayer != currentPlayer {
		gr.unlock()
		return false
	}

	select {
	case <-gr.stopChan:
		gr.unlock()
		return false
	default:
	}

	piece := gr.game.Board.GetPieceAt(move.GetFrom())
	if piece == nil || piece.GetOwner() != gr.game.CurrentPlayer {
		gr.game.SetWinner(gr.getOpponent(gr.game.CurrentPlayer), WinCauseNoMovablePieces)
		gr.unlock()
		return false
	}

	if !gr.game.Board.IsValidMove(&move) || gr.game.IsTwoSquareViolation(gr.game.CurrentPlayer.GetID(), move) {
		gr.game.SetWinner(gr.getOpponent(gr.game.CurrentPlayer), WinCauseNoMovablePieces)
		gr.unlock()
		return false
	}

	gr.game.MakeMove(&move, piece)
	gr.unlock()

	if gr.onMoveExecuted != nil {
		gr.onMoveExecuted()
	}
	return true
}

// applyTurnPacing sleeps for the configured turn delay minus elapsed time.
// Returns false if stopped or paused during the wait.
func (gr *Runner) applyTurnPacing(elapsed time.Duration) bool {
	if gr.turnDelay <= 0 {
		return true
	}
	delay := gr.turnDelay
	if gr.turnDelay >= 100*time.Millisecond {
		// #nosec G404 - weak random is sufficient for turn delay pacing
		delay = time.Duration(float64(gr.turnDelay)*0.8 + float64(rand.Intn(int(gr.turnDelay)))*0.4)
	}
	if sleepTime := delay - elapsed; sleepTime > 0 {
		select {
		case <-gr.stopChan:
			return false
		case <-time.After(sleepTime):
		}
	}
	return !gr.isPaused()
}

// Step executes a single turn even if the game is paused
func (gr *Runner) Step() bool {
	return gr.executeTurn(true)
}
