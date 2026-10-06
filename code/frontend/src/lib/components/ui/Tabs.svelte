<script module lang="ts">
	export interface TabItem<K = string> {
		id: K;
		label: string;
		count?: number;
	}
</script>

<script lang="ts" generics="T extends string = string">
	import type { Snippet } from 'svelte';

	interface Props {
		tabs: TabItem<T>[];
		active?: T;
		onchange?: (id: T) => void;
		children?: Snippet<[T]>;
		class?: string;
	}

	let {
		tabs,
		active = $bindable(tabs[0]?.id),
		onchange,
		children,
		class: className = ''
	}: Props = $props();

	function selectTab(id: T) {
		active = id;
		onchange?.(id);
	}
</script>

<div class="space-y-6 {className}">
	<div
		class="flex items-center gap-2 p-1.5 rounded-xl bg-surface-elevated/30 border border-white/10 w-fit backdrop-blur-md"
	>
		{#each tabs as tab (tab.id)}
			<button
				type="button"
				class="px-4 py-2 rounded-lg text-xs font-semibold uppercase tracking-wider transition-all duration-200 cursor-pointer {active ===
				tab.id
					? 'bg-brand-primary text-black shadow-sm font-bold'
					: 'text-white/60 hover:text-white hover:bg-white/5'}"
				onclick={() => selectTab(tab.id)}
			>
				{tab.label}
				{#if tab.count !== undefined}
					<span
						class="ml-1 px-1.5 py-0.5 rounded-full text-[10px] {active === tab.id
							? 'bg-black/20 text-black'
							: 'bg-white/10 text-white/70'}"
					>
						{tab.count}
					</span>
				{/if}
			</button>
		{/each}
	</div>

	{#if children}
		{@render children(active)}
	{/if}
</div>
