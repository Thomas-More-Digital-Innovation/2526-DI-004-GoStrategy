<script lang="ts">
    import { onMount } from "svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { boardSetups } from "$lib/api/client";
    import { flipSetup } from "$lib/utils/board-binary";
    import type { BoardSetup } from "$lib/types/board-setup";
    import { gamemodes } from "$lib/data/gamemodes.data";
    import { Clock, FolderOpen, Dices } from "@lucide/svelte";
    import LoadSavedSetup from "./LoadSavedSetup.svelte";
    import AiPlayerSelector from "./AiPlayerSelector.svelte";
    import { useGameSession } from "../_state/context";

    const session = useGameSession();

    const ownerId = $derived(
        session.store.gameMode.mode === gamemodes.ai_vs_ai.mode
            ? session.setupSelectedPlayer === 0
                ? 2
                : 1
            : session.viewerId === -1
              ? 2
              : session.viewerId === 0
                ? 2
                : 1,
    );

    let savedSetups = $state<BoardSetup[]>([]);
    let loadingSetups = $state(true);
    let showSelector = $state(false);
    let headless = $state(false);

    onMount(async () => {
        try {
            const result = await boardSetups.list();
            savedSetups = result ?? [];
        } catch {
            // Silently fail
        } finally {
            loadingSetups = false;
        }
    });

    function selectSetup(setupData: string) {
        let finalSetup = setupData;
        if (
            session.store.gameMode.mode === gamemodes.ai_vs_ai.mode &&
            session.setupSelectedPlayer === 1
        ) {
            finalSetup = flipSetup(setupData);
        }
        session.handleLoadSetup(
            finalSetup,
            session.store.gameMode.mode === gamemodes.ai_vs_ai.mode
                ? session.setupSelectedPlayer
                : undefined,
        );
        showSelector = false;
    }

    function formatTime(seconds: number | null): string {
        if (seconds === null) return "00:00";
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${mins}:${secs.toString().padStart(2, "0")}`;
    }
</script>

{#snippet TimeLimit()}
    <div class="group cursor-help relative">
        <p
            class="text-white/40 text-[10px] font-bold uppercase tracking-wider flex items-center gap-1.5 w-fit"
        >
            <Clock class="size-3 text-white/50" />
            <span>Time Remaining:</span>
            <span class="text-white font-mono"
                >{formatTime(session.store.setupRemainingSecs)}</span
            >
        </p>
        <div
            class="hidden group-hover:block absolute top-full left-0 mt-2 bg-surface-elevated rounded-lg p-3 text-xs w-64 text-white shadow-xl border border-white/10 z-50 animate-in fade-in slide-in-from-top-1 duration-200"
        >
            A time limit is set for this game to prevent infinite setup times.
            The game will start automatically when the time runs out.
        </div>
    </div>
{/snippet}

<div class="fixed top-0 left-0 right-0 z-50 pointer-events-none">
    <div
        class="glass pointer-events-auto flex items-center justify-between gap-6 px-8 py-4 border-b border-white/10"
    >
        <div class="flex items-center gap-3">
            <Button variant="outline" onclick={() => session.abandonAndQuit()}>
                Back To Menu
            </Button>

            <div>
                <div class="flex gap-4 items-center">
                    <h2
                        class="text-lg font-black text-white uppercase tracking-tighter"
                    >
                        Setup Phase
                    </h2>
                    {@render TimeLimit()}
                </div>

                {#if session.store.gameMode.mode === gamemodes.ai_vs_ai.mode}
                    <AiPlayerSelector
                        selectedPlayer={session.setupSelectedPlayer}
                        onSelectPlayer={(p: number) => {
                            session.setupSelectedPlayer = p;
                            session.setupSwapPos1 = null;
                        }}
                    />
                {:else}
                    <div class="flex flex-col gap-0.5">
                        <p class="text-white/40 text-xs font-medium">
                            Arrange your pieces or load a formation
                        </p>
                    </div>
                {/if}
            </div>
        </div>

        <div class="flex gap-3 items-center">
            {#if session.store.gameMode.mode === gamemodes.ai_vs_ai.mode}
                <div
                    class="flex items-center gap-2 px-3 py-2 bg-white/5 rounded-xl border border-white/10"
                >
                    <input
                        type="checkbox"
                        id="headless-mode"
                        bind:checked={headless}
                        class="accent-brand-primary cursor-pointer"
                    />
                    <label
                        for="headless-mode"
                        class="text-[10px] font-bold text-white/60 uppercase tracking-wider cursor-pointer"
                    >
                        Headless Mode
                    </label>
                </div>
            {/if}
            <Button
                variant="ghost"
                onclick={() => (showSelector = true)}
                class="bg-white/5! hover:bg-white/10!"
                disabled={loadingSetups || savedSetups.length === 0}
            >
                <FolderOpen class="size-4" />
                &nbsp; Load Setup
            </Button>

            <Button
                variant="outline"
                onclick={() =>
                    session.handleRandomize(
                        session.store.gameMode.mode === gamemodes.ai_vs_ai.mode
                            ? session.setupSelectedPlayer
                            : undefined,
                    )}
            >
                <Dices class="size-4" />
                &nbsp; Randomize
            </Button>
            <Button
                variant="primary"
                onclick={() => session.handleStartGame(headless)}
            >
                Start Game
            </Button>
        </div>
    </div>
</div>

{#if showSelector}
    <LoadSavedSetup
        {savedSetups}
        {ownerId}
        selectedPlayer={session.setupSelectedPlayer}
        onSelectSetup={selectSetup}
        bind:showSelector
    />
{/if}
