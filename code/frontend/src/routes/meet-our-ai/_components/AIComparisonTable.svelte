<script lang="ts">
	import type { AIDossier } from '$lib/types/ai';

	let { dossiers }: { dossiers: AIDossier[] } = $props();
</script>

<div class="border border-white/10 bg-surface-elevated/40 rounded-xl overflow-hidden">
	<div class="p-5 border-b border-white/10">
		<h2 class="text-base font-bold text-white uppercase tracking-wider">Engine Specifications</h2>
	</div>

	<div class="overflow-x-auto">
		<table class="w-full text-left text-xs">
			<thead>
				<tr
					class="border-b border-white/10 bg-black/20 text-[11px] uppercase tracking-wider text-white/50"
				>
					<th class="py-3 px-4 font-semibold">Engine</th>
					<th class="py-3 px-4 font-semibold">Category</th>
					<th class="py-3 px-4 font-semibold">Complexity</th>
					<th class="py-3 px-4 font-semibold">State Model</th>
					<th class="py-3 px-4 font-semibold">Training Support</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-white/5">
				{#each dossiers as ai (ai.id)}
					<tr class="hover:bg-white/5 transition-colors">
						<td class="py-3 px-4 font-bold text-white flex items-center gap-3">
							<img
								src={ai.image}
								alt={ai.name}
								class="w-7 h-7 rounded object-cover bg-black border border-white/10 shrink-0"
							/>
							<span class="text-sm">{ai.name}</span>
						</td>
						<td class="py-3 px-4 text-white/80">
							{ai.category}
						</td>
						<td class="py-3 px-4 font-mono text-white/70">
							{ai.complexity}
						</td>
						<td class="py-3 px-4 text-white/70">
							{ai.stateModel}
						</td>
						<td class="py-3 px-4 text-white/80">
							{ai.trainingSupport}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
