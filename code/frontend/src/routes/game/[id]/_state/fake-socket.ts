import { vi } from 'vitest';
import type { GameSocket } from '#lib/api/websocket.js';

/** Test double for GameSocket: records sends and lets tests emit server messages. */
export function createFakeSocket() {
	const handlers = new Map<string, (data: unknown) => void>();
	let closeHandler: (() => void) | null = null;

	const fake = {
		connect: vi.fn<GameSocket['connect']>().mockResolvedValue(undefined),
		disconnect: vi.fn(),
		on: vi.fn((type: string, handler: (data: unknown) => void) => handlers.set(type, handler)),
		onClose: vi.fn((handler: () => void) => (closeHandler = handler)),
		sendMove: vi.fn(),
		requestValidMoves: vi.fn(),
		sendSwapPieces: vi.fn(),
		sendRandomizeSetup: vi.fn(),
		sendStartGame: vi.fn(),
		sendAnimationComplete: vi.fn(),
		sendPause: vi.fn(),
		sendUnpause: vi.fn(),
		sendLoadSetup: vi.fn(),
		sendSetSpeed: vi.fn(),
		sendStep: vi.fn(),
		emit: (type: string, data?: unknown) => handlers.get(type)?.(data),
		close: () => closeHandler?.()
	};
	return fake as typeof fake & GameSocket;
}
