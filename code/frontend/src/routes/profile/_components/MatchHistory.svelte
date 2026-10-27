<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { games } from '#lib/api/client.js';
	import { toastStore } from '#lib/state/toast.svelte.js';
	import Card from '#lib/components/ui/Card.svelte';
	import Button from '#lib/components/ui/Button.svelte';
	import { Play, Swords } from '@lucide/svelte';
	import type { GameRecord } from '#lib/types/game.js';
	import PaginationControlls from './PaginationControlls.svelte';

	interface Props {
		currentUserId: number;
	}

	let { currentUserId }: Props = $props();

	const PAGE_SIZE = 30;
	let matchHistory = $state<GameRecord[]>([]);
	let totalGames = $state(0);
	let currentPage = $state(0);
	let isLoading = $state(true);

	const totalPages = $derived(Math.max(1, Math.ceil(totalGames / PAGE_SIZE)));

	async function loadPage(page: number) {
		if (page < 0 || (totalGames > 0 && page >= totalPages)) return;
		isLoading = true;
		try {
			const res = await games.listMine(PAGE_SIZE, page * PAGE_SIZE);
			matchHistory = res.games || [];
			console.log(matchHistory);
			totalGames = res.total || 0;
			currentPage = page;
		} catch (e) {
			toastStore.handleApiMessage(e, 'Failed to load match history');
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		loadPage(0);
	});

	function getOutcome(game: GameRecord): { label: string; class: string } {
		if (game.winner_id === null || game.winner_id === undefined) {
			return { label: 'Draw', class: 'bg-white/10 text-white/70 border-white/20' };
		}
		const isP1 = game.player1_user_id === currentUserId;
		const won = (isP1 && game.winner_id === 0) || (!isP1 && game.winner_id === 1);
		if (won) {
			return {
				label: 'Victory',
				class: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30'
			};
		}
		return { label: 'Defeat', class: 'bg-rose-500/15 text-rose-400 border-rose-500/30' };
	}

	function formatGameType(type: string): string {
		switch (type) {
			case 'human_vs_ai':
				return 'Human vs AI';
			case 'ai_vs_ai':
				return 'AI vs AI';
			case 'human_vs_human':
				return 'Human vs Human';
			default:
				return type.replace(/_/g, ' ');
		}
	}

	function formatDate(dateStr: string): string {
		const d = new Date(dateStr);
		return d.toLocaleDateString(undefined, {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function handleWatchReplay(gameId: string) {
		goto(resolve(`game/${gameId}?replay=true`));
	}
</script>

<Card class="space-y-4">
	<div class="flex items-center justify-between border-b border-white/10 pb-4 mb-0">
		<h2 class="text-lg font-bold text-white tracking-wide">Game History</h2>
		<span class="text-xs text-white/40">
			{totalGames}
			{totalGames === 1 ? 'match' : 'matches'} recorded
		</span>
	</div>

	{#if isLoading}
		<div class="py-8 text-center text-sm text-white/40">Loading recent matches...</div>
	{:else if matchHistory.length === 0}
		<div class="py-8 text-center space-y-2">
			<Swords class="size-8 mx-auto text-white/20" />
			<p class="text-sm text-white/40">No recorded matches found.</p>
			<p class="text-xs text-white/20">Completed games will appear here for full replay review.</p>
		</div>
	{:else}
		<div class="divide-y divide-white/5">
			<PaginationControlls {totalPages} {currentPage} {isLoading} {loadPage} />
			{#each matchHistory as game (game.id)}
				{@const outcome = getOutcome(game)}
				<div
					class="py-3 flex items-center justify-between gap-4 hover:bg-white/2 px-2 rounded-lg transition-colors"
				>
					<div class="space-y-1">
						<div class="flex items-center gap-2">
							<span class="text-sm font-semibold text-white capitalize">
								{formatGameType(game.game_type)}
							</span>
							<span
								class="text-[10px] font-bold border px-1.5 py-0.5 rounded uppercase tracking-wider {outcome.class}"
							>
								{outcome.label}
							</span>
						</div>
						<div class="text-xs text-white/40">
							{formatDate(game.created_at)}
						</div>
					</div>

					{#if game.has_replay}
						<Button variant="outline" size="sm" onclick={() => handleWatchReplay(game.id)}>
							<Play class="mr-1.5 size-3.5" />
							Watch Replay
						</Button>
					{:else}
						<span class="text-xs text-white/30 italic px-2 py-1 select-none">
							No replay available
						</span>
					{/if}
				</div>
			{/each}

			<PaginationControlls {totalPages} {currentPage} {isLoading} {loadPage} />

			<p class="text-xs text-white/40 mt-4">
				Note: AI vs AI game replays are kept for approximately 24 hours due to storage constraints.
				Other game replays are kept for 30 days.
			</p>
		</div>
	{/if}
</Card>
