<script lang="ts">
	import type { Piece as PieceType } from '$lib/types/game';
	import { PIECE_INVENTORY, type PieceInfo } from '$lib/types/board-setup';

	interface Props {
		piece: PieceType | null;
		isSelected?: boolean;
		isHighlighted?: boolean;
		isLake?: boolean;
		viewerId?: number;
		scale?: number;
	}

	let {
		piece,
		isSelected = false,
		isHighlighted = false,
		isLake = false,
		viewerId = 0,
		scale = 1
	}: Props = $props();

	const isOccupied = $derived(
		Boolean(
			piece &&
			piece.ownerId !== undefined &&
			piece.ownerId >= 0 &&
			(piece.ownerName || piece.rank || piece.type)
		)
	);

	const canSeePiece = $derived(
		Boolean(isOccupied && piece && (piece.ownerId === viewerId || piece.revealed))
	);

	const isEnemy = $derived(Boolean(isOccupied && piece && piece.ownerId !== viewerId));

	const pieceRank = $derived.by(() => {
		if (!piece || !piece.rank) return null;
		return piece.rank;
	});

	const pieceIcon = $derived.by(() => {
		if (!piece) return null;

		const directIcon = piece.ownerId === 1 ? piece.iconBlue : piece.iconRed;
		if (directIcon) return directIcon;

		if (!piece.rank) return null;

		let inventoryItem: PieceInfo | undefined = PIECE_INVENTORY[piece.rank];
		if (!inventoryItem) {
			inventoryItem = Object.values(PIECE_INVENTORY).find((item) => item.rank === piece.rank);
		}

		if (!inventoryItem) return null;

		return piece.ownerId === 1 ? inventoryItem.icon_blue : inventoryItem.icon_red;
	});

	const isAsset = $derived.by(() => {
		const icon = pieceIcon;
		return Boolean(icon && (icon.includes('/') || icon.includes('.')));
	});

	const fontSize = $derived(18 * scale);
</script>

<div
	class="piece unselectable"
	class:selected={isSelected}
	class:highlighted={isHighlighted}
	class:empty={!isOccupied && !isLake}
	class:lake={isLake}
	class:player1={isOccupied && piece && piece.ownerId === 1}
	class:player2={isOccupied && piece && (piece.ownerId === 2 || piece.ownerId === 0)}
	style="--scale: {scale}"
>
	{#if isLake}
		<span style="font-size: {24 * scale}px">🌊</span>
	{:else if isOccupied}
		{#if canSeePiece}
			<div class="flex flex-col items-center justify-center w-full h-full overflow-hidden relative">
				{#if isAsset}
					<img
						src={pieceIcon}
						alt={pieceRank}
						class="w-full h-full object-fill pointer-events-none"
					/>
					{#if pieceRank === '0'}
						<span class="piece-rank border-2 border-green-400">F</span>
					{:else}
						<span class="piece-rank">{pieceRank}</span>
					{/if}
				{:else if pieceRank}
					<span class="font-bold text-white" style="font-size: {fontSize}px">
						{pieceRank}
					</span>
				{/if}
			</div>
		{:else if isEnemy}
			<span style="font-size: {20 * scale}px; font-weight: bold; color: rgba(255,255,255,0.8)"
				>?</span
			>
		{/if}
	{/if}
</div>

<style>
	.piece {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.15s;
		border-radius: calc(6px * var(--scale));

		img {
			border-radius: calc(6px * var(--scale));
		}

		.piece-rank {
			position: absolute;
			bottom: calc(2px * var(--scale));
			right: calc(2px * var(--scale));
			width: calc(28px * var(--scale));
			height: calc(28px * var(--scale));
			background: rgba(0, 0, 0, 0.4);
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: calc(14px * var(--scale));
			font-weight: bold;
			color: white;
			transition: all 0.15s;
			z-index: 2;
		}

		&:hover .piece-rank {
			bottom: 0;
			right: 0;
			width: 100%;
			height: 100%;
			border-radius: calc(6px * var(--scale));
			background: rgba(0, 0, 0, 0.6);
			font-size: calc(20px * var(--scale));
		}
	}

	.piece.empty {
		border: 1px solid rgba(255, 255, 255, 0.05);
		background: rgba(255, 255, 255, 0.05);
		cursor: default;
	}

	.piece.lake {
		background: oklch(0.6 0.15 240);
		opacity: 0.4;
		cursor: default;
	}

	.piece.player1 {
		background: oklch(0.6 0.2 260);
		cursor: pointer;
	}

	.piece.player2 {
		background: oklch(0.6 0.2 20);
		cursor: pointer;
	}

	.piece.player1:hover,
	.piece.player2:hover {
		transform: scale(1.05);
		z-index: 1;
	}

	.piece.selected {
		border: 2px solid oklch(0.7 0.2 150);
		box-shadow: 0 0 12px oklch(0.7 0.2 150 / 0.4);
	}

	.piece.highlighted {
		border: 2px solid oklch(0.75 0.18 145);
	}
</style>
