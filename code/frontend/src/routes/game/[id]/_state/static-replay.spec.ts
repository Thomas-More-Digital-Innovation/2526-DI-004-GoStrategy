import { beforeEach, describe, expect, it, vi } from 'vitest';
import { goto } from '$app/navigation';
import { games } from '#lib/api/client.js';
import { gamemodes } from '#lib/data/gamemodes.data.js';
import { toastStore } from '#lib/state/toast.svelte.js';
import { GameConnectionManager } from './game-connection.svelte';
import { GameSessionController } from './game-session.svelte';
import { initStaticReplay } from './static-replay';
import { createFakeSocket } from './fake-socket';

vi.mock('$app/navigation', () => ({ goto: vi.fn() }));
vi.mock('$app/paths', () => ({ resolve: (p: string) => `/${p}` }));
vi.mock('#lib/api/client.js', () => ({ games: { getHistory: vi.fn() } }));

const history = {
	gameId: 'g',
	initialState: [[{ type: 'Scout', rank: '2', ownerId: 0 }, null]],
	moves: [{ moveIndex: 0, playerId: 1, fromX: 0, fromY: 0, toX: 1, toY: 0, result: 'move' }],
	winnerId: 1,
	gameType: 'human_vs_ai',
	player2Username: 'bob'
};

describe('initStaticReplay', () => {
	let session: GameSessionController;
	const store = new Map<string, string>();

	beforeEach(() => {
		vi.clearAllMocks();
		store.clear();
		toastStore.toasts = [];
		vi.stubGlobal('sessionStorage', { getItem: (k: string) => store.get(k) ?? null });
		session = new GameSessionController(new GameConnectionManager(createFakeSocket()));
	});

	it('loads a server replay and jumps to the last move', async () => {
		vi.mocked(games.getHistory).mockResolvedValue(history as never);
		await initStaticReplay(session, 'g');
		expect(session.isStaticReplay).toBe(true);
		expect(session.store.gameMode).toBe(gamemodes.human_vs_ai);
		expect(session.store.gameState).toMatchObject({
			isGameOver: true,
			winnerId: 1,
			winCause: 'Victory achieved',
			currentPlayerId: 1,
			player1Username: 'Player 1',
			player2Username: 'bob'
		});
		expect(session.store.isReplaying).toBe(true);
		expect(session.store.currentHistoryIndex).toBe(0);
	});

	it('loads a local replay from sessionStorage', async () => {
		store.set(
			'gostrategy_local_replay',
			JSON.stringify({ ...history, winnerId: null, gameType: undefined })
		);
		await initStaticReplay(session, 'local');
		expect(games.getHistory).not.toHaveBeenCalled();
		expect(session.store.gameMode).toBe(gamemodes.ai_vs_ai);
		expect(session.store.gameState).toMatchObject({ winnerName: 'Draw', winCause: 'Draw' });
	});

	it('redirects to profile when no local replay exists', async () => {
		await initStaticReplay(session, 'local');
		expect(toastStore.toasts[0].type).toBe('warning');
		expect(goto).toHaveBeenCalledWith('/profile', { replaceState: true });
	});

	it('redirects to the live game when loading fails', async () => {
		vi.mocked(games.getHistory).mockRejectedValue(new Error('gone'));
		await initStaticReplay(session, 'g');
		expect(toastStore.toasts[0].message).toBe('gone');
		expect(goto).toHaveBeenCalledWith('/game/g', { replaceState: true });
	});
});
