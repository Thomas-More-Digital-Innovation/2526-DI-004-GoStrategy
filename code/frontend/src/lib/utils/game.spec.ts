import { describe, expect, it } from 'vitest';
import { PIECE_INVENTORY } from '#lib/types/board-setup.js';
import { getPieceIcon } from './game';

describe('getPieceIcon', () => {
	it('returns null for missing or hidden pieces', () => {
		expect(getPieceIcon(undefined)).toBeNull();
		expect(getPieceIcon({ rank: '', ownerId: 0 } as never)).toBeNull();
	});

	it('looks up by inventory key and picks colour by owner', () => {
		expect(getPieceIcon({ rank: 'B', ownerId: 0 } as never)).toBe(PIECE_INVENTORY.B.icon_red);
		expect(getPieceIcon({ rank: 'B', ownerId: 1 } as never)).toBe(PIECE_INVENTORY.B.icon_blue);
	});

	it('falls back to lookup by display rank', () => {
		expect(getPieceIcon({ rank: 'F', ownerId: 0 } as never)).toBe(PIECE_INVENTORY['0'].icon_red);
	});

	it('returns null for unknown ranks', () => {
		expect(getPieceIcon({ rank: 'Z', ownerId: 0 } as never)).toBeNull();
	});
});
