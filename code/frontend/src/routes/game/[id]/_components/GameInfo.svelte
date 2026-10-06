<script lang="ts">
    import Card from "$lib/components/ui/Card.svelte";
    import type { GameMode, GameState } from "$lib/types/game";
    import { gamemodes } from "$lib/data/gamemodes.data";
    import { AIs } from "$lib/data/AI.data";
    import { useGameSession } from "../_state/context";

    interface Props {
        gameState?: GameState | null;
        gameMode?: GameMode;
    }

    let { gameState: propGameState, gameMode: propGameMode }: Props = $props();

    const session = useGameSession();
    const gameState = $derived(propGameState ?? session.store.gameState);
    const gameMode = $derived(propGameMode ?? session.store.gameMode);

    function formatPossessiveUsername(username: string) {
        return username.endsWith("s") ? `${username}'` : `${username}'s`;
    }

    function formatAiName(raw?: string): string {
        if (!raw) return "AI";
        const numMatch = raw.match(/\s+\d+$/);
        const numSuffix = numMatch ? numMatch[0] : "";
        const base = raw.replace(/\s+\d+$/, "").trim();
        const match = AIs.find(
            (a) =>
                a.id.toLowerCase() === base.toLowerCase() ||
                a.name.toLowerCase() === base.toLowerCase(),
        );
        const name = match ? match.name : base;
        return `${name}${numSuffix}`;
    }

    const title = $derived.by(() => {
        if (!gameState) return gameMode.title;

        if (
            gameMode.mode === gamemodes.human_vs_ai.mode ||
            gameMode.mode === "human_vs_ai"
        ) {
            if (!gameState.player2Username) return gameMode.title;
            const ai = formatAiName(gameState.player2Username);
            return `you vs ${ai} ai`;
        }

        if (
            gameMode.mode === gamemodes.ai_vs_ai.mode ||
            gameMode.mode === "ai_vs_ai"
        ) {
            if (!gameState.player1Username || !gameState.player2Username) {
                return gameMode.title;
            }
            const ai1 = formatAiName(gameState.player1Username);
            const ai2 = formatAiName(gameState.player2Username);
            return `${ai1} ai vs ${ai2} ai`;
        }
        // TODO: add human vs human title once implemented
        return gameMode.title;
    });
</script>

<Card class="space-y-4">
    <div class="border-b border-white/10 pb-3">
        <h2
            class="text-md font-bold text-brand-accent uppercase tracking-wider"
        >
            {title}
        </h2>
    </div>

    {#if gameState}
        <div class="space-y-2 text-sm">
            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50">Round</span>
                <span class="font-semibold text-white">{gameState.round}</span>
            </div>

            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50">Current Turn</span>
                <span
                    class="font-semibold"
                    class:text-brand-secondary={gameState.currentPlayerId === 0}
                    class:text-brand-primary={gameState.currentPlayerId === 1}
                >
                    {gameState.currentPlayerName}
                </span>
            </div>

            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50"
                    >{formatPossessiveUsername(gameState.player1Username)}
                    Pieces</span
                >
                <span class="font-semibold text-brand-secondary"
                    >{gameState.player1AlivePieces}</span
                >
            </div>

            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50"
                    >{formatPossessiveUsername(gameState.player2Username)} Pieces</span
                >
                <span class="font-semibold text-brand-primary"
                    >{gameState.player2AlivePieces}</span
                >
            </div>

            {#if !gameState.isSetupPhase}
                <div class="flex justify-between items-center py-1">
                    <span class="text-white/50">Moves</span>
                    <span class="font-semibold text-white"
                        >{gameState.moveCount}</span
                    >
                </div>
            {/if}
        </div>

        {#if gameState.isGameOver}
            <div
                class="rounded-xl bg-brand-accent/20 border border-brand-accent/30 p-4 text-center space-y-1"
            >
                <h3
                    class="text-brand-accent font-bold uppercase tracking-wider"
                >
                    Game Over
                </h3>
                <p class="text-white text-sm">
                    Winner: <strong>{gameState.winnerName || "Draw"}</strong>
                </p>
                <p class="text-white/60 text-xs">{gameState.winCause}</p>
            </div>
        {/if}
    {:else}
        <p class="text-white/30 text-center py-4">No game data</p>
    {/if}
</Card>
