import type { PieceData } from '$lib/replayEngine';
import { PIECE_INVENTORY, type PieceInfo } from '$lib/types/board-setup';

export const getPieceIcon = (piece: PieceData | undefined) => {
	if (!piece || !piece.rank || piece.rank === '') return null;
	let inventoryItem: PieceInfo | undefined = PIECE_INVENTORY[piece.rank];
	if (!inventoryItem) {
		inventoryItem = Object.values(PIECE_INVENTORY).find((item) => item.rank === piece.rank);
	}
	if (!inventoryItem) return null;
	return piece.ownerId === 0 ? inventoryItem.icon_red : inventoryItem.icon_blue;
};
