<script lang="ts">
	import Piece from './Piece.svelte';
	import type { Piece as PieceType } from '#lib/types/game.js';
	import { handleBoardKeyDown, a11yStore } from '#lib/utils/a11y.svelte.js';

	interface Props {
		x: number;
		y: number;
		piece: PieceType | null;
		isSelected: boolean;
		isValidMove: boolean;
		isLake: boolean;
		interactive: boolean;
		isFocused: boolean;
		isDisabled: boolean;
		isVisualDisabled: boolean;
		responsive: boolean;
		draggable: boolean;
		viewerId: number;
		scale: number;
		ariaLabel: string;
		cols: number;
		rows: number;
		onCellClick?: (x: number, y: number) => void;
		onCellDragStart?: (e: DragEvent, x: number, y: number) => void;
		onCellDragOver?: (e: DragEvent, x: number, y: number) => void;
		onCellDragLeave?: (e: DragEvent, x: number, y: number) => void;
		onCellDrop?: (e: DragEvent, x: number, y: number) => void;
	}

	let {
		x,
		y,
		piece,
		isSelected,
		isValidMove,
		isLake,
		interactive,
		isFocused,
		isDisabled,
		isVisualDisabled,
		responsive,
		draggable,
		viewerId,
		scale,
		ariaLabel,
		cols,
		rows,
		onCellClick,
		onCellDragStart,
		onCellDragOver,
		onCellDragLeave,
		onCellDrop
	}: Props = $props();
</script>

<button
	id="cell-{x}-{y}"
	class="cell focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-accent focus-visible:ring-offset-2 focus-visible:ring-offset-black/50"
	class:lake={isLake}
	class:interactive
	class:valid-move={isValidMove}
	class:visual-disabled={isVisualDisabled}
	class:responsive
	role="gridcell"
	tabindex={isFocused ? 0 : -1}
	aria-disabled={isDisabled ? 'true' : undefined}
	aria-label={ariaLabel}
	{draggable}
	onclick={() => {
		if (isDisabled) return;
		onCellClick?.(x, y);
	}}
	onkeydown={(e) => handleBoardKeyDown(e, { x, y, cols, rows })}
	onfocus={() => {
		a11yStore.focusedCell = { x, y };
	}}
	ondragstart={(e) => {
		if (isDisabled) return;
		onCellDragStart?.(e, x, y);
	}}
	ondragover={(e) => {
		if (isDisabled) return;
		if (onCellDrop) {
			e.preventDefault();
			onCellDragOver?.(e, x, y);
		}
	}}
	ondragleave={(e) => {
		if (isDisabled) return;
		onCellDragLeave?.(e, x, y);
	}}
	ondrop={(e) => {
		if (isDisabled) return;
		onCellDrop?.(e, x, y);
	}}
>
	<Piece {piece} {isSelected} isHighlighted={isValidMove} {isLake} {viewerId} {scale} />
</button>

<style>
	.cell {
		width: var(--cell-size);
		height: var(--cell-size);
		position: relative;
		padding: 0;
		border: none;
		border-radius: 6px;
		background: none;
		cursor: default;
	}
	.cell.responsive {
		width: 100%;
		height: auto;
		aspect-ratio: 1;
	}
	.cell.interactive {
		cursor: pointer;
	}
	.cell.visual-disabled {
		background-color: rgba(255, 0, 0, 0.1);
	}
	.cell.visual-disabled:not(.lake)::before {
		content: '';
		position: absolute;
		inset: 0;
		background: linear-gradient(
			to top right,
			transparent calc(50% - 1px),
			rgba(255, 255, 255, 0.2) 50%,
			transparent calc(50% + 1px)
		);
		border-radius: 6px;
		pointer-events: none;
		z-index: 1;
	}
	.cell.valid-move::after {
		content: '';
		position: absolute;
		inset: 3px;
		border: 2px solid oklch(0.75 0.18 145);
		border-radius: 4px;
		pointer-events: none;
		animation: pulse 1.5s ease-in-out infinite;
	}
	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.4;
		}
	}
</style>
