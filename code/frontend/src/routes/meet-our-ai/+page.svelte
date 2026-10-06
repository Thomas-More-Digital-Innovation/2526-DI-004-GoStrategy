<script lang="ts">
	import { aiDossiers, tournamentBenchmarks } from '$lib/data/ai-intel.data';
	import AIDossierCard from './_components/AIDossierCard.svelte';
	import AIComparisonTable from './_components/AIComparisonTable.svelte';
	import AITournamentBenchmark from './_components/AITournamentBenchmark.svelte';
	import Tabs from '$lib/components/ui/Tabs.svelte';

	type ViewMode = 'dossiers' | 'matrix' | 'benchmarks';
	let activeView = $state<ViewMode>('dossiers');

	const tabs: { id: ViewMode; label: string }[] = [
		{ id: 'dossiers', label: 'Engines' },
		{ id: 'matrix', label: 'Specifications' },
		{ id: 'benchmarks', label: 'Tournament Results' }
	];
</script>

<svelte:head>
	<title>GoStrategy — Meet our AI</title>
</svelte:head>

<div class="space-y-6 pb-12">
	<header>
		<h1 class="text-2xl font-bold text-white uppercase tracking-wider">Meet our AI</h1>
		<p class="text-white/50 text-xs mt-1">
			Technical specifications and tournament benchmarks for GoStrategy AI opponents.
		</p>
	</header>

	<Tabs {tabs} bind:active={activeView}>
		{#snippet children(current)}
			{#if current === 'dossiers'}
				<div class="space-y-4">
					{#each aiDossiers as dossier}
						<AIDossierCard {dossier} />
					{/each}
				</div>
			{:else if current === 'matrix'}
				<AIComparisonTable dossiers={aiDossiers} />
			{:else if current === 'benchmarks'}
				<AITournamentBenchmark benchmarks={tournamentBenchmarks} />
			{/if}
		{/snippet}
	</Tabs>
</div>
