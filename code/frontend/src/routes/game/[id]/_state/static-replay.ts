import { goto } from '$app/navigation';
import { games } from '#lib/api/client.js';
import { toastStore } from '#lib/state/toast.svelte.js';
import { gamemodes } from '#lib/data/gamemodes.data.js';
import type { GameHistoryResponse } from '#lib/types/game.js';
import type { GameSessionController } from './game-session.svelte';
import { resolve } from '$app/paths';

export async function initStaticReplay(
	session: GameSessionController,
	gameId: string
): Promise<void> {
	session.isStaticReplay = true;
	session.store.gameMode = gamemodes.ai_vs_ai;

	try {
		let history: GameHistoryResponse;

		if (gameId === 'local') {
			const stored = sessionStorage.getItem('gostrategy_local_replay');
			if (!stored) {
				toastStore.warning('No local replay found in session');
				goto(resolve('profile'), { replaceState: true });
				return;
			}
			const parsed = JSON.parse(stored);
			history = {
				gameId: parsed.gameId || 'local',
				initialState: parsed.initialState,
				moves: parsed.moves,
				winnerId: parsed.winnerId ?? null,
				gameType: parsed.gameType,
				player1Username: parsed.player1Username,
				player2Username: parsed.player2Username,
				winCause: parsed.winCause,
				createdAt: parsed.createdAt
			};
		} else {
			history = await games.getHistory(gameId);
		}

		if (history.gameType) {
			session.store.gameMode = gamemodes.fromString(history.gameType);
		} else {
			session.store.gameMode = gamemodes.ai_vs_ai;
		}

		session.store.loadMoveHistory(
			{
				initialState: history.initialState,
				fullHistory: history.moves,
				moves: history.moves.map((m) => ({
					from: { x: m.fromX, y: m.fromY },
					to: { x: m.toX, y: m.toY }
				}))
			},
			gameId,
			-1
		);

		const p1Name = history.player1Username || 'Player 1';
		const p2Name = history.player2Username || 'Player 2';
		const winnerName =
			history.winnerId === 0
				? `${p1Name} (Blue)`
				: history.winnerId === 1
					? `${p2Name} (Red)`
					: 'Draw';

		session.store.gameState = {
			round: Math.ceil(history.moves.length / 2) || 1,
			currentPlayerId:
				history.moves.length > 0 ? history.moves[history.moves.length - 1].playerId : 0,
			currentPlayerName: '',
			isGameOver: true,
			winnerId:
				history.winnerId !== null && history.winnerId !== undefined ? history.winnerId : undefined,
			winnerName,
			winCause:
				history.winCause ||
				(history.winnerId !== null && history.winnerId !== undefined ? 'Victory achieved' : 'Draw'),
			player1Score: 0,
			player2Score: 0,
			waitingForInput: false,
			paused: false,
			moveCount: history.moves.length,
			player1AlivePieces: 0,
			player2AlivePieces: 0,
			isSetupPhase: false,
			headless: false,
			player1Username: p1Name,
			player2Username: p2Name
		};

		session.store.isReplaying = true;

		if (session.store.history.length > 0) {
			session.store.goToMove(session.store.history.length - 1);
		}
	} catch (err) {
		toastStore.handleApiMessage(err, 'Failed to load game replay');
		goto(resolve(gameId === 'local' ? 'profile' : `game/${gameId}`), { replaceState: true });
	}
}
