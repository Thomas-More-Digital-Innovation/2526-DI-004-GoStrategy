<script lang="ts">
	import { ArrowLeft, ArrowRight } from '@lucide/svelte';
	import Card from '#lib/components/ui/Card.svelte';

	let {
		html,
		activeSlug,
		items,
		onSelect
	}: {
		html: string;
		activeSlug: string;
		items: { title: string; slug: string }[];
		onSelect: (slug: string) => void;
	} = $props();

	const currentIndex = $derived(items.findIndex((i) => i.slug === activeSlug));
	const prevItem = $derived(currentIndex > 0 ? items[currentIndex - 1] : null);
	const nextItem = $derived(
		currentIndex >= 0 && currentIndex < items.length - 1 ? items[currentIndex + 1] : null
	);
</script>

<Card padding="lg" class="space-y-8 min-h-125 flex flex-col justify-between">
	<div class="space-y-6">
		<div class="flex items-center justify-between pb-4 border-b border-white/10 gap-4">
			<a
				href={`https://github.com/${__REPO_NAME__}/blob/main/docs/${activeSlug}.md`}
				target="_blank"
				rel="noopener noreferrer"
				class="font-mono hover:underline text-xs text-white/40">docs/{activeSlug}.md</a
			>
		</div>

		<article class="docs-prose max-w-none text-white/80 leading-relaxed text-sm">
			<!-- eslint-disable-next-line svelte/no-at-html-tags -->
			{@html html}
		</article>
	</div>

	<footer
		class="pt-6 border-t border-white/10 flex items-center justify-between gap-4 text-xs font-medium"
	>
		{#if prevItem}
			<button
				type="button"
				onclick={() => onSelect(prevItem.slug)}
				class="inline-flex items-center gap-2 text-white/60 hover:text-brand-primary transition-colors cursor-pointer"
			>
				<ArrowLeft class="w-4 h-4" />
				<span>Previous: {prevItem.title}</span>
			</button>
		{:else}
			<div></div>
		{/if}

		{#if nextItem}
			<button
				type="button"
				onclick={() => onSelect(nextItem.slug)}
				class="inline-flex items-center gap-2 text-white/60 hover:text-brand-primary transition-colors cursor-pointer"
			>
				<span>Next: {nextItem.title}</span>
				<ArrowRight class="w-4 h-4" />
			</button>
		{/if}
	</footer>
</Card>
