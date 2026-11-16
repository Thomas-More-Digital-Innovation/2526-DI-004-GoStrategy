import { describe, expect, it } from 'vitest';
import { getBoardAtMove, validateReplayJson, type GameHistory } from './replayEngine';

const red = { type: 'Scout', rank: '2', ownerId: 0 };
const blue = { type: 'Bomb', rank: 'B', ownerId: 1 };

function history(result: GameHistory['moves'][number]['result']): GameHistory {
	return {
		gameId: 'g',
		initialState: [
			[red, blue],
			[null, null]
		],
		moves: [{ moveIndex: 0, playerId: 0, fromX: 0, fromY: 0, toX: 1, toY: 0, result }]
	};
}

describe('replayEngine', () => {
	describe('getBoardAtMove', () => {
		it('returns a copy of the initial state at index 0', () => {
			const h = history('move');
			const board = getBoardAtMove(h, 0);
			expect(board).toEqual(h.initialState);
			expect(board[0][0]).not.toBe(h.initialState[0][0]);
		});

		it.each([
			['move', [null, red]],
			['win', [null, red]],
			['capture', [null, red]],
			['loss', [null, blue]],
			['tie', [null, null]]
		] as const)('applies a %s result', (result, row) => {
			expect(getBoardAtMove(history(result), 1)[0]).toEqual(row);
		});

		it('ignores indices past the last move', () => {
			expect(getBoardAtMove(history('move'), 99)[0]).toEqual([null, red]);
		});

		it('skips moves from an empty square', () => {
			const h = history('move');
			h.moves[0] = { ...h.moves[0], fromX: 0, fromY: 1 };
			expect(getBoardAtMove(h, 1)).toEqual(h.initialState);
		});
	});

	describe('validateReplayJson', () => {
		it.each([
			[null, 'empty'],
			['str', 'empty'],
			[{ initialState: [] }, 'initialState'],
			[{ initialState: [[]] }, 'moves list'],
			[{ initialState: [[]], moves: [{ fromX: 0 }] }, 'Move #1']
		])('rejects %j', (data, error) => {
			const result = validateReplayJson(data);
			expect(result.valid).toBe(false);
			expect(result.error).toContain(error);
		});

		it('fills defaults for optional fields', () => {
			const result = validateReplayJson({ initialState: [[]], moves: [], gameType: 5 });
			expect(result).toEqual({
				valid: true,
				data: {
					gameId: 'local-replay',
					initialState: [[]],
					moves: [],
					winnerId: null,
					gameType: undefined,
					player1Username: undefined,
					player2Username: undefined,
					winCause: undefined,
					createdAt: undefined
				}
			});
		});

		it('keeps valid optional fields', () => {
			const h = { ...history('win'), winnerId: 0, player1Username: 'a', createdAt: 'now' };
			const result = validateReplayJson(h);
			expect(result.data).toMatchObject({
				gameId: 'g',
				winnerId: 0,
				player1Username: 'a',
				createdAt: 'now'
			});
		});
	});
});
