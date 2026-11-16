import { describe, expect, it } from 'vitest';
import {
	decodeSetup,
	encodeSetup,
	formatBoardSetupForExport,
	sanitizeFilename
} from './board-binary';

const rows = ['0B12345678', '9MBBBBB111', '2222233333', '4444555566'];
const encoded = rows.join('');

describe('board-binary', () => {
	// why: the round-trip and flip are covered by contracts.spec.ts; these are frontend-only edge cases
	describe('encodeSetup / decodeSetup', () => {
		it('pads short and missing rows with dots, truncates long rows', () => {
			const result = encodeSetup(['12', '12345678901234']);
			expect(result).toHaveLength(40);
			expect(result.slice(0, 10)).toBe('12........');
			expect(result.slice(10, 20)).toBe('1234567890');
			expect(result.slice(20)).toBe('.'.repeat(20));
		});

		it('returns an empty setup for invalid length', () => {
			expect(decodeSetup('abc')).toEqual(Array(4).fill('.'.repeat(10)));
		});
	});

	describe('formatBoardSetupForExport', () => {
		it('adds decoded rows and defaults description', () => {
			expect(formatBoardSetupForExport({ name: 'n', setup_data: encoded })).toEqual({
				name: 'n',
				description: '',
				setup_data: encoded,
				rows
			});
		});
	});

	describe('sanitizeFilename', () => {
		it('lowercases and replaces unsafe characters', () => {
			expect(sanitizeFilename('  My Setup/../v2!  ')).toBe('my_setup____v2_');
		});

		it('falls back when empty', () => {
			expect(sanitizeFilename('   ')).toBe('board_setup');
		});
	});
});
