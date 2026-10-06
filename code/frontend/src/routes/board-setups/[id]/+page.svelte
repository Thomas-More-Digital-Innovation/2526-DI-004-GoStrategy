<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import SetupEditor from '$lib/components/setup/SetupEditor.svelte';
	import BoardSetupMetaForm from '$lib/components/setup/BoardSetupMetaForm.svelte';
	import Loading from '$lib/components/ui/Loading.svelte';
	import { boardSetups } from '$lib/api/client';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { toastStore } from '$lib/state/toast.svelte';
	import type { BoardSetup } from '$lib/types/board-setup';

	let id = $derived(Number(page.params.id));
	let setup = $state<BoardSetup | null>(null);
	let name = $state('');
	let description = $state('');
	let isDefault = $state(false);
	let loading = $state(true);
	let saving = $state(false);
	let loadError = $state('');

	onMount(async () => {
		try {
			const found = await boardSetups.getOne(id);
			if (found) {
				setup = found;
				name = found.name || '';
				description = found.description || '';
				isDefault = found.is_default || false;
			} else {
				loadError = 'Setup not found';
				toastStore.error('Setup not found');
			}
		} catch (e: any) {
			loadError = e.message || 'Failed to load setup';
			toastStore.handleApiMessage(e, 'Failed to load setup');
		} finally {
			loading = false;
		}
	});

	async function handleSave(setupData: string) {
		saving = true;
		try {
			await boardSetups.update(id, {
				name,
				description,
				setup_data: setupData,
				is_default: isDefault
			});
			goto('/board-setups');
		} catch (e: any) {
			toastStore.handleApiMessage(e, 'Failed to update setup');
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>GoStrategy — Edit Setup</title>
</svelte:head>

<div class="max-w-6xl mx-auto space-y-8">
	<div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
		<div class="flex-1 space-y-1">
			<h1 class="text-3xl font-black text-white uppercase tracking-tighter">Edit Setup</h1>
			<p class="text-white/40">Adjust your formation for upcoming battles.</p>
		</div>
	</div>

	{#if loading}
		<Loading
			title="Loading Setup"
			description="Retrieving your saved piece placements..."
			subtitle="Loading"
		/>
	{:else if loadError && !setup}
		<Card class="text-center py-12">
			<p class="text-brand-secondary">{loadError}</p>
			<Button variant="ghost" class="mt-4" onclick={() => goto('/board-setups')}
				>Back to List</Button
			>
		</Card>
	{:else if setup}
		<BoardSetupMetaForm bind:name bind:description bind:isDefault disabled={saving} />

		{#if saving}
			<Loading
				title="Saving Setup"
				description="Storing your updated board formation..."
				subtitle="Saving"
			/>
		{:else}
			<SetupEditor
				initialSetup={setup.setup_data}
				onSave={handleSave}
				onCancel={() => goto('/board-setups')}
			/>
		{/if}
	{/if}
</div>
