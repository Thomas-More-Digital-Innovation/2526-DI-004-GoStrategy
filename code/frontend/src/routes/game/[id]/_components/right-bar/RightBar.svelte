<script lang="ts">
    import Card from "$lib/components/ui/Card.svelte";
    import GameControls from "./GameControls.svelte";
    import GameHistory from "./GameHistory.svelte";

    interface Props {
        currentMoveIndex: number;
        totalMoves: number;
        isReplaying: boolean;
        onPrevious: () => void;
        onNext: () => void;
        onGoToMove: (index: number) => void;
        onExitReplay: () => void;
        onTogglePause: () => void;
        onSetSpeed: (speedMs: number) => void;
        onStep: () => void;
    }

    let {
        currentMoveIndex,
        totalMoves,
        isReplaying,
        onPrevious,
        onNext,
        onGoToMove,
        onExitReplay,
        onTogglePause,
        onSetSpeed,
        onStep,
    }: Props = $props();
</script>

<Card class="space-y-3 max-h-125 flex flex-col">
    <GameControls
        {isReplaying}
        {onSetSpeed}
        {onTogglePause}
        {onStep}
        {onExitReplay}
    />
    <hr class="border-white/10" />
    <GameHistory
        {currentMoveIndex}
        {totalMoves}
        {isReplaying}
        {onPrevious}
        {onNext}
        {onGoToMove}
        {onExitReplay}
    />
</Card>
