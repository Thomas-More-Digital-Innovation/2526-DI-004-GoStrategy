<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { page } from "$app/stores";
    import { gameStore } from "$lib/state/game.svelte";
    import { gamemodes } from "$lib/data/gamemodes.data";
    import Board from "$lib/components/game/Board.svelte";
    import Loading from "$lib/components/ui/Loading.svelte";
    import { GameSessionController } from "./_state/game-session.svelte";
    import ConnectionOverlay from "./_components/ConnectionOverlay.svelte";
    import CombatAnimation from "./_components/CombatAnimation.svelte";
    import SetupBanner from "./_components/SetupBanner.svelte";
    import GameInfo from "./_components/GameInfo.svelte";
    import GameTopBar from "./_components/GameTopBar.svelte";
    import SetupInstructions from "./_components/right-bar/SetupInstructions.svelte";
    import RightBar from "./_components/right-bar/RightBar.svelte";

    const session = new GameSessionController();

    onMount(() => {
        session.init(
            $page.params.id || "",
            new URLSearchParams(window.location.search),
        );
    });

    onDestroy(() => {
        session.destroy();
    });
</script>

<svelte:head>
    <title>GoStrategy — Game {session.gameId}</title>
</svelte:head>

{#if !session.connected}
    <ConnectionOverlay
        isReconnecting={session.isReconnecting}
        gameId={session.gameId}
        reconnectAttempts={session.reconnectAttempts}
        maxReconnectAttempts={session.maxReconnectAttempts}
        onRetry={() => session.attemptReconnect()}
        onReturnToMenu={() => session.abandonAndQuit()}
    />
{:else if gameStore.gameState?.headless && !gameStore.gameState?.isGameOver}
    <Loading
        title="Game in progress"
        description="2 AI's are having the battle of their lives"
        subtitle="AI is thinking"
    />
{:else}
    <GameTopBar
        isGameOver={gameStore.gameState?.isGameOver ?? false}
        connected={session.connected}
        onAbandonAndQuit={() => session.abandonAndQuit()}
        onSaveGame={() => session.saveGame()}
    />

    <div class="grid grid-cols-[280px_1fr_280px] gap-6 items-start">
        <div>
            <GameInfo
                gameState={gameStore.gameState}
                gameMode={gameStore.gameMode}
            />
        </div>

        <div class="flex justify-center">
            <Board
                boardState={gameStore.boardState}
                selectedPosition={session.isSetupPhase
                    ? session.setupSwapPos1
                    : gameStore.selectedPosition}
                onCellClick={(x, y) => session.handleCellClick(x, y)}
                onCellDragStart={(e, x, y) =>
                    session.handleCellDragStart(e, x, y)}
                onCellDrop={(e, x, y) => session.handleCellDrop(e, x, y)}
                isInteractive={!gameStore.isReplaying &&
                    (session.isHumanTurn || session.isSetupPhase)}
                viewerId={session.viewerId}
                validMoves={session.validMoves}
                disabledRows={session.disabledRows}
                visualDisabledRows={session.visualDisabledRows}
                highlightedRows={session.highlightedRows}
                highlightColor={session.highlightColor}
                scale={1.3}
                lastMove={gameStore.lastMove}
            />
        </div>

        <div>
            {#if !session.isSetupPhase}
                <RightBar
                    currentMoveIndex={gameStore.currentHistoryIndex}
                    totalMoves={gameStore.history.length}
                    isReplaying={gameStore.isReplaying}
                    onPrevious={() => session.handlePreviousMove()}
                    onNext={() => session.handleNextMove()}
                    onGoToMove={(index) => session.handleGoToMove(index)}
                    onExitReplay={() => session.handleExitReplay()}
                    onTogglePause={() => session.handleTogglePause()}
                    onSetSpeed={(speed) => session.handleSetSpeed(speed)}
                    onStep={() => session.handleStep()}
                />
            {:else}
                <SetupInstructions />
            {/if}
        </div>
    </div>
{/if}

{#if gameStore.combatAnimation}
    <CombatAnimation
        attacker={gameStore.combatAnimation.attacker}
        defender={gameStore.combatAnimation.defender}
        attackerWon={gameStore.combatAnimation.attackerWon}
        defenderWon={gameStore.combatAnimation.defenderWon}
        onComplete={() => session.handleAnimationComplete()}
    />
{/if}

{#if session.isSetupPhase && (gameStore.gameMode.mode === gamemodes.human_vs_ai.mode || gameStore.gameMode.mode === gamemodes.ai_vs_ai.mode)}
    <SetupBanner
        onRandomize={(p) => session.handleRandomize(p)}
        onStart={(h) => session.handleStartGame(h)}
        onLoadSetup={(setup, p) => session.handleLoadSetup(setup, p)}
        onBackToMenu={() => session.abandonAndQuit()}
        viewerId={session.viewerId}
        gameMode={gameStore.gameMode}
        selectedPlayer={session.setupSelectedPlayer}
        onSelectPlayer={(p: number) => {
            session.setupSelectedPlayer = p;
            session.setupSwapPos1 = null;
        }}
    />
{/if}
