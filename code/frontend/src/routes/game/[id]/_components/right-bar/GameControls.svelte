<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import { gamemodes } from '$lib/data/gamemodes.data';
	import { Play, Pause, StepForward } from '@lucide/svelte';
	import { useGameSession } from '../../_state/context';

	const session = useGameSession();

	let speedMs = $state(1000);
</script>

<h3 class="text-sm font-bold text-brand-accent uppercase tracking-wider">Controls</h3>
<div class="flex flex-col gap-2">
	<div class="flex gap-2">
		{#if session.store.isReplaying && !session.store.isGameOver}
			<Button
				variant="secondary"
				size="sm"
				class="flex-1"
				onclick={() => session.handleExitReplay()}
			>
				Exit Replay
			</Button>
		{:else if !session.store.isGameOver}
			<Button
				variant="outline"
				size="sm"
				onclick={() => session.handleTogglePause()}
				class="flex-1"
			>
				{#if session.store.isPaused}
					<Play class="mr-1.5 size-3.5" />
					Resume
				{:else}
					<Pause class="mr-1.5 size-3.5" />
					Pause
				{/if}
			</Button>
			{#if session.store.isPaused && session.store.gameMode.mode === gamemodes.ai_vs_ai.mode}
				<Button
					variant="ghost"
					size="sm"
					onclick={() => session.handleStep()}
					loading={session.store.isStepping}
					disabled={session.store.isStepping}
					disabledMessage="Processing move..."
				>
					<StepForward class="mr-1.5 size-3.5" />
					Step
				</Button>
			{/if}
		{/if}
	</div>

	<div class="space-y-1.5 px-1">
		<div class="flex justify-between text-[10px] text-white/40">
			<span>Speed</span>
			<span>{speedMs}ms</span>
		</div>
		<input
			type="range"
			min="500"
			max="5000"
			step="100"
			bind:value={speedMs}
			onchange={() => session.handleSetSpeed(speedMs)}
			class="w-full accent-brand-primary h-1 bg-white/10 rounded-lg appearance-none cursor-pointer"
		/>
	</div>
</div>
