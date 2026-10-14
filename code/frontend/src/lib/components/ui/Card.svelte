<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		class?: string;
		padding?: 'none' | 'sm' | 'md' | 'lg';
		onclick?: () => void;
		as?: 'div' | 'button';
		disabled?: boolean;
		type?: 'button' | 'submit' | 'reset';
	}

	let {
		children,
		class: className = '',
		padding = 'md',
		onclick,
		as = 'div',
		disabled = false,
		type = 'button'
	}: Props = $props();

	const paddings = {
		none: '',
		sm: 'p-3',
		md: 'p-6',
		lg: 'p-10'
	};
</script>

{#if as === 'button'}
	<button
		{type}
		{disabled}
		class="glass rounded-2xl overflow-hidden text-left transition-all {paddings[
			padding
		]} {className}"
		{onclick}
	>
		{@render children()}
	</button>
{:else}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="glass rounded-2xl overflow-hidden {paddings[padding]} {className}" {onclick}>
		{@render children()}
	</div>
{/if}
