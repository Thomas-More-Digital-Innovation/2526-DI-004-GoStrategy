<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { FriendEntry } from '#lib/types/generated/models.js';

	interface Props {
		entry: FriendEntry;
		actions: Snippet;
	}

	let { entry, actions }: Props = $props();
</script>

<li class="flex items-center gap-4 rounded-xl bg-white/5 px-4 py-3">
	<div
		class="w-10 h-10 rounded-full bg-brand-secondary/30 border border-brand-accent/50 flex items-center justify-center text-brand-accent font-bold shrink-0"
	>
		{entry.user.username[0]?.toUpperCase() || '?'}
	</div>
	<span class="flex-1 font-bold text-white">{entry.user.username}</span>
	<div class="flex gap-2">
		{@render actions()}
	</div>
</li>
