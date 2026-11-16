import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { GameState } from '#lib/types/game.js';
import { gamemodes } from '#lib/data/gamemodes.data.js';
import { toastStore } from '#lib/state/toast.svelte.js';
import { GameConnectionManager } from './game-connection.svelte';
import { GameSessionController } from './game-session.svelte';
import { initStaticReplay } from './static-replay';
import { createFakeSocket } from './fake-socket';

vi.mock('./static-replay', () => ({ initStaticReplay: vi.fn() }));

const state = (over: Partial<GameState> = {}) =>
	({ isSetupPhase: false, currentPlayerId: 0, paused: false, ...over }) as GameState;

describe('GameSessionController', () => {
	let socket: ReturnType<typeof createFakeSocket>;
	let session: GameSessionController;

	beforeEach(() => {
		toastStore.toasts = [];
		socket = createFakeSocket();
		session = new GameSessionController(new GameConnectionManager(socket));
	});

	describe('init', () => {
		it('delegates replays to initStaticReplay', async () => {
			await session.init('g', new URLSearchParams('replay=true'));
			expect(initStaticReplay).toHaveBeenCalledWith(session, 'g');
			expect(socket.connect).not.toHaveBeenCalled();
		});

		it.each([
			['mode=human_vs_ai', 0],
			['mode=ai_vs_ai', -1],
			['mode=ai_vs_ai&seat=1', 1]
		])('%s connects with seat %i', async (query, seat) => {
			await session.init('g', new URLSearchParams(query));
			expect(socket.connect).toHaveBeenCalledWith('g', seat);
			expect(session.connected).toBe(true);
		});
	});

	describe('socket messages', () => {
		beforeEach(() => session.init('g', new URLSearchParams('mode=human_vs_ai')));

		it('updates the store from game and board state', () => {
			socket.emit('gameState', state({ round: 3 }));
			expect(session.store.gameState?.round).toBe(3);
			expect(session.isHumanTurn).toBe(true);
		});

		it('stores valid moves', () => {
			socket.emit('validMoves', { validMoves: [{ x: 1, y: 1 }] });
			expect(session.validMoves).toEqual([{ x: 1, y: 1 }]);
			socket.emit('validMoves', {});
			expect(session.validMoves).toEqual([]);
		});

		it.each(['moveResult', 'error'])('%s failure clears selection and toasts', (type) => {
			session.validMoves = [{ x: 1, y: 1 }];
			session.store.setSelectedPosition({ x: 0, y: 0 });
			socket.emit(type, { success: false, error: 'illegal' });
			expect(session.validMoves).toEqual([]);
			expect(session.store.selectedPosition).toBeNull();
			expect(toastStore.toasts[0].message).toBe('illegal');
		});

		it('shows combat animations', () => {
			socket.emit('combat', {
				attacker: 'a',
				defender: 'd',
				attackerWon: true,
				defenderWon: false
			});
			expect(session.store.combatAnimation?.attackerWon).toBe(true);
		});

		it('reconnects on unexpected close but not after game over', () => {
			const reconnect = vi.spyOn(session.connection, 'attemptReconnect').mockResolvedValue();
			socket.close();
			expect(reconnect).toHaveBeenCalledWith('g', 0);
			session.connection.connected = true;
			session.store.gameState = state({ isGameOver: true });
			socket.close();
			expect(reconnect).toHaveBeenCalledTimes(1);
		});
	});

	describe('derived board highlighting', () => {
		it('reflects setup phase and mode', () => {
			session.store.gameMode = gamemodes.human_vs_ai;
			session.store.gameState = state({ isSetupPhase: true });
			expect(session.isSetupPhase).toBe(true);
			expect(session.viewerId).toBe(0);
			expect(session.visualDisabledRows).toEqual([4, 5]);
			expect(session.highlightedRows).toEqual([6, 7, 8, 9]);
			expect(session.highlightColor).toBe('red');
		});
	});

	describe('handleCellClick', () => {
		it('swaps pieces during setup', () => {
			session.store.gameMode = gamemodes.human_vs_ai;
			session.store.gameState = state({ isSetupPhase: true });
			session.handleCellClick(0, 6);
			session.handleCellClick(1, 6);
			expect(socket.sendSwapPieces).toHaveBeenCalledWith({ x: 0, y: 6 }, { x: 1, y: 6 });
		});

		it('warns when paused', () => {
			session.store.gameState = state({ paused: true });
			session.handleCellClick(0, 0);
			expect(toastStore.toasts[0].type).toBe('warning');
		});

		it('selects on the human turn only', () => {
			session.store.gameMode = gamemodes.human_vs_ai;
			session.store.gameState = state();
			session.store.boardState = {
				board: [[{ ownerId: 0, ownerName: 'me', revealed: true, position: { x: 0, y: 0 } }]],
				width: 10,
				height: 10
			};
			session.handleCellClick(0, 0);
			expect(session.store.selectedPosition).toEqual({ x: 0, y: 0 });
			session.store.gameState = state({ currentPlayerId: 1 });
			session.handleCellClick(0, 0);
			expect(socket.requestValidMoves).toHaveBeenCalledTimes(1);
		});
	});

	describe('controls', () => {
		it('forwards setup and playback commands to the socket', () => {
			session.setupSwapPos1 = { x: 0, y: 0 };
			session.handleRandomize(1);
			expect(socket.sendRandomizeSetup).toHaveBeenCalledWith(1);
			expect(session.setupSwapPos1).toBeNull();
			session.handleStartGame(true);
			expect(socket.sendStartGame).toHaveBeenCalledWith(true);
			session.handleLoadSetup('data', 0);
			expect(socket.sendLoadSetup).toHaveBeenCalledWith('data', 0);
			session.handleSetSpeed(100);
			expect(socket.sendSetSpeed).toHaveBeenCalledWith(100);
			session.handleStep();
			expect(session.store.isStepping).toBe(true);
			session.handleAnimationComplete();
			expect(socket.sendAnimationComplete).toHaveBeenCalled();
		});

		it('toggles pause based on current state', () => {
			session.store.gameState = state();
			session.handleTogglePause();
			expect(socket.sendPause).toHaveBeenCalled();
			session.store.gameState = state({ paused: true });
			session.handleTogglePause();
			expect(socket.sendUnpause).toHaveBeenCalled();
		});

		it('pauses a live game when browsing history', () => {
			session.store.gameState = state();
			session.store.history = [0, 1, 2].map((i) => ({
				moveNumber: i,
				boardState: [],
				move: { moveIndex: i } as never
			}));
			session.handleGoToMove(0);
			expect(socket.sendPause).toHaveBeenCalled();
			expect(session.store.currentHistoryIndex).toBe(0);
		});

		it('destroy disconnects and resets', () => {
			session.store.gameState = state();
			session.destroy();
			expect(socket.disconnect).toHaveBeenCalled();
			expect(session.store.gameState).toBeNull();
		});

		it('warns when saving without history', () => {
			session.saveGame();
			expect(toastStore.toasts[0].message).toContain('No history');
		});
	});
});
