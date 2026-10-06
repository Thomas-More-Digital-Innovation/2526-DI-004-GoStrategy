import { getContext, setContext } from 'svelte';
import type { GameSessionController } from './game-session.svelte';

const GAME_SESSION_KEY = Symbol('GAME_SESSION');

export function setGameSession(session: GameSessionController): GameSessionController {
	return setContext(GAME_SESSION_KEY, session);
}

export function useGameSession(): GameSessionController {
	const session = getContext<GameSessionController>(GAME_SESSION_KEY);
	if (!session) {
		throw new Error('useGameSession must be used within a GameSession context');
	}
	return session;
}
