<script lang="ts">
	import type {
		BoardState,
		Position,
		Piece as PieceType,
		HistoricalMove
	} from '#lib/types/game.js';
	import MoveVisualization from './move-visualization/MoveVisualization.svelte';
	import BoardCell from './BoardCell.svelte';
	import { BOARD_CONFIG } from '#lib/data/board.data.js';
	import { a11yStore, getCellAriaLabel, isLake, announceMove } from '#lib/utils/a11y.svelte.js';

	interface Props {
		boardState?: BoardState | null;
		board?: (PieceType | null)[][]; // Allow passing raw board array
		selectedPosition?: Position | null;
		onCellClick?: (x: number, y: number) => void;
		isInteractive?: boolean;
		viewerId?: number;
		disabledRows?: number[];
		validMoves?: Position[];
		rows?: number;
		cols?: number;
		scale?: number; // Scale factor (e.g., 0.5 for small preview)
		onCellDragStart?: (e: DragEvent, x: number, y: number) => void;
		onCellDragOver?: (e: DragEvent, x: number, y: number) => void;
		onCellDragLeave?: (e: DragEvent, x: number, y: number) => void;
		onCellDrop?: (e: DragEvent, x: number, y: number) => void;
		isLakeCell?: (x: number, y: number) => boolean;
		responsive?: boolean;
		visualDisabledRows?: number[];
		lastMove?: HistoricalMove | null;
		highlightedRows?: number[];
		highlightColor?: 'red' | 'blue' | '';
		isSetupPhase?: boolean;
		currentPlayerId?: number;
		isGameOver?: boolean;
		winnerId?: number | null;
	}

	let {
		boardState,
		board,
		selectedPosition,
		onCellClick,
		isInteractive = true,
		viewerId = 0,
		disabledRows = [],
		validMoves = [],
		rows = BOARD_CONFIG.rows,
		cols = BOARD_CONFIG.cols,
		scale = 1,
		onCellDragStart,
		onCellDragOver,
		onCellDragLeave,
		onCellDrop,
		isLakeCell,
		responsive = false,
		visualDisabledRows = [],
		lastMove = null,
		highlightedRows = [],
		highlightColor = '',
		isSetupPhase = false,
		currentPlayerId,
		isGameOver = false,
		winnerId = null
	}: Props = $props();

	const displayBoard = $derived(board || boardState?.board || []);

	const isSelected = (x: number, y: number) =>
		selectedPosition?.x === x && selectedPosition?.y === y;
	const isValidMove = (x: number, y: number) => validMoves.some((m) => m.x === x && m.y === y);
	const isCellFocused = (x: number, y: number) =>
		(a11yStore.focusedCell?.x === x && a11yStore.focusedCell?.y === y) ||
		(a11yStore.focusedCell === null && x === 0 && y === 0);
	const isCellDisabled = (x: number, y: number) =>
		!isInteractive || isLake(x, y, rows, isLakeCell) || disabledRows.includes(y);

	const cellSize = $derived(BOARD_CONFIG.baseCellSize * scale);

	const highlightStyle = $derived.by(() => {
		if (highlightedRows.length === 0) return '';
		const minRow = Math.min(...highlightedRows);
		const maxRow = Math.max(...highlightedRows);
		const rowCount = maxRow - minRow + 1;
		const gap = BOARD_CONFIG.gap * scale;

		const marginY = 12;
		const top = minRow * (cellSize + gap);
		const height = rowCount * cellSize + (rowCount - 1) * gap + marginY;

		return `top: ${top}px; height: ${height}px;`;
	});

	$effect(() => {
		if (currentPlayerId !== undefined && !isGameOver && !isSetupPhase) {
			const turnName = currentPlayerId === 0 ? 'Your turn (Blue)' : "Opponent's turn (Red)";
			a11yStore.announce(turnName);
		}
	});

	$effect(() => {
		if (lastMove) {
			const msg = announceMove(lastMove, viewerId);
			if (msg) {
				a11yStore.announce(msg);
			}
		}
	});

	$effect(() => {
		if (isGameOver) {
			const msg =
				winnerId === 0
					? 'Game over! You won!'
					: winnerId === 1
						? 'Game over! Red won!'
						: 'Game over!';
			a11yStore.announce(msg);
		}
	});
