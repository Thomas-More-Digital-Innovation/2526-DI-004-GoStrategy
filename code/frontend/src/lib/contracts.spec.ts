import { describe, expect, it } from 'vitest';
import auth from '../../../contracts/auth.json';
import setup from '../../../contracts/setup.json';
import board from '../../../contracts/board.json';
import type { GameSocket } from '#lib/api/websocket.js';
import { gamemodes } from '#lib/data/gamemodes.data.js';
import { PIECE_INVENTORY } from '#lib/types/board-setup.js';
import { isLake } from '#lib/utils/a11y.svelte.js';
import {
	isPasswordStrong,
	isUsernameFormatValid,
	isUsernameLengthValid,
	validatePasswordInput
} from '#lib/utils/authValidation.js';
import { decodeSetup, encodeSetup, flipSetup } from '#lib/utils/board-binary.js';
import { getHighlightedRows } from '../routes/game/[id]/_state/setup-helpers';
import { handleSetupPieceSwap } from '../routes/game/[id]/_state/board-interaction';

describe('backend contracts', () => {
	describe('auth', () => {
		it.each(auth.username)('username $input', ({ input, formatValid, valid }) => {
			expect(isUsernameFormatValid(input)).toBe(formatValid);
			expect(isUsernameFormatValid(input) && isUsernameLengthValid(input)).toBe(valid);
		});

		it.each(auth.password)('password $input', ({ input, formatValid, accepted }) => {
			expect(validatePasswordInput(input).isValidFormat).toBe(formatValid);
			expect(isPasswordStrong(input)).toBe(accepted);
		});
	});

	describe('setup', () => {
		it('piece counts match the inventory', () => {
			const counts = Object.fromEntries(
				Object.entries(PIECE_INVENTORY).map(([rank, info]) => [rank, info.count])
			);
			expect(counts).toEqual(setup.pieceCounts);
		});

		it.each(setup.cases)('encodes, decodes and flips $setupData', (tc) => {
			expect(encodeSetup(tc.rows)).toBe(tc.setupData);
			expect(decodeSetup(tc.setupData)).toEqual(tc.rows);
			expect(flipSetup(tc.setupData)).toBe(tc.flippedSetupData);
		});
	});

	describe('board', () => {
		it('lakes', () => {
			const lakes = [];
			for (let y = 0; y < board.size; y++) {
				for (let x = 0; x < board.size; x++) {
					if (isLake(x, y)) lakes.push([x, y]);
				}
			}
			expect(lakes).toEqual(board.lakes);
		});

		it.each(Object.entries(board.setupRows))('setup rows for player %s', (player, [start, end]) => {
			const id = Number(player);
			const range = Array.from({ length: end - start + 1 }, (_, i) => start + i);
			const socket = { sendSwapPieces: () => {} } as unknown as GameSocket;
			const accepted = Array.from({ length: board.size }, (_, y) => y).filter(
				(y) => handleSetupPieceSwap(0, y, id, null, socket) !== null
			);
			expect(accepted).toEqual(range);
			expect(getHighlightedRows(true, gamemodes.ai_vs_ai, id)).toEqual(range);
		});
	});
});
