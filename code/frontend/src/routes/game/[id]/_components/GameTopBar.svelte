<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";
    import Title from "$lib/components/Title.svelte";

    interface Props {
        isGameOver: boolean;
        connected: boolean;
        onAbandonAndQuit: () => void;
        onSaveGame: () => void;
    }

    let { isGameOver, connected, onAbandonAndQuit, onSaveGame }: Props = $props();
</script>

<div class="grid grid-cols-[1fr_auto_1fr] items-center mb-6">
    <div class="flex justify-start">
        {#if isGameOver}
            <Button
                onclick={() => {
                    window.location.href = "/";
                }}
                variant="secondary"
            >
                Return to menu
            </Button>
        {:else}
            <Button
                variant="ghost"
                onclick={() => {
                    if (confirm("Are you sure you want to quit?")) {
                        onAbandonAndQuit();
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
            onclick={onSaveGame}
            disabled={!connected || !isGameOver}
            disabledMessage="Game must be finished to save a replay"
        >
            💾 Save Replay
        </Button>
    </div>
</div>
