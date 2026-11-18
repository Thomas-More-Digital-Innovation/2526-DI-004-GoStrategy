import type { GameHistory } from './generated/game';
import type {
	BoardStateMessage,
	GameStateMessage,
	MoveDTO,
	PieceDTO,
	PositionDTO
} from './generated/dto';
import type {
	CreateGameResponse,
	Game,
	UserDTO,
	UserStats as GeneratedUserStats
} from './generated/models';

export type { HistoricalMove, MoveResultType, PieceData } from './generated/game';
import type { HistoricalMove } from './generated/game';

export type Position = PositionDTO;
export type Move = MoveDTO;
export type GameState = GameStateMessage;
export type User = UserDTO;
export type UserStats = GeneratedUserStats;
export type GameRecord = Game;
export type GameInfo = CreateGameResponse;
// why: static replays are assembled client-side from a game record, so they carry metadata the backend history lacks
export interface GameHistoryResponse extends GameHistory {
	gameType?: string;
	player1Username?: string;
	player2Username?: string;
	winCause?: string;
	createdAt?: string;
}

export interface Piece extends PieceDTO {
	iconBlue?: string;
	iconRed?: string;
}

export interface BoardState extends Omit<BoardStateMessage, 'board'> {
	board: Piece[][];
}

export type MoveVisualizationHighlightState = 'move' | 'win' | 'loss';

export interface HistoryMove {
	moveNumber: number;
	move: HistoricalMove;
	piece?: Piece;
	boardState: Piece[][];
}

export interface CombatAnimation {
	attacker: Piece;
	defender: Piece;
	attackerWon: boolean;
	defenderWon: boolean;
}

export interface GameMode {
	mode: string;
	icon: string;
	title: string;
	desc: string;
	disabled?: boolean;
}

export interface AI {
	name: string;
	id: string;
	description: string;
	image?: string;
}
