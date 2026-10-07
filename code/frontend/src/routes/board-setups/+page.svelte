<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { authStore } from '$lib/state/auth.svelte';
	import { boardSetups } from '$lib/api/client';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import type { BoardSetup } from '$lib/types/board-setup';
	import { MAX_BOARD_SETUPS } from '$lib/types/board-setup';
	import BoardSetupCard from '$lib/components/setup/BoardSetupCard.svelte';
	import { exportSetupAsJson } from '$lib/utils/board-binary';

	let setups = $state<BoardSetup[]>([]);
	let error = $state('');
	let loading = $state(true);

	onMount(async () => {
		await authStore.check();
		if (!authStore.user) {
			goto(resolve('/login'));
			return;
		}
		await loadSetups();
	});

	async function loadSetups() {
		loading = true;
		try {
			const result = await boardSetups.list();
			setups = result ?? [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load setups';
		} finally {
			loading = false;
		}
	}

	async function deleteSetup(id: number) {
		if (!confirm('Delete this setup?')) return;
		try {
			await boardSetups.delete(id);
			await loadSetups();
		} catch (e) {
			error = 'Failed to delete: ' + (e instanceof Error ? e.message : 'Unknown error');
		}
	}

	async function handleExportAll() {
		try {
			await boardSetups.exportAllZip();
		} catch (e) {
			error = 'Failed to export setups: ' + (e instanceof Error ? e.message : 'Unknown error');
		}
	}
</script>

<svelte:head>
	<title>GoStrategy — Board Setups</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-extrabold text-white uppercase tracking-widest">Board Setups</h1>
			<p class="text-white/40 text-sm mt-1">
				{setups.length}/{MAX_BOARD_SETUPS} setups saved
			</p>
		</div>
		<div class="flex items-center gap-3">
			{#if setups.length > 0}
				<Button variant="outline" onclick={handleExportAll}>Export All (ZIP)</Button>
			{/if}
			<Button
				variant="primary"
				disabled={setups.length >= MAX_BOARD_SETUPS}
				onclick={() => goto(resolve('/board-setups/new'))}
			>
				+ Create New
			</Button>
		</div>
	</div>

	{#if error}
		<Alert variant="error" message={error} />
	{/if}

	{#if loading}
		<div class="text-center py-12 text-white/40">Loading formations...</div>
	{:else if setups.length === 0}
		<Card class="text-center py-12">
			<p class="text-white/40">No board setups yet. Create your first battle formation!</p>
		</Card>
	{:else}
		<div class="flex flex-wrap gap-6">
			{#each setups as setup (setup.id)}
				<BoardSetupCard {setup} ownerId={1}>
					{#snippet actions()}
						<div class="flex gap-2">
							<Button
								variant="outline"
								size="sm"
								class="flex-1"
								onclick={() => goto(resolve(`/board-setups/${setup.id}`))}
							>
								Edit
							</Button>
							<Button variant="outline" size="sm" onclick={() => exportSetupAsJson(setup)}>
								Export
							</Button>
							<Button variant="ghost" size="sm" onclick={() => deleteSetup(setup.id)}>
								Delete
							</Button>
						</div>
					{/snippet}
				</BoardSetupCard>
			{/each}
		</div>
	{/if}
</div>
