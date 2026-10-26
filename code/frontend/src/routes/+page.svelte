<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { authStore } from '#lib/state/auth.svelte.js';
	import Button from '#lib/components/ui/Button.svelte';
	import Card from '#lib/components/ui/Card.svelte';
	import Alert from '#lib/components/ui/Alert.svelte';
	import type { GameMode } from '#lib/types/game.js';
	import { allGamemodes, gamemodes } from '#lib/data/gamemodes.data.js';
	import ChangelogButton from '#lib/components/changelog/ChangelogButton.svelte';
	import ReconnectBanner from '#lib/components/ReconnectBanner.svelte';

	let selectedMode = $state<GameMode>(gamemodes.human_vs_ai);
	let error = $state('');

	function startFlow() {
		error = '';
		if (selectedMode.mode === gamemodes.human_vs_human.mode) {
			error = 'Coming Soon';
			return;
		}
		goto(resolve(`select-ai?mode=${selectedMode.mode}`));
	}
</script>

<svelte:head>
	<title>GoStrategy — Command Center</title>
</svelte:head>

<div class="space-y-8">
	<header>
		<h1 class="text-3xl font-extrabold text-white uppercase tracking-widest">Command Center</h1>
		<p class="text-white/50 mt-1">
			{#if authStore.user}
				Welcome back, <span class="text-brand-accent font-semibold">{authStore.user.username}</span
				>.
			{:else}
				<a href={resolve('login')} class="text-brand-primary hover:underline">Sign in</a>

				to start playing.
			{/if}
		</p>
	</header>

	<ReconnectBanner />

	{#if error}
		<Alert variant="error" message={error} />
	{/if}

	{#if selectedMode}
		<!-- Game Mode Selection -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
			{#each allGamemodes as gamemode (gamemode.mode)}
				<Card
					as="button"
					padding="md"
					disabled={gamemode.disabled}
					onclick={() => (selectedMode = gamemode)}
					class="cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed transition-all duration-200 border-2 {selectedMode.mode ===
					gamemode.mode
						? 'border-brand-primary! bg-brand-primary/10! scale-[1.02]'
						: 'border-white/5 hover:border-white/20 hover:bg-white/5'}"
				>
					<div class="text-3xl mb-3">{gamemode.icon}</div>
					<h3 class="font-bold text-white text-lg tracking-wide uppercase">
						{gamemode.title}
					</h3>
					<p class="text-white/50 text-xs mt-1 leading-relaxed">
						{gamemode.desc}
					</p>
				</Card>
			{/each}
		</div>

		<Button
			variant="primary"
			size="lg"
			class="w-full"
			onclick={() => {
				if (!authStore.user) {
					goto(resolve('login'));
				} else {
					startFlow();
				}
			}}
		>
			{#if !authStore.user}
				Sign In to Play
			{:else}
				Choose AI Opponent
			{/if}
		</Button>
	{/if}

	<ChangelogButton />
</div>
