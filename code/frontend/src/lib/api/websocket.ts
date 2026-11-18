import type { Position } from '#lib/types/game.js';
import {
	MsgTypeAnimationComplete,
	MsgTypeBoardState,
	MsgTypeCombat,
	MsgTypeError,
	MsgTypeGameOver,
	MsgTypeGameState,
	MsgTypeGetValidMoves,
	MsgTypeLoadSetup,
	MsgTypeMove,
	MsgTypeMoveHistory,
	MsgTypeMoveResult,
	MsgTypePause,
	MsgTypeRandomizeSetup,
	MsgTypeSetSpeed,
	MsgTypeStartGame,
	MsgTypeStep,
	MsgTypeSwapPieces,
	MsgTypeUnpause,
	MsgTypeValidMoves,
	type BoardStateMessage,
	type CombatMessage,
	type ErrorMessage,
	type GameOverMessage,
	type GameStateMessage,
	type GetValidMovesMessage,
	type LoadSetupMessage,
	type MoveHistoryMessage,
	type MoveMessage,
	type MoveResultMessage,
	type RandomizeSetupMessage,
	type SetSpeedMessage,
	type StartGameMessage,
	type SwapPiecesMessage,
	type ValidMovesMessage
} from '#lib/types/generated/dto.js';

const WS_BASE = import.meta.env.VITE_WS_BASE || 'ws://localhost:8080';

/** Payload of every server-to-client message, keyed by message type. */
export interface ServerMessages {
	[MsgTypeGameState]: GameStateMessage;
	[MsgTypeBoardState]: BoardStateMessage;
	[MsgTypeMoveHistory]: MoveHistoryMessage;
	[MsgTypeMoveResult]: MoveResultMessage;
	[MsgTypeValidMoves]: ValidMovesMessage;
	[MsgTypeCombat]: CombatMessage;
	[MsgTypeGameOver]: GameOverMessage;
	[MsgTypeError]: ErrorMessage;
}

type MessageHandler<K extends keyof ServerMessages> = (data: ServerMessages[K]) => void;

export class GameSocket {
	private ws: WebSocket | null = null;
	private handlers = new Map<string, (data: unknown) => void>();
	private closeHandler: (() => void) | null = null;

	connect(gameId: string, playerId: number = -1): Promise<void> {
		return new Promise((resolve, reject) => {
			const url = `${WS_BASE}/game/${gameId}?player=${playerId}`;
			this.ws = new WebSocket(url);

			this.ws.onopen = () => resolve();
			this.ws.onerror = (e) => reject(e);

			this.ws.onmessage = (event) => {
				try {
					const message = JSON.parse(event.data);
					this.handlers.get(message.type)?.(message.data);
				} catch (error) {
					console.error('Failed to parse WebSocket message:', error);
				}
			};

			this.ws.onclose = () => {
				console.log('WebSocket closed');
				this.closeHandler?.();
			};
		});
	}

	onClose(handler: () => void) {
		this.closeHandler = handler;
	}

	on<K extends keyof ServerMessages>(type: K, handler: MessageHandler<K>) {
		// why: the server guarantees the payload shape per type, so the erased handler is safe to store
		this.handlers.set(type, handler as (data: unknown) => void);
	}

	private send(type: string, data?: object) {
		if (!this.ws) return;
		this.ws.send(JSON.stringify({ type, data }));
	}

	sendMove(from: Position, to: Position) {
		this.send(MsgTypeMove, { from, to } satisfies MoveMessage);
	}

	requestValidMoves(position: Position) {
		this.send(MsgTypeGetValidMoves, { position } satisfies GetValidMovesMessage);
	}

	sendSwapPieces(pos1: Position, pos2: Position) {
		this.send(MsgTypeSwapPieces, { pos1, pos2 } satisfies SwapPiecesMessage);
	}

	sendRandomizeSetup(playerId?: number) {
		this.send(MsgTypeRandomizeSetup, { playerId } satisfies RandomizeSetupMessage);
	}

	sendStartGame(headless: boolean = false) {
		this.send(MsgTypeStartGame, { headless } satisfies StartGameMessage);
	}

	sendAnimationComplete() {
		this.send(MsgTypeAnimationComplete);
	}

	sendPause() {
		this.send(MsgTypePause);
	}

	sendUnpause() {
		this.send(MsgTypeUnpause);
	}

	sendLoadSetup(setupData: string, playerId?: number) {
		this.send(MsgTypeLoadSetup, { setupData, playerId } satisfies LoadSetupMessage);
	}

	sendSetSpeed(speedMs: number) {
		this.send(MsgTypeSetSpeed, { speedMs } satisfies SetSpeedMessage);
	}

	sendStep() {
		this.send(MsgTypeStep);
	}

	disconnect() {
		this.ws?.close();
		this.ws = null;
		this.handlers.clear();
		this.closeHandler = null;
	}
}
