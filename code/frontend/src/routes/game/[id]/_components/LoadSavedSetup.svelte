<script lang="ts">
	import type { BoardSetup } from '$lib/types/board-setup';
	import BoardSetupCard from '$lib/components/setup/BoardSetupCard.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import { useGameSession } from '../_state/context';

	interface Props {
		savedSetups: BoardSetup[];
		ownerId: number;
		selectedPlayer: number;
		onSelectSetup: (setupData: string) => void;
		showSelector: boolean;
	}

	let {
		savedSetups,
		ownerId,
		selectedPlayer,
		showSelector = $bindable(),
		onSelectSetup
	}: Props = $props();

	const session = useGameSession();

	const playerName = $derived(
		selectedPlayer === 0
			? session.store.gameState?.player1Username || 'Player 1'
			: session.store.gameState?.player2Username || 'Player 2'
	);
</script>

<Modal
	bind:isOpen={showSelector}
	title={`Formation for ${playerName}`}
	description="Select a saved battle formation to deploy"
	onClose={() => (showSelector = false)}
	maxWidth="5xl"
	class={ownerId === 1 ? 'border-brand-primary/40' : 'border-brand-secondary/40'}
>
	<div class="flex justify-center flex-wrap gap-6 py-2">
		{#each savedSetups as setup (setup.id)}
			<BoardSetupCard
				{setup}
				{ownerId}
				isInteractive={true}
				onclick={() => onSelectSetup(setup.setup_data)}
			/>
		{/each}
	</div>
</Modal>
