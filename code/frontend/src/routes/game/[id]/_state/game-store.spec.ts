import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { BoardState, GameState, Piece } from '#lib/types/game.js';
import type { HistoricalMove } from '#lib/replayEngine.js';
import { gamemodes } from '#lib/data/gamemodes.data.js';
import { GameStore } from './game-store.svelte';

const state = (over: Partial<GameState> = {}) => ({ isSetupPhase: false, ...over }) as GameState;
const move = (moveIndex: number, over: Partial<HistoricalMove> = {}): HistoricalMove => ({
	moveIndex,
	playerId: 0,
	fromX: 0,
	fromY: 0,
	toX: 1,
	toY: 0,
	result: 'move',
	...over
});
const piece = (ownerId: number): Piece => ({
	ownerId,
	rank: '5',
	type: 'Lieutenant',
	revealed: true,
	position: { x: 0, y: 0 }
});
const board = (lastMove?: HistoricalMove): BoardState => ({
	board: [[piece(0), piece(1)]],
	width: 10,
	height: 10,
	lastMove
});

describe('GameStore', () => {
	let store: GameStore;

	beforeEach(() => {
		store = new GameStore();
	});

	describe('setup countdown', () => {
		beforeEach(() => vi.useFakeTimers());
		afterEach(() => vi.useRealTimers());

		it('counts down during setup and stops at zero', () => {
			store.updateGameState(state({ isSetupPhase: true, setupRemainingSecs: 2 }));
			vi.advanceTimersByTime(2000);
			expect(store.setupRemainingSecs).toBe(0);
			vi.advanceTimersByTime(1000);
			expect(store.setupRemainingSecs).toBeNull();
		});

		it('clears the countdown once setup ends', () => {
			store.updateGameState(state({ isSetupPhase: true, setupRemainingSecs: 5 }));
			store.updateGameState(state());
			expect(store.setupRemainingSecs).toBeNull();
		});
	});

	describe('updateBoardState', () => {
		it('records live moves into history once per move index', () => {
			store.updateGameState(state());
			store.updateBoardState(board(move(1)));
			store.updateBoardState(board(move(1)));
			store.updateBoardState(board(move(2)));
			expect(store.history).toHaveLength(2);
			expect(store.currentHistoryIndex).toBe(1);
			expect(store.lastMove).toEqual(move(2));
		});

		it('ignores moves during setup', () => {
			store.updateGameState(state({ isSetupPhase: true }));
			store.updateBoardState(board(move(1)));
			expect(store.history).toEqual([]);
		});

		it('hides opponent pieces from the viewer in history', () => {
			store.updateGameState(state());
			const attacker = { type: 'Spy', rank: '1', ownerId: 0 };
			const defender = { type: 'Bomb', rank: 'B', ownerId: 1 };
			store.updateBoardState(board(move(1, { attacker, defender })), 0);
			const [entry] = store.history;
			expect(entry.boardState[0][0].rank).toBe('5');
			expect(entry.boardState[0][1].rank).toBeUndefined();
			expect(entry.move.attacker).toEqual(attacker);
			expect(entry.move.defender).toEqual({ ...defender, type: '', rank: '' });
		});

		it('keeps the replayed board while replaying', () => {
			store.updateGameState(state());
			store.updateBoardState(board(move(1)));
			store.goToMove(0);
			const live = board(move(2));
			store.updateBoardState(live);
			expect(store.lastLiveBoard).toBe(live);
			expect(store.currentHistoryIndex).toBe(0);
		});
	});

	describe('replay navigation', () => {
		beforeEach(() => {
			store.updateGameState(state());
			[1, 2, 3].forEach((i) => store.updateBoardState(board(move(i))));
		});

		it('steps back and forward within bounds', () => {
			store.previousMove();
			expect(store.isReplaying).toBe(true);
			expect(store.currentHistoryIndex).toBe(1);
			store.goToMove(0);
			store.previousMove();
			expect(store.currentHistoryIndex).toBe(0);
			store.nextMove();
			store.nextMove();
			store.nextMove();
			expect(store.currentHistoryIndex).toBe(2);
		});

		it('ignores out of range indices', () => {
			store.goToMove(-1);
			store.goToMove(3);
			expect(store.isReplaying).toBe(false);
		});

		it('returns to the live board on exit', () => {
			store.goToMove(0);
			store.exitReplay();
			expect(store.isReplaying).toBe(false);
			expect(store.currentHistoryIndex).toBe(2);
			expect(store.lastMove).toEqual(move(3));
		});

		it('stays in replay when the game is over', () => {
			store.gameState = state({ isGameOver: true });
			store.goToMove(0);
			store.exitReplay();
			expect(store.isReplaying).toBe(true);
		});
	});

	describe('loadMoveHistory', () => {
		const initialState = [[{ type: 'Scout', rank: '2', ownerId: 0 }, null]];
		const moves = [{ from: { x: 0, y: 0 }, to: { x: 1, y: 0 } }];

		it('rebuilds boards from the full history', () => {
			store.gameMode = gamemodes.ai_vs_ai;
			store.loadMoveHistory({ moves, fullHistory: [move(0)], initialState }, 'g');
			expect(store.history[0].boardState[0][1]).toMatchObject({ rank: '2', revealed: true });
			expect(store.currentHistoryIndex).toBe(0);
			expect(store.rawHistory?.gameId).toBe('g');
		});

		it('hides opponent pieces for a human viewer', () => {
			store.gameMode = gamemodes.human_vs_ai;
			store.loadMoveHistory({ moves, fullHistory: [move(0)], initialState }, 'g', 1);
			expect(store.history[0].boardState[0][1].revealed).toBe(false);
		});

		it('falls back to bare moves without full history', () => {
			store.loadMoveHistory({ moves, fullHistory: undefined as never, initialState: [] });
			expect(store.history[0].move).toMatchObject({ fromX: 0, toX: 1, result: 'move' });
			expect(store.history[0].boardState).toEqual([]);
		});

		it('jumps to the last move when the game is over', () => {
			store.gameState = state({ isGameOver: true });
			store.loadMoveHistory({ moves, fullHistory: [move(0)], initialState });
			expect(store.isReplaying).toBe(true);
			expect(store.boardState?.lastMove).toEqual(move(0));
		});
	});

	it('exports nothing without history, otherwise a replay payload', () => {
		expect(store.exportGame()).toBeNull();
		store.gameMode = gamemodes.ai_vs_ai;
		store.loadMoveHistory({
			moves: [],
			fullHistory: [],
			initialState: [[null]]
		});
		const exported = JSON.parse(store.exportGame()!);
		expect(exported).toMatchObject({
			gameType: gamemodes.ai_vs_ai.mode,
			player1Username: 'Player 1',
			winnerId: null
		});
		expect(exported.createdAt).toEqual(expect.any(String));
	});

	it('reset clears everything', () => {
		store.updateGameState(state());
		store.updateBoardState(board(move(1)));
		store.reset();
		expect(store.gameState).toBeNull();
		expect(store.history).toEqual([]);
		expect(store.gameMode).toBe(gamemodes.unknown);
	});
});
