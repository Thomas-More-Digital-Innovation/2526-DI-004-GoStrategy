<script lang="ts">
	import { Cpu, ExternalLink } from '@lucide/svelte';

	const stackLayers = [
		{
			name: 'Game Engine & API',
			tech: 'Go 1.26+ / Gin / WebSockets',
			description:
				'Real-time match loops, turn validation, Gorilla WebSocket pump, and Prometheus metrics.',
			path: 'code/backend'
		},
		{
			name: 'Web Interface',
			tech: 'SvelteKit 5 / Tailwind CSS / Runes',
			description: 'Reactive board interaction, piece movement animations, and profile management.',
			path: 'code/frontend'
		},
		{
			name: 'AI Simulation Binary',
			tech: 'Go Headless CLI Runner',
			description:
				'Benchmarking and automated tournament engine supporting MCTS, Minimax, Heuristic, and FATO algorithms.',
			path: 'code/backend/cmd/simulation'
		},
		{
			name: 'Persistence & Schema',
			tech: 'PostgreSQL 15 / GORM',
			description:
				'User authentication, board setup templates, and match history storage with migration guards.',
			path: 'code/backend/internal/db'
		}
	];
</script>

<section class="space-y-4">
	<div class="flex items-center gap-2 text-white/90">
		<Cpu class="w-5 h-5 text-brand-primary" />
		<h2 class="text-lg font-semibold tracking-wide">System Architecture</h2>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
		{#each stackLayers as layer (layer.name)}
			<div class="border border-white/10 bg-black/30 rounded-xl p-4 flex flex-col justify-between">
				<div>
					<div class="flex items-center justify-between">
						<span class="text-sm font-semibold text-white/90">{layer.name}</span>
						<a
							href={`https://github.com/${__REPO_NAME__}/tree/main/${layer.path}`}
							target="_blank"
							rel="noopener noreferrer"
							class="text-xs text-white/40 hover:text-white flex items-center gap-1 font-mono transition-colors"
							title="View source on GitHub"
						>
							<span>{layer.path}</span>
							<ExternalLink class="w-3 h-3" />
						</a>
					</div>
					<div class="text-xs font-mono text-brand-primary mt-1">{layer.tech}</div>
					<p class="text-xs text-white/60 mt-2 leading-relaxed">
						{layer.description}
					</p>
				</div>
			</div>
		{/each}
	</div>
</section>
