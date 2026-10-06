<script lang="ts">
	import { PIECE_INVENTORY } from '$lib/types/board-setup';
	import Card from '../ui/Card.svelte';

	interface Props {
		remainingCounts: Record<string, number>;
		selectedPieceRank: string | null;
		ownerId?: number;
		isComplete: boolean;
		onSelectRank: (rank: string) => void;
		onDragStart: (e: DragEvent, rank: string) => void;
	}

	let {
		remainingCounts,
		selectedPieceRank,
		ownerId = 1,
		isComplete,
		onSelectRank,
		onDragStart
	}: Props = $props();
</script>

<Card class="w-full lg:w-80 h-fit sticky top-6">
	<h3 class="font-bold text-white mb-4">Inventory</h3>
	<div class="grid grid-cols-2 gap-2">
		{#each Object.entries(PIECE_INVENTORY) as [rank, info]}
			{@const count = remainingCounts[rank] ?? 0}
			<button
				class="inventory-item flex items-center gap-2 p-2 rounded-xl border transition-all"
				class:active={selectedPieceRank === rank}
				class:out-of-stock={count === 0}
				onclick={() => onSelectRank(rank)}
				draggable={count > 0}
				ondragstart={(e) => onDragStart(e, rank)}
			>
				<div
					class="w-8 h-8 flex items-center justify-center text-xl bg-white/5 rounded-lg overflow-hidden p-1 shrink-0"
				>
					{#if ownerId === 2}
						{#if info.icon_red && (info.icon_red.includes('/') || info.icon_red.includes('.'))}
							<img src={info.icon_red} alt={info.name} class="w-full h-full object-contain" />
						{:else}
							{info.rank}
						{/if}
					{:else if info.icon_blue && (info.icon_blue.includes('/') || info.icon_blue.includes('.'))}
						<img src={info.icon_blue} alt={info.name} class="w-full h-full object-contain" />
					{:else}
						{info.rank}
					{/if}
				</div>
				<div class="flex-1 text-left min-w-0">
					<div class="text-[10px] font-bold text-white/40 uppercase leading-none mb-1 truncate">
						{info.name}
					</div>
					<div class="text-sm font-bold text-white flex justify-between">
						<span>{info.rank}</span>
						<span class={count > 0 ? 'text-brand-accent' : 'text-white/20'}>
							x{count}
						</span>
					</div>
				</div>
			</button>
		{/each}
	</div>

	{#if !isComplete}
		<p class="mt-4 text-[10px] text-white/30 italic">
			* You must place all 40 pieces to save the setup.
		</p>
	{/if}
</Card>

<style>
	.inventory-item {
		background: rgba(255, 255, 255, 0.03);
		border-color: rgba(255, 255, 255, 0.05);
	}

	.inventory-item:hover:not(.out-of-stock) {
		background: rgba(255, 255, 255, 0.07);
		border-color: rgba(255, 255, 255, 0.1);
	}

	.inventory-item.active {
		background: oklch(0.7 0.2 150 / 0.1);
		border-color: oklch(0.7 0.2 150);
	}

	.inventory-item.out-of-stock {
		opacity: 0.5;
		cursor: default;
	}
</style>
