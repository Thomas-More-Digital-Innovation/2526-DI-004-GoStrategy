import { gameStore } from "$lib/state/game.svelte";
import { toastStore } from "$lib/state/toast.svelte";
import { gamemodes } from "$lib/data/gamemodes.data";
import type { Position } from "$lib/types/game";
import { GameConnectionManager } from "./game-connection.svelte";
import {
    getDisabledRows,
    getHighlightedRows,
    getHighlightColor,
} from "./setup-helpers";
import {
    handlePieceCellClick,
    handleSetupPieceSwap,
} from "./board-interaction";

export class GameSessionController {
    connection = new GameConnectionManager();
    gameId = $state("");
    seatIndex = $state(-1);

    validMoves = $state<Position[]>([]);
    setupSwapPos1 = $state<Position | null>(null);
    setupSelectedPlayer = $state(0);

    get connected() {
        return this.connection.connected;
    }
    get isReconnecting() {
        return this.connection.isReconnecting;
    }
    get reconnectAttempts() {
        return this.connection.reconnectAttempts;
    }
    get maxReconnectAttempts() {
        return this.connection.maxReconnectAttempts;
    }

    isSetupPhase = $derived(gameStore.gameState?.isSetupPhase ?? false);

    isHumanTurn = $derived.by(() => {
        return (
            gameStore.gameState?.currentPlayerId === 0 &&
            gameStore.gameMode.mode === gamemodes.human_vs_ai.mode &&
            !gameStore.gameState?.isGameOver &&
            !this.isSetupPhase
        );
    });

    viewerId = $derived(
        gameStore.gameMode.mode === gamemodes.human_vs_ai.mode ? 0 : -1,
    );

    disabledRows = $derived(
        getDisabledRows(
            this.isSetupPhase,
            gameStore.gameMode,
            this.setupSelectedPlayer,
        ),
    );

    visualDisabledRows = $derived(this.isSetupPhase ? [4, 5] : []);

    highlightedRows = $derived(
        getHighlightedRows(
            this.isSetupPhase,
            gameStore.gameMode,
            this.setupSelectedPlayer,
        ),
    );

    highlightColor = $derived(
        getHighlightColor(
            this.isSetupPhase,
            gameStore.gameMode,
            this.setupSelectedPlayer,
            gameStore.gameState?.currentPlayerId,
        ),
    );

    async init(id: string, params: URLSearchParams) {
        this.gameId = id;
        gameStore.gameMode = gamemodes.fromString(params.get("mode") || "");

        const seatParam = params.get("seat");
        const defaultSeat = gameStore.gameMode.mode === "human_vs_ai" ? 0 : -1;
        this.seatIndex = seatParam !== null ? parseInt(seatParam) : defaultSeat;

        this.setupHandlers();
        await this.connection.connect(this.gameId, this.seatIndex);
    }

    destroy() {
        this.connection.disconnect();
        gameStore.reset();
    }

    setupHandlers() {
        const socket = this.connection.socket;
        socket.on("gameState", (data) => gameStore.updateGameState(data));
        socket.on("boardState", (data) =>
            gameStore.updateBoardState(data, this.viewerId),
        );
        socket.on("moveHistory", (data) =>
            gameStore.loadMoveHistory(data, this.gameId, this.viewerId),
        );

        socket.on("moveResult", (data) => {
            if (!data.success) {
                toastStore.handleApiMessage(data.error || data, "Move failed");
                gameStore.setSelectedPosition(null);
                this.validMoves = [];
            }
        });

        socket.on("validMoves", (data) => {
            this.validMoves = data.validMoves || [];
        });

        socket.on("combat", (data) => {
            gameStore.showCombatAnimation({
                attacker: data.attacker,
                defender: data.defender,
                attackerWon: data.attackerWon,
                defenderWon: data.defenderWon,
            });
        });

        socket.on("gameOver", () => {
            toastStore.success(
                "Game Over! You can review the game replay or go back to the menu.",
                5000,
            );
        });

        socket.on("error", (data) => {
            toastStore.handleApiMessage(data.error || data);
            gameStore.setSelectedPosition(null);
            this.validMoves = [];
        });

        socket.onClose(() => {
            if (this.connection.connected && !gameStore.gameState?.isGameOver) {
                this.connection.connected = false;
                this.connection.attemptReconnect(this.gameId, this.seatIndex);
            }
        });
    }

