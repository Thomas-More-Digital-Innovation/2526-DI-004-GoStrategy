<script lang="ts">
    import { tick } from "svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { ChevronLeft, ChevronRight } from "@lucide/svelte";
    import { useGameSession } from "../../_state/context";

    const session = useGameSession();

    const currentMoveIndex = $derived(session.store.currentHistoryIndex);
    const totalMoves = $derived(session.store.history.length);
    const isReplaying = $derived(session.store.isReplaying);

    const canGoPrevious = $derived(currentMoveIndex > 0);
    const canGoNext = $derived(currentMoveIndex < totalMoves - 1);

    let scrollContainer = $state<HTMLDivElement | null>(null);
    let isInitialMount = true;

    $effect(() => {
        if (totalMoves > 0 && !isReplaying) {
            const smooth = !isInitialMount;
            isInitialMount = false;
            tick().then(() => {
                if (scrollContainer) {
                    scrollContainer.scrollTo({
                        top: scrollContainer.scrollHeight,
                        behavior: smooth ? "smooth" : "auto",
                    });
                }
            });
        }
    });
</script>

<div class="flex items-center justify-between">
    <h3 class="text-sm font-bold text-brand-accent uppercase tracking-wider">
        Move History
    </h3>
    {#if session.store.isGameOver}
        <span
            class="text-[10px] font-bold border border-green-500/30 bg-green-500/10 text-green-400 px-2 py-0.5 rounded-md uppercase tracking-wider"
        >
            Finished
        </span>
    {:else if isReplaying}
        <span
            class="text-[10px] font-bold border border-brand-secondary/40 bg-brand-secondary/15 text-brand-secondary px-2 py-0.5 rounded-md uppercase tracking-wider"
        >
            Replay
        </span>
    {/if}
    {#if session.store.isPaused && !session.store.isGameOver}
        <span
            class="text-[10px] font-bold border border-amber-500/40 bg-amber-500/15 text-amber-400 px-2 py-0.5 rounded-md uppercase tracking-wider"
        >
            Paused
        </span>
    {/if}
</div>

{#if totalMoves > 0}
    <p class="text-white/40 text-xs text-center">
        Move {currentMoveIndex + 1} of {totalMoves}
    </p>

    <div class="grid grid-cols-2 gap-2">
        <Button
            variant="outline"
            size="sm"
            onclick={() => session.handlePreviousMove()}
            disabled={!canGoPrevious}
        >
            <ChevronLeft class="mr-1 size-3.5" />
            Prev
        </Button>
        <Button
            variant="outline"
            size="sm"
            onclick={() => session.handleNextMove()}
            disabled={!canGoNext}
        >
            Next
            <ChevronRight class="ml-1 size-3.5" />
        </Button>
    </div>

    <div
        bind:this={scrollContainer}
        class="custom-scrollbar flex-1 overflow-y-auto space-y-1 min-h-0 pr-1"
    >
        {#each Array(totalMoves) as _, index}
            <button
                class="w-full text-left px-3 py-1.5 rounded-lg text-xs transition-all {index ===
                currentMoveIndex
                    ? 'bg-brand-primary/20 text-brand-primary font-semibold'
                    : 'text-white/40 hover:bg-white/5 hover:text-white/70'}"
                onclick={() => {
                    if (index === totalMoves - 1 && !session.store.isGameOver) {
                        session.handleExitReplay();
                    } else {
                        session.handleGoToMove(index);
                    }
                }}
            >
                Move {index + 1}
            </button>
        {/each}
    </div>
{:else}
    <p class="text-white/30 text-center py-4 text-sm">No moves yet</p>
{/if}

<style>
    .custom-scrollbar {
        scrollbar-width: thin;
        scrollbar-color: rgba(255, 255, 255, 0.18) transparent;
    }
    .custom-scrollbar::-webkit-scrollbar {
        width: 5px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.18);
        border-radius: 9999px;
    }
</style>
