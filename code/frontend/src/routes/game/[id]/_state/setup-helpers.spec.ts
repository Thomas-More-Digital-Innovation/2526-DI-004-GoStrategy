import { describe, expect, it } from 'vitest';
import { gamemodes } from '#lib/data/gamemodes.data.js';
import { getDisabledRows, getHighlightColor, getHighlightedRows } from './setup-helpers';

const vsAi = gamemodes.human_vs_ai;
const local = gamemodes.ai_vs_ai;
const own = [6, 7, 8, 9];
const top = [0, 1, 2, 3];

describe('setup-helpers', () => {
	it('disables nothing and highlights everything outside setup', () => {
		expect(getDisabledRows(false, vsAi, 0)).toEqual([]);
		expect(getHighlightedRows(false, vsAi, 0)).toEqual([...top, 4, 5, ...own]);
	});

	it('limits a human vs ai player to the bottom rows', () => {
		expect(getDisabledRows(true, vsAi, 1)).toEqual([...top, 4, 5]);
		expect(getHighlightedRows(true, vsAi, 1)).toEqual(own);
		expect(getHighlightColor(true, vsAi, 1)).toBe('red');
	});

	it('follows the selected player in other modes', () => {
		expect(getDisabledRows(true, local, 0)).toEqual([...top, 4, 5]);
		expect(getDisabledRows(true, local, 1)).toEqual([4, 5, ...own]);
		expect(getHighlightColor(true, local, 0)).toBe('red');
		expect(getHighlightColor(true, local, 1)).toBe('blue');
	});

	it('colours by current player during play', () => {
		expect(getHighlightColor(false, vsAi, 0, 0)).toBe('red');
		expect(getHighlightColor(false, vsAi, 0, 1)).toBe('blue');
	});
});