    attemptReconnect() {
        return this.connection.attemptReconnect(this.gameId, this.seatIndex);
    }

    abandonAndQuit() {
        return this.connection.abandonAndQuit(this.gameId);
    }

    handleCellClick(x: number, y: number) {
        if (this.isSetupPhase) {
            const pid =
                gameStore.gameMode.mode === gamemodes.human_vs_ai.mode
                    ? 0
                    : this.setupSelectedPlayer;
            this.setupSwapPos1 = handleSetupPieceSwap(
                x,
                y,
                pid,
                this.setupSwapPos1,
                this.connection.socket,
            );
            return;
        }

        if (this.isHumanTurn) {
            const { newSelected, newValidMoves } = handlePieceCellClick(
                x,
                y,
                gameStore.selectedPosition,
                this.validMoves,
                this.connection.socket,
            );
            gameStore.setSelectedPosition(newSelected);
            this.validMoves = newValidMoves;
        }
    }

    handleCellDragStart(e: DragEvent, x: number, y: number) {
        if (!this.isSetupPhase) return;
        e.dataTransfer?.setData("text/plain", JSON.stringify({ x, y }));
    }

    handleCellDrop(e: DragEvent, x: number, y: number) {
        if (!this.isSetupPhase) return;
        const data = e.dataTransfer?.getData("text/plain");
        if (!data) return;

        try {
            const from = JSON.parse(data) as Position;
            if (from.x === x && from.y === y) return;
            this.connection.socket.sendSwapPieces(from, { x, y });
        } catch (err) {
            console.error("Failed to parse drop data", err);
        }
    }

    handleRandomize(playerId?: number) {
        this.connection.socket.sendRandomizeSetup(playerId);
        this.setupSwapPos1 = null;
    }

    handleStartGame(headless: boolean = false) {
        this.connection.socket.sendStartGame(headless);
        this.setupSwapPos1 = null;
    }

    handleLoadSetup(setupData: string, playerId?: number) {
        this.connection.socket.sendLoadSetup(setupData, playerId);
        this.setupSwapPos1 = null;
    }

    handleSetSpeed(speedMs: number) {
        this.connection.socket.sendSetSpeed(speedMs);
    }

    handleStep() {
        gameStore.isStepping = true;
        this.connection.socket.sendStep();
    }

    handleTogglePause() {
        if (gameStore.isPaused) {
            this.connection.socket.sendUnpause();
        } else {
            this.connection.socket.sendPause();
        }
    }

    handleAnimationComplete() {
        this.connection.socket.sendAnimationComplete();
        gameStore.hideCombatAnimation();
    }

    handlePreviousMove() {
        if (!gameStore.isPaused) this.connection.socket.sendPause();
        gameStore.previousMove();
    }

    handleNextMove() {
        if (!gameStore.isPaused) this.connection.socket.sendPause();
        gameStore.nextMove();
    }

    handleGoToMove(index: number) {
        if (!gameStore.isPaused) this.connection.socket.sendPause();
        gameStore.goToMove(index);
    }

    handleExitReplay() {
        this.connection.socket.sendUnpause();
        gameStore.exitReplay();
    }

    saveGame() {
        try {
            const data = gameStore.exportGame();
            if (!data) {
                toastStore.warning("No history available to save");
                return;
            }
            const blob = new Blob([data], { type: "application/json" });
            const url = URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;
            a.download = `gostrategy-${this.gameId}-${Date.now()}.json`;
            a.click();
            URL.revokeObjectURL(url);
        } catch {
            toastStore.error("Failed to save game");
        }
    }
}
