<script lang="ts">
	import { BookOpen, ExternalLink } from '@lucide/svelte';
	import Card from '#lib/components/ui/Card.svelte';

	let {
		items,
		activeSlug,
		class: className,
		onSelect
	}: {
		items: { title: string; slug: string }[];
		activeSlug: string;
		class?: string;
		onSelect: (slug: string) => void;
	} = $props();
</script>

<Card padding="md" class="{className} space-y-4">
	<div class="flex items-center gap-2 pb-3 border-b border-white/10 text-brand-primary">
		<BookOpen class="w-4 h-4" />
		<span class="text-xs font-bold uppercase tracking-wider text-white">Documentation</span>
	</div>

	<nav class="space-y-1">
		{#each items as item (item.slug)}
			<button
				type="button"
				onclick={() => onSelect(item.slug)}
				class="w-full cursor-pointer text-left px-3 py-2 rounded-lg text-xs font-medium transition-colors flex items-center justify-between {activeSlug ===
				item.slug
					? 'bg-brand-primary/20 text-brand-primary border border-brand-primary/30'
					: 'text-white/70 hover:text-white hover:bg-white/5'}"
			>
				<span>{item.title}</span>
			</button>
		{/each}
	</nav>

	<div class="pt-3 border-t border-white/10">
		<a
			href={`https://github.com/${__REPO_NAME__}/wiki`}
			target="_blank"
			class="w-full cursor-pointer text-left px-3 py-2 rounded-lg text-xs font-medium transition-colors flex items-center justify-between text-white/70 hover:text-white hover:bg-white/5"
		>
			<span>View on GitHub Wiki</span>
			<ExternalLink class="w-3.5 h-3.5" />
		</a>
	</div>
</Card>
