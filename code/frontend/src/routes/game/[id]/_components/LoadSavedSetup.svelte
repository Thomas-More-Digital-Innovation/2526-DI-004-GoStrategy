<script lang="ts">
    import type { BoardSetup } from "$lib/types/board-setup";
    import { gameStore } from "$lib/state/game.svelte";
    import BoardSetupCard from "$lib/components/setup/BoardSetupCard.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";

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
        onSelectSetup,
    }: Props = $props();

    const playerName = $derived(
        selectedPlayer === 0
            ? gameStore.gameState?.player1Username || "Player 1"
            : gameStore.gameState?.player2Username || "Player 2",
    );
</script>

<Modal
    bind:isOpen={showSelector}
    title={`Formation for ${playerName}`}
    description="Select a saved battle formation to deploy"
    onClose={() => (showSelector = false)}
    maxWidth="5xl"
    class={ownerId === 1 ? "border-brand-primary/40" : "border-brand-secondary/40"}
>
    <div class="flex justify-center flex-wrap gap-6 py-2">
        {#each savedSetups as setup}
            <BoardSetupCard
                {setup}
                {ownerId}
                isInteractive={true}
                onclick={() => onSelectSetup(setup.setup_data)}
            />
        {/each}
    </div>
</Modal>
