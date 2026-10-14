<script lang="ts">
	import { Network, Bot, Database, Code } from '@lucide/svelte';
	import Card from '$lib/components/ui/Card.svelte';

	const pillars = [
		{
			title: 'Private & LAN Matchmaking',
			icon: Network,
			description:
				'Run dedicated game servers across private local networks or VPNs. Zero external cloud dependencies, zero latency jitter, and complete network isolation.'
		},
		{
			title: 'Headless AI Lab & Benchmarking',
			icon: Bot,
			description:
				'Simulate thousands of bot-vs-bot games at maximum CPU speed using the CLI engine. Benchmark MCTS, Minimax, Heuristics, or plug in your own decision models.'
		},
		{
			title: 'Full Data Sovereignty',
			icon: Database,
			description:
				'All user accounts, match histories, turn logs, and tactical board setups remain in your own PostgreSQL instance with full direct SQL access.'
		},
		{
			title: 'Extensible & Transparent Engine',
			icon: Code,
			description:
				'Inspect every line of game logic—turn lifecycles, piece movements, combat matrices, and WebSocket framing—under the permissive MIT license.'
		}
	];
</script>

<section class="space-y-4">
	<div class="flex items-center gap-2 text-white/90">
		<h2 class="text-lg font-semibold tracking-wide">Why Self-Host GoStrategy?</h2>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
		{#each pillars as pillar (pillar.title)}
			{@const Icon = pillar.icon}
			<Card padding="md" class="flex flex-col justify-between">
				<div>
					<div class="flex items-center gap-3">
						<div
							class="p-2 rounded-lg bg-white/5 border border-white/10 text-brand-primary shrink-0"
						>
							<Icon class="w-5 h-5" />
						</div>
						<h3 class="text-sm font-semibold text-white/90">{pillar.title}</h3>
					</div>
					<p class="text-xs text-white/60 mt-3 leading-relaxed">
						{pillar.description}
					</p>
				</div>
			</Card>
		{/each}
	</div>
</section>