</script>

<div
	class="board-wrapper"
	style="--cell-size: {cellSize}px; --cols: {cols}; --rows: {rows}; --scale: {scale}; --gap: {BOARD_CONFIG.gap}px;"
	class:non-interactive={!isInteractive}
>
	<div class="sr-only" aria-live="polite" aria-atomic="true">
		{a11yStore.announcement}
	</div>

	{#if displayBoard.length > 0}
		<div
			class="board glass rounded-2xl p-3"
			class:responsive
			role="grid"
			aria-label="Stratego Game Board"
		>
			{#each { length: rows }, y (y)}
				<div role="row" style="display: contents;">
					{#each { length: cols }, x (x)}
						{@const piece = displayBoard[y]?.[x]}
						{@const lake = isLake(x, y, rows, isLakeCell)}
						<BoardCell
							{x}
							{y}
							{piece}
							isSelected={isSelected(x, y)}
							isValidMove={isValidMove(x, y)}
							isLake={lake}
							interactive={isInteractive && !lake}
							isFocused={isCellFocused(x, y)}
							isDisabled={isCellDisabled(x, y)}
							isVisualDisabled={visualDisabledRows.includes(y)}
							{responsive}
							draggable={isInteractive &&
								!!piece &&
								piece.ownerId !== undefined &&
								piece.ownerId >= 0 &&
								!lake &&
								!disabledRows.includes(y)}
							{viewerId}
							{scale}
							{cols}
							{rows}
							ariaLabel={getCellAriaLabel({
								x,
								y,
								piece,
								isSelected: isSelected(x, y),
								isValidMove: isValidMove(x, y),
								isLake: lake,
								viewerId,
								isSetupPhase
							})}
							{onCellClick}
							{onCellDragStart}
							{onCellDragOver}
							{onCellDragLeave}
							{onCellDrop}
						/>
					{/each}
				</div>
			{/each}

			{#if lastMove && !selectedPosition}
				<MoveVisualization move={lastMove} {cellSize} {scale} />
			{/if}

			{#if highlightedRows.length > 0}
				<div
					class="row-highlight"
					class:red={highlightColor === 'red'}
					class:blue={highlightColor === 'blue'}
					style={highlightStyle}
				></div>
			{/if}
		</div>
	{:else}
		<div class="flex items-center justify-center p-10 text-white/40">No board data available</div>
	{/if}
</div>

<style>
	.board-wrapper {
		display: flex;
		justify-content: center;
		align-items: center;
	}
	.board-wrapper.non-interactive {
		pointer-events: none;
	}
	.board {
		position: relative;
		display: grid;
		grid-template-columns: repeat(var(--cols), var(--cell-size));
		grid-template-rows: repeat(var(--rows), var(--cell-size));
		gap: calc(var(--gap) * var(--scale, 1));
	}
	.board.responsive {
		grid-template-columns: repeat(var(--cols), 1fr);
		grid-template-rows: auto;
		width: 100%;
		gap: 2px;
	}
	.row-highlight {
		position: absolute;
		left: 6px;
		right: 6px;
		margin-top: 6px;
		border: 2px solid white;
		border-radius: 10px;
		pointer-events: none;
		z-index: 5;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		opacity: 0;
		animation: fade-in 0.3s forwards;
	}
	@keyframes fade-in {
		from {
			opacity: 0;
			transform: scale(0.98);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}
	.row-highlight.red {
		border-color: oklch(0.6688 0.1971 39.15 / 0.5);
		background: oklch(0.6688 0.1971 39.15 / 0.05);
		box-shadow:
			0 0 30px oklch(0.6688 0.1971 39.15 / 0.15),
			inset 0 0 20px oklch(0.6688 0.1971 39.15 / 0.05);
	}
	.row-highlight.blue {
		border-color: oklch(0.7113 0.1044 226.48 / 0.5);
		background: oklch(0.7113 0.1044 226.48 / 0.05);
		box-shadow:
			0 0 30px oklch(0.7113 0.1044 226.48 / 0.15),
			inset 0 0 20px oklch(0.7113 0.1044 226.48 / 0.05);
	}
</style>
