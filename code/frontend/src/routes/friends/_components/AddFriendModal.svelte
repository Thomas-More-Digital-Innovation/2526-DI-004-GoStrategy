<script lang="ts">
	import { friends } from '#lib/api/client.js';
	import { authStore } from '#lib/state/auth.svelte.js';
	import { toastStore } from '#lib/state/toast.svelte.js';
	import Button from '#lib/components/ui/Button.svelte';
	import Modal from '#lib/components/ui/Modal.svelte';
	import type { UserDTO } from '#lib/types/generated/models.js';

	interface Props {
		isOpen: boolean;
		code: string;
		onClose: () => void;
		onSent: () => void;
	}

	let { isOpen = $bindable(false), code, onClose, onSent }: Props = $props();
	let target = $state<UserDTO | null>(null);
	let problem = $state('');
	let sending = $state(false);

	const isSelf = $derived(target !== null && target.id === authStore.user?.id);

	$effect(() => {
		if (isOpen) lookUp(code);
	});

	async function lookUp(requested: string) {
		target = null;
		problem = '';
		try {
			const result = await friends.preview(requested);
			if (requested === code) target = result.user;
		} catch (e) {
			if (requested === code) problem = e instanceof Error ? e.message : 'Unknown friend code';
		}
	}

	// why: the request is sent on click only, because link previews and prefetchers open URLs
	async function send() {
		sending = true;
		try {
			const result = await friends.sendRequest(code);
			toastStore.success(result.status === 'accepted' ? 'You are now friends' : 'Request sent');
			isOpen = false;
			onSent();
		} catch (e) {
			toastStore.handleApiMessage(e, 'Failed to send request');
		} finally {
			sending = false;
		}
	}
</script>

<Modal
	bind:isOpen
	title={target && !isSelf ? `Add ${target.username}?` : 'Add a friend'}
	maxWidth="md"
	{onClose}
>
	{#if problem}
		<p class="text-white/70">{problem}</p>
	{:else if isSelf}
		<p class="text-white/70">That is your own friend code.</p>
	{:else if target}
		<p class="text-white/70">They get a request and must accept it before you are friends.</p>
	{:else}
		<p class="text-white/50">Looking up friend code…</p>
	{/if}
	{#snippet actions()}
		<div class="flex justify-end gap-3">
			<Button variant="ghost" onclick={onClose}>Cancel</Button>
			{#if target && !isSelf}
				<Button loading={sending} onclick={send}>Send request</Button>
			{/if}
		</div>
	{/snippet}
</Modal>
