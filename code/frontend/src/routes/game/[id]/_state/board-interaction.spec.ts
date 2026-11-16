import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { GameSocket } from '#lib/api/websocket.js';
import type { BoardState } from '#lib/types/game.js';
import { handlePieceCellClick, handleSetupPieceSwap } from './board-interaction';

const socket = {
	requestValidMoves: vi.fn(),
	sendMove: vi.fn(),
	sendSwapPieces: vi.fn()
};
const fakeSocket = socket as unknown as GameSocket;

const own = { ownerId: 0, ownerName: 'me' };
const enemy = { ownerId: 1, ownerName: 'them' };
const boardState = {
	board: [
		[own, enemy],
		[null, own]
	]
} as unknown as BoardState;

function click(x: number, y: number, selected: { x: number; y: number } | null, valid = []) {
	return handlePieceCellClick(x, y, selected, valid, fakeSocket, false, boardState);
}

// why: legality of moves is the backend's job; these tests only cover selection and which socket call is sent
describe('handlePieceCellClick', () => {
	beforeEach(() => vi.clearAllMocks());

	it('does nothing while replaying', () => {
		const sel = { x: 0, y: 0 };
		const res = handlePieceCellClick(1, 1, sel, [sel], fakeSocket, true, boardState);
		expect(res).toEqual({ newSelected: sel, newValidMoves: [sel] });
		expect(socket.requestValidMoves).not.toHaveBeenCalled();
	});

	it('clears selection without a board', () => {
		const res = handlePieceCellClick(0, 0, null, [], fakeSocket, false, null);
		expect(res).toEqual({ newSelected: null, newValidMoves: [] });
	});

	it('selects an own piece and requests its moves', () => {
		expect(click(0, 0, null).newSelected).toEqual({ x: 0, y: 0 });
		expect(socket.requestValidMoves).toHaveBeenCalledWith({ x: 0, y: 0 });
	});

	it('ignores enemy pieces and empty squares when nothing is selected', () => {
		expect(click(1, 0, null).newSelected).toBeNull();
		expect(click(0, 1, null).newSelected).toBeNull();
		expect(socket.requestValidMoves).not.toHaveBeenCalled();
	});

	it('deselects when clicking the selected piece', () => {
		expect(click(0, 0, { x: 0, y: 0 }).newSelected).toBeNull();
	});

	it('sends a move to a valid target', () => {
		const res = handlePieceCellClick(
			1,
			0,
			{ x: 0, y: 0 },
			[{ x: 1, y: 0 }],
			fakeSocket,
			false,
			boardState
		);
		expect(res.newSelected).toBeNull();
		expect(socket.sendMove).toHaveBeenCalledWith({ x: 0, y: 0 }, { x: 1, y: 0 });
	});

	it('switches selection to another own piece', () => {
		expect(click(1, 1, { x: 0, y: 0 }).newSelected).toEqual({ x: 1, y: 1 });
		expect(socket.requestValidMoves).toHaveBeenCalledWith({ x: 1, y: 1 });
	});

	it('clears selection on an invalid target', () => {
		expect(click(1, 0, { x: 0, y: 0 }).newSelected).toBeNull();
		expect(socket.sendMove).not.toHaveBeenCalled();
	});
});

describe('handleSetupPieceSwap', () => {
	beforeEach(() => vi.clearAllMocks());

	it('ignores clicks outside the player rows', () => {
		expect(handleSetupPieceSwap(0, 5, 0, null, fakeSocket)).toBeNull();
		expect(handleSetupPieceSwap(0, 6, 1, { x: 1, y: 1 }, fakeSocket)).toEqual({ x: 1, y: 1 });
	});

	it('picks the first square, then swaps on the second', () => {
		const first = handleSetupPieceSwap(0, 6, 0, null, fakeSocket);
		expect(first).toEqual({ x: 0, y: 6 });
		expect(handleSetupPieceSwap(1, 9, 0, first, fakeSocket)).toBeNull();
		expect(socket.sendSwapPieces).toHaveBeenCalledWith({ x: 0, y: 6 }, { x: 1, y: 9 });
	});

	it('uses the top rows for player 1', () => {
		expect(handleSetupPieceSwap(0, 3, 1, null, fakeSocket)).toEqual({ x: 0, y: 3 });
	});
});
