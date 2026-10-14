<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import { Copy, Check, ExternalLink, Terminal } from '@lucide/svelte';

	let copied = $state(false);
	const composeCmd = 'docker compose up -d';

	async function copyCommand() {
		try {
			await navigator.clipboard.writeText(composeCmd);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		} catch {
			// ponytail: fallback if clipboard API unavailable
		}
	}
</script>

<Card padding="md" class="space-y-4">
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div class="space-y-1">
			<h2 class="text-base font-semibold text-white tracking-wide">
				Ready to run your own instance?
			</h2>
			<p class="text-xs text-white/60 max-w-xl">
				The repository includes production and development Docker Compose configurations for the Go
				backend, SvelteKit web client, and PostgreSQL database.
			</p>
		</div>

		<div class="flex items-center gap-3 shrink-0">
			<Button
				href={`https://github.com/${__REPO_NAME__}#running-with-docker-recommended`}
				target="_blank"
				variant="primary"
				size="md"
			>
				<span>Docker Guide in README</span>
				<ExternalLink class="ml-2 w-3.5 h-3.5" />
			</Button>
		</div>
	</div>

	<div
		class="flex items-center justify-between gap-3 bg-black/70 border border-white/10 rounded-xl px-4 py-3 font-mono text-xs text-brand-accent/90"
	>
		<div class="flex items-center gap-2.5">
			<Terminal class="w-4 h-4 text-white/40" />
			<span class="select-all">{composeCmd}</span>
		</div>
		<Button variant="ghost" size="md" onclick={copyCommand}>
			{#if copied}
				<Check class="w-4 h-4 text-emerald-400" />
			{:else}
				<Copy class="w-4 h-4" />
			{/if}
		</Button>
	</div>
</Card>
