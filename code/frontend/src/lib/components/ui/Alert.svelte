<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'error' | 'warning' | 'info' | 'success';
		message?: string;
		children?: Snippet;
		action?: Snippet;
		class?: string;
	}

	let {
		variant = 'error',
		message = '',
		children,
		action,
		class: className = ''
	}: Props = $props();

	const styles = {
		error: 'bg-brand-secondary/15 border-brand-secondary/30 text-brand-secondary',
		warning: 'bg-amber-500/15 border-amber-500/30 text-amber-300',
		info: 'bg-brand-primary/15 border-brand-primary/30 text-brand-primary',
		success: 'bg-emerald-500/15 border-emerald-500/30 text-emerald-300'
	};
</script>

<div
	class="rounded-xl border px-4 py-3 text-sm font-medium flex flex-col sm:flex-row items-center justify-between gap-3 animate-fade-in {styles[
		variant
	]} {className}"
	role="alert"
>
	<div class="flex-1 text-center sm:text-left">
		{#if message}
			<span>{message}</span>
		{/if}
		{#if children}
			{@render children()}
		{/if}
	</div>
	{#if action}
		<div class="shrink-0">
			{@render action()}
		</div>
	{/if}
</div>
