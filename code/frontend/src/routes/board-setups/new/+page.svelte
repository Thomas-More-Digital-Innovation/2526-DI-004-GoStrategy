<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import SetupEditor from '$lib/components/setup/SetupEditor.svelte';
	import BoardSetupMetaForm from '$lib/components/setup/BoardSetupMetaForm.svelte';
	import { boardSetups } from '$lib/api/client';
	import { toastStore } from '$lib/state/toast.svelte';

	let name = $state('My New Setup');
	let description = $state('');
	let isDefault = $state(false);
	let saving = $state(false);

	async function handleSave(setupData: string) {
		saving = true;
		try {
			await boardSetups.create({
				name,
				description,
				setup_data: setupData,
				is_default: isDefault
			});
			goto(resolve('/board-setups'));
		} catch (e) {
			toastStore.handleApiMessage(e, 'Failed to save setup');
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>GoStrategy — New Setup</title>
</svelte:head>

<div class="max-w-6xl mx-auto space-y-8">
	<div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
		<div class="flex-1 space-y-1">
			<h1 class="text-3xl font-black text-white uppercase tracking-tighter">Create New Setup</h1>
			<p class="text-white/40">Design your starting positions for the battlefield.</p>
		</div>
	</div>

	<BoardSetupMetaForm bind:name bind:description bind:isDefault disabled={saving} />

	{#if saving}
		<div class="text-center py-12 text-white/40">Saving setup...</div>
	{:else}
		<SetupEditor onSave={handleSave} onCancel={() => goto(resolve('/board-setups'))} />
	{/if}
</div>
