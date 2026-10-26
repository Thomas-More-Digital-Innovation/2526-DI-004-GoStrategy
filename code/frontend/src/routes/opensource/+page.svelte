<script lang="ts">
	import { FileText, Bug } from '@lucide/svelte';
	import Button from '#lib/components/ui/Button.svelte';
	import PillarsGrid from './_components/PillarsGrid.svelte';
	import ArchitectureStack from './_components/ArchitectureStack.svelte';
	import DeploymentCTA from './_components/DeploymentCTA.svelte';
</script>

<svelte:head>
	<title>GoStrategy — Open Source & Self-Hosting</title>
	<meta
		name="description"
		content="GoStrategy is 100% open-source software under the MIT license. Learn why self-hosting gives you private matchmaking, headless AI benchmarking, and full data sovereignty."
	/>
</svelte:head>

<div class="space-y-12 pb-12">
	<header
		class="min-h-[60vh] flex flex-col justify-center border-b border-white/10 py-12 space-y-8"
	>
		<div class="space-y-4 max-w-4xl">
			<h1 class="text-4xl sm:text-5xl lg:text-6xl font-black text-white tracking-tight">
				GoStrategy is opensource!
			</h1>
			<p class="text-xl sm:text-2xl font-medium text-brand-primary">Self host it, or contribute</p>
		</div>

		<div class="flex flex-wrap items-center gap-4 pt-2">
			<Button
				href={`https://github.com/${__REPO_NAME__}`}
				target="_blank"
				rel="noopener noreferrer"
				variant="primary"
				size="md"
			>
				<span class="flex items-center gap-2">
					<svg class="w-4 h-4 fill-current" viewBox="0 0 24 24" aria-hidden="true">
						<path
							fill-rule="evenodd"
							clip-rule="evenodd"
							d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.53 1.032 1.53 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
						/>
					</svg>
					<span>View on GitHub</span>
				</span>
			</Button>

			<Button
				href={`https://github.com/${__REPO_NAME__}/issues`}
				target="_blank"
				rel="noopener noreferrer"
				variant="outline"
				size="md"
			>
				<span class="flex items-center gap-2">
					<Bug class="w-4 h-4" />
					<span>Issues</span>
				</span>
			</Button>

			<Button
				href={`https://github.com/${__REPO_NAME__}/blob/main/LICENSE`}
				target="_blank"
				rel="noopener noreferrer"
				variant="ghost"
				size="md"
			>
				<span class="flex items-center gap-2">
					<FileText class="w-4 h-4" />
					<span>MIT</span>
				</span>
			</Button>
		</div>
	</header>

	<PillarsGrid />

	<ArchitectureStack />

	<DeploymentCTA />

	<footer
		class="border-t border-white/10 pt-6 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-white/50"
	>
		<div>
			Looking to extend game mechanics or add new AI algorithms? Read the
			<a
				href={`https://github.com/${__REPO_NAME__}/blob/main/CONTRIBUTING.md`}
				target="_blank"
				rel="noopener noreferrer"
				class="text-brand-primary underline hover:text-brand-primary/80"
			>
				Contributing Guide
			</a>.
		</div>
		<div class="font-mono text-[11px] text-white/40">
			License: MIT &bull; Go 1.26+ &bull; SvelteKit 5
		</div>
	</footer>
</div>
