package game

import "fmt"

// IsWaitingForInput returns true if the game is waiting for human input
func (gr *Runner) IsWaitingForInput() bool {
	gr.lock()
	defer gr.unlock()
	return gr.isWaitingForInput()
}

func (gr *Runner) isWaitingForInput() bool {
	gr.stateMu.RLock()
	defer gr.stateMu.RUnlock()
	return gr.waitingForHumanInput
}

// DebugSetWaitingForInput sets the waiting for human input flag to the given value.
// This is for debugging (& testing) purposes only and should not be used in production code.
func (gr *Runner) DebugSetWaitingForInput(value bool) {
	gr.stateMu.Lock()
	defer gr.stateMu.Unlock()
	gr.waitingForHumanInput = value
}

// Stop stops the game runner
func (gr *Runner) Stop() {
	gr.stateMu.Lock()
	gr.stopped = true
	gr.stateMu.Unlock()

	select {
	case gr.stopChan <- true:
	default:
	}
}

// SubmitHumanMove allows external code to submit a human player's move
func (gr *Runner) SubmitHumanMove(move Move) error {
	gr.lock()

	if !gr.isWaitingForInput() {
		gr.unlock()
		return fmt.Errorf("not waiting for input")
	}

	controller := gr.game.GetCurrentController()
	if controller.GetControllerType() != HumanController {
		gr.unlock()
		return fmt.Errorf("current player is not human")
	}

	humanController, ok := controller.(*HumanPlayerController)
	if !ok {
		gr.unlock()
		return fmt.Errorf("invalid controller type")
	}

	if move.GetPlayer() != gr.game.CurrentPlayer {
		gr.unlock()
		return fmt.Errorf("move player does not match current player")
	}

	if !gr.game.Board.IsValidMove(&move) || gr.game.IsTwoSquareViolation(gr.game.CurrentPlayer.GetID(), move) {
		gr.unlock()
		return fmt.Errorf("invalid move")
	}

	humanController.SetPendingMove(move)
	gr.unlock()

	gr.ExecuteTurn()
	return nil
}

// Pause pauses the game runner
func (gr *Runner) Pause() {
	gr.lock()
	defer gr.unlock()
	gr.setPaused(true)
}

func (gr *Runner) setPaused(paused bool) {
	gr.stateMu.Lock()
	defer gr.stateMu.Unlock()
	gr.paused = paused
}

// Unpause unpauses the game runner
func (gr *Runner) Unpause() {
	gr.lock()
	defer gr.unlock()
	gr.setPaused(false)
}

// IsPaused returns whether the game runner is paused
func (gr *Runner) IsPaused() bool {
	gr.lock()
	defer gr.unlock()
	return gr.isPaused()
}

func (gr *Runner) isPaused() bool {
	gr.stateMu.RLock()
	defer gr.stateMu.RUnlock()
	return gr.paused
}

// IsStopped returns whether the game runner has been stopped
func (gr *Runner) IsStopped() bool {
	gr.stateMu.RLock()
	defer gr.stateMu.RUnlock()
	return gr.stopped
}
