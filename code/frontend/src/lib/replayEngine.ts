export interface PieceData {
	type: string;
	rank: string;
	ownerId: number;
}

export interface HistoricalMove {
	moveIndex: number;
	playerId: number;
	fromX: number;
	fromY: number;
	toX: number;
	toY: number;
	attacker?: PieceData;
	defender?: PieceData;
	result: 'move' | 'win' | 'loss' | 'tie' | 'capture';
}

export interface GameHistory {
	gameId: string;
	initialState: (PieceData | null)[][];
	moves: HistoricalMove[];
	winnerId?: number | null;
	gameType?: string;
	player1Username?: string;
	player2Username?: string;
	winCause?: string;
	createdAt?: string;
}

/**
 * Reconstructs the board state at a specific move index.
 * @param history The full game history containing initial setup and moves.
 * @param targetMoveIndex The index (1-based) representing which state to return.
 *                        0 returns the initial layout.
 */
export function getBoardAtMove(
	history: GameHistory,
	targetMoveIndex: number
): (PieceData | null)[][] {
	// Deep clone initial setup (10x10)
	const board: (PieceData | null)[][] = history.initialState.map((row) =>
		row.map((cell) => (cell ? { ...cell } : null))
	);

	// Apply moves sequentially
	for (let i = 0; i < targetMoveIndex && i < history.moves.length; i++) {
		applyMove(board, history.moves[i]);
	}

	return board;
}

function applyMove(board: (PieceData | null)[][], move: HistoricalMove) {
	const { fromX, fromY, toX, toY, result } = move;
	const attacker = board[fromY][fromX];

	if (!attacker) return;

	// Movement: Clear the original square
	board[fromY][fromX] = null;

	switch (result) {
		case 'move':
		case 'win':
		case 'capture':
			board[toY][toX] = attacker;
			break;
		case 'loss':
			break;
		case 'tie':
			board[toY][toX] = null;
			break;
	}
}

export interface ValidationResult {
	valid: boolean;
	error?: string;
	data?: GameHistory;
}

export function validateReplayJson(data: unknown): ValidationResult {
	if (!data || typeof data !== 'object') {
		return { valid: false, error: 'Replay file is empty or not a valid JSON object' };
	}

	const candidate = data as Partial<GameHistory>;

	if (!Array.isArray(candidate.initialState) || candidate.initialState.length === 0) {
		return { valid: false, error: 'Replay missing valid initialState grid' };
	}

	if (!Array.isArray(candidate.moves)) {
		return { valid: false, error: 'Replay missing valid moves list' };
	}

	for (let i = 0; i < candidate.moves.length; i++) {
		const m = candidate.moves[i];
		if (
			!m ||
			typeof m.fromX !== 'number' ||
			typeof m.fromY !== 'number' ||
			typeof m.toX !== 'number' ||
			typeof m.toY !== 'number' ||
			typeof m.playerId !== 'number'
		) {
			return { valid: false, error: `Move #${i + 1} contains invalid coordinates or player ID` };
		}
	}

	return {
		valid: true,
		data: {
			gameId: typeof candidate.gameId === 'string' ? candidate.gameId : 'local-replay',
			initialState: candidate.initialState,
			moves: candidate.moves,
			winnerId: candidate.winnerId ?? null,
			gameType: typeof candidate.gameType === 'string' ? candidate.gameType : undefined,
			player1Username:
				typeof candidate.player1Username === 'string' ? candidate.player1Username : undefined,
			player2Username:
				typeof candidate.player2Username === 'string' ? candidate.player2Username : undefined,
			winCause: typeof candidate.winCause === 'string' ? candidate.winCause : undefined,
			createdAt: typeof candidate.createdAt === 'string' ? candidate.createdAt : undefined
		}
	};
}
