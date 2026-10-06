<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/state';
	import { gamemodes } from '$lib/data/gamemodes.data';
	import Board from '$lib/components/game/Board.svelte';
	import Loading from '$lib/components/ui/Loading.svelte';
	import { GameSessionController } from './_state/game-session.svelte';
	import { setGameSession } from './_state/context';
	import ConnectionOverlay from './_components/ConnectionOverlay.svelte';
	import CombatAnimation from './_components/CombatAnimation.svelte';
	import SetupBanner from './_components/SetupBanner.svelte';
	import GameInfo from './_components/GameInfo.svelte';
	import GameTopBar from './_components/GameTopBar.svelte';
	import SetupInstructions from './_components/right-bar/SetupInstructions.svelte';
	import RightBar from './_components/right-bar/RightBar.svelte';

	const session = setGameSession(new GameSessionController());

	onMount(() => {
		session.init(page.params.id || '', new URLSearchParams(window.location.search));
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
{:else if session.store.gameState?.headless && !session.store.isGameOver}
	<Loading
		title="Game in progress"
		description="2 AI's are having the battle of their lives"
		subtitle="AI is thinking"
	/>
{:else}
	<GameTopBar />

	<div class="grid grid-cols-[280px_1fr_280px] gap-6 items-start">
		<div>
			<GameInfo />
		</div>

		<div class="flex justify-center">
			<Board
				boardState={session.store.boardState}
				selectedPosition={session.isSetupPhase
					? session.setupSwapPos1
					: session.store.selectedPosition}
				onCellClick={(x, y) => session.handleCellClick(x, y)}
				onCellDragStart={(e, x, y) => session.handleCellDragStart(e, x, y)}
				onCellDrop={(e, x, y) => session.handleCellDrop(e, x, y)}
				isInteractive={!session.store.isReplaying && (session.isHumanTurn || session.isSetupPhase)}
				viewerId={session.viewerId}
				validMoves={session.validMoves}
				disabledRows={session.disabledRows}
				visualDisabledRows={session.visualDisabledRows}
				highlightedRows={session.highlightedRows}
				highlightColor={session.highlightColor}
				scale={1.3}
				lastMove={session.store.lastMove}
				isSetupPhase={session.isSetupPhase}
				currentPlayerId={session.store.gameState?.currentPlayerId}
				isGameOver={session.store.isGameOver}
				winnerId={session.store.gameState?.winnerId}
			/>
		</div>

		<div>
			{#if !session.isSetupPhase}
				<RightBar />
			{:else}
				<SetupInstructions />
			{/if}
		</div>
	</div>
{/if}

{#if session.store.combatAnimation}
	<CombatAnimation
		attacker={session.store.combatAnimation.attacker}
		defender={session.store.combatAnimation.defender}
		attackerWon={session.store.combatAnimation.attackerWon}
		defenderWon={session.store.combatAnimation.defenderWon}
		onComplete={() => session.handleAnimationComplete()}
	/>
{/if}

{#if session.isSetupPhase && (session.store.gameMode.mode === gamemodes.human_vs_ai.mode || session.store.gameMode.mode === gamemodes.ai_vs_ai.mode)}
	<SetupBanner />
{/if}
