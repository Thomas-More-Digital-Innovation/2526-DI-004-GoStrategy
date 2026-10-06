import type { Position, BoardState } from '$lib/types/game';
import type { GameSocket } from '$lib/api/websocket';

export function handlePieceCellClick(
	x: number,
	y: number,
	selected: Position | null,
	validMoves: Position[],
	socket: GameSocket,
	isReplaying: boolean,
	boardState: BoardState | null
): { newSelected: Position | null; newValidMoves: Position[] } {
	if (isReplaying) {
		return { newSelected: selected, newValidMoves: validMoves };
	}

	const board = boardState?.board;
	if (!board) return { newSelected: null, newValidMoves: [] };

	const clickedPiece = board[y][x];
	const hasValidPiece = clickedPiece && clickedPiece.ownerName;

	if (!selected) {
		if (hasValidPiece && clickedPiece.ownerId === 0) {
			socket.requestValidMoves({ x, y });
			return { newSelected: { x, y }, newValidMoves: [] };
		}
		return { newSelected: null, newValidMoves: [] };
	}

	if (selected.x === x && selected.y === y) {
		return { newSelected: null, newValidMoves: [] };
	}

	const isValid = validMoves.some((m) => m.x === x && m.y === y);
	if (isValid) {
		socket.sendMove(selected, { x, y });
		return { newSelected: null, newValidMoves: [] };
	}

	if (hasValidPiece && clickedPiece.ownerId === 0) {
		socket.requestValidMoves({ x, y });
		return { newSelected: { x, y }, newValidMoves: [] };
	}

	return { newSelected: null, newValidMoves: [] };
}

export function handleSetupPieceSwap(
	x: number,
	y: number,
	playerId: number,
	swapPos1: Position | null,
	socket: GameSocket
): Position | null {
	const startRow = playerId === 0 ? 6 : 0;
	const endRow = playerId === 0 ? 9 : 3;
	if (y < startRow || y > endRow) return swapPos1;

	if (!swapPos1) {
		return { x, y };
	}
	socket.sendSwapPieces(swapPos1, { x, y });
	return null;
}
