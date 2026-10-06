<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Title from '$lib/components/Title.svelte';
	import { Save } from '@lucide/svelte';
	import { useGameSession } from '../_state/context';

	const session = useGameSession();
</script>

<div class="grid grid-cols-[1fr_auto_1fr] items-center mb-6">
	<div class="flex justify-start">
		{#if session.store.isGameOver}
			<Button
				onclick={() => {
					window.location.href = '/';
				}}
				variant="secondary"
			>
				Return to menu
			</Button>
		{:else}
			<Button
				variant="ghost"
				onclick={() => {
					if (confirm('Are you sure you want to quit?')) {
						session.abandonAndQuit();
					}
				}}
			>
				Quit Game
			</Button>
		{/if}
	</div>

	<Title />

	<div class="flex justify-end">
		<Button
			variant="outline"
			size="sm"
			onclick={() => session.saveGame()}
			disabled={!session.connected || !session.store.isGameOver}
			disabledMessage="Game must be finished to save a replay"
		>
			<Save class="mr-1.5 size-4" />
			&nbsp;Save Replay
		</Button>
	</div>
</div>
