import type { GameMode } from '$lib/types/game';
import { gamemodes } from '$lib/data/gamemodes.data';

export function getDisabledRows(
	isSetup: boolean,
	mode: GameMode,
	selectedPlayer: number
): number[] {
	if (!isSetup) return [];
	if (mode.mode === gamemodes.human_vs_ai.mode) return [0, 1, 2, 3, 4, 5];
	return selectedPlayer === 0 ? [0, 1, 2, 3, 4, 5] : [4, 5, 6, 7, 8, 9];
}

export function getHighlightedRows(
	isSetup: boolean,
	mode: GameMode,
	selectedPlayer: number
): number[] {
	if (!isSetup) return [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
	if (mode.mode === gamemodes.human_vs_ai.mode) return [6, 7, 8, 9];
	return selectedPlayer === 0 ? [6, 7, 8, 9] : [0, 1, 2, 3];
}

export function getHighlightColor(
	isSetup: boolean,
	mode: GameMode,
	selectedPlayer: number,
	currentPlayerId?: number
): 'red' | 'blue' {
	if (isSetup) {
		if (mode.mode === gamemodes.human_vs_ai.mode) return 'red';
		return selectedPlayer === 0 ? 'red' : 'blue';
	}
	return currentPlayerId === 0 ? 'red' : 'blue';
}
